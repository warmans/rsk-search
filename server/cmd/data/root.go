package data

import (
	"github.com/spf13/cobra"
	"os"
)

type dataConfig struct {
	dataDir  string
	audioDir string
	videoDir string
	imageDir string
}

var cfg = dataConfig{}

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "data",
		Short: "commands related to to the search index",
	}

	root.PersistentFlags().StringVarP(&cfg.dataDir, "data-dir", "d", "./var/data/episodes", "Path to the raw data files")
	root.PersistentFlags().StringVarP(&cfg.audioDir, "audio-dir", "a", os.Getenv("AUDIO_DIR"), "Path to the audio files")
	root.PersistentFlags().StringVarP(&cfg.videoDir, "video-dir", "v", os.Getenv("VIDEO_DIR"), "Path to the video files")
	root.PersistentFlags().StringVarP(&cfg.imageDir, "image-dir", "", "./var/images", "Path to generated images")

	root.AddCommand(InitCmd())
	root.AddCommand(InitFromAudioFilesCmd())
	root.AddCommand(ImportPilkipediaRaw())
	root.AddCommand(ImportSpotifyData())
	root.AddCommand(InitFromSrtCmd())

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
	root.AddCommand(AssemblyAI2Srt())

	//openai testing
	root.AddCommand(TranscribeOpenAICommand())

	//misc
	root.AddCommand(CountWords())
	root.AddCommand(MergeTimestampsAAICommand())
	root.AddCommand(RefreshAudioMetadataCmd())
	root.AddCommand(AddRemasterAudioCmd())

	//video
	root.AddCommand(ExtractVideoImages())

	return root
}
