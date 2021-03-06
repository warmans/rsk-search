package transcription

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	index := &cobra.Command{
		Use:   "transcription",
		Short: "commands related to audio transcription",
	}

	index.AddCommand(GcloudCmd())
	index.AddCommand(ImportRawCmd())
	index.AddCommand(MapChunksCmd())

	return index
}
