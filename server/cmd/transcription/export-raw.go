package transcription

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/transcript"
	"go.uber.org/zap"
	"os"
)

func ExportRaw() *cobra.Command {

	var inputFile string

	cmd := &cobra.Command{
		Use:   "export-raw",
		Short: "Convert a transcript to HTML (mediawiki compatible)",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					panic("failed to sync logger: "+err.Error())
				}
			}()

			episode, err := data.LoadEpisodePath(inputFile)
			if err != nil {
				return err
			}

			raw, err := transcript.Export(episode.Transcript, episode.Synopsis, episode.Trivia)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, raw)
			return err
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "./var/data/episodes/ep-xfm-S2E28.json", "Path input JSON")

	return cmd
}
