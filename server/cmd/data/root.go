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
	root := &cobra.Command{
		Use:   "data",
		Short: "commands related to to the search index",
	}

	root.PersistentFlags().StringVarP(&cfg.dataDir, "data-dir", "d", "./var/data/episodes", "Path to the raw data files")
	root.PersistentFlags().StringVarP(&cfg.audioDir, "audio-dir", "a", "", "Path to the audio files")

	root.AddCommand(InitCmd())
	root.AddCommand(InitFromAudioFilesCmd())
	root.AddCommand(ImportPilkipediaRaw())
	root.AddCommand(ImportSpotifyData())

	// exports
	root.AddCommand(GenerateHTMLCmd())
	root.AddCommand(InferMissingOffsetsCmd())
	root.AddCommand(RefreshCmd())
	root.AddCommand(DumpPlaintext())

	// index
	root.AddCommand(PopulateBlugeIndex())

	// assembly ai testing
	root.AddCommand(TranscribeAssemblyAICmd())
	root.AddCommand(AssemblyAI2Dialog())

	//openai testing
	root.AddCommand(TranscribeOpenAICommand())

	//misc
	root.AddCommand(CountWords())
	root.AddCommand(MergeTimestampsAAICommand())
	root.AddCommand(RefreshAudioMetadataCmd())

	return root
}
