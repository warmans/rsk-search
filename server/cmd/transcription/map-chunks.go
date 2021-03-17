package transcription

import (
	"github.com/spf13/cobra"
)

func MapChunksCmd() *cobra.Command {

	var incompleteFilesDir string

	cmd := &cobra.Command{
		Use:   "map-chunks",
		Short: "Take raw transcript and map it to a series of chunks that can be consumed by the API",
		RunE: func(cmd *cobra.Command, args []string) error {

			// todo for each file scan the lines and create a series of ~2-3 min chunks
			// just dump them all to JSON for now and they can be loaded into a DB

			return nil
		},
	}

	cmd.Flags().StringVarP(&incompleteFilesDir, "incomplete-dir", "d", "./var/data/incomplete", "Path to incomplete transcripts output by gcloud")

	return cmd
}
