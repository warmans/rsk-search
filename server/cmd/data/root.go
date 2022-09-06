package data

import (
	"github.com/spf13/cobra"
)

type dataConfig struct {
	dataDir  string
	audioDir string
}

var cfg = dataConfig{}

func RootCmd() *cobra.Command {
	index := &cobra.Command{
		Use:   "data",
		Short: "commands related to to the search index",
	}

	index.PersistentFlags().StringVarP(&cfg.dataDir, "data-dir", "d", "./var/data/episodes", "Path to the raw data files")
	index.PersistentFlags().StringVarP(&cfg.audioDir, "audio-dir", "a", "", "Path to the audio files")

	index.AddCommand(InitCmd())
	index.AddCommand(ImportPilkipediaRaw())
	index.AddCommand(ImportSpotifyData())

	// exports
	index.AddCommand(GenerateHTMLCmd())
	index.AddCommand(InferMissingOffsetsCmd())
	index.AddCommand(RefreshCmd())
	index.AddCommand(DumpPlaintext())

	// index

	index.AddCommand(PopulateBlugeIndex())

	return index
}
