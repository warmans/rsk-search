package transcription

import (
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"os"
)

func GcloudCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "gcloud",
		Short: "use google cloud transcription service",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Argument should be a google cloud storage gs:// URI to a .wav file.
			if len(args) != 1 {
				return fmt.Errorf("gs URI was missing")
			}

			// GOOGLE_APPLICATION_CREDENTIALS should be a path to a credentials file e.g. /home/foo/credentials.json
			if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
				return fmt.Errorf("no google application credentials")
			}

			ctx := context.Background()
			client, err := speech.NewClient(ctx)
			if err != nil {
				return err
			}

			// Send the contents of the audio file with the encoding and
			// and sample rate information to be transcripted.
			req := &speechpb.LongRunningRecognizeRequest{
				Config: &speechpb.RecognitionConfig{
					Encoding:                   speechpb.RecognitionConfig_LINEAR16,
					LanguageCode:               "en-US",
					AudioChannelCount:          1,
					EnableWordTimeOffsets:      true,
					EnableAutomaticPunctuation: true,
					Model:                      "video",
				},
				Audio: &speechpb.RecognitionAudio{
					AudioSource: &speechpb.RecognitionAudio_Uri{Uri: args[0]},
				},
			}

			fmt.Fprintln(os.Stderr, "requesting...")
			op, err := client.LongRunningRecognize(ctx, req)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stderr, "waiting...")
			resp, err := op.Wait(ctx)
			if err != nil {
				return err
			}

			// Print the results.
			for _, result := range resp.Results {
				for _, alt := range result.Alternatives {
					for _, v := range alt.Words {
						if _, err := fmt.Fprintf(os.Stdout, "#OFFSET: %d\n", v.StartTime.Seconds); err != nil {
							return err
						}
						break
					}
					if _, err := fmt.Fprintf(os.Stdout, "Unknown: %s\n", alt.Transcript); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	return cmd
}
