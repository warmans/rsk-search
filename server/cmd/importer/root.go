package importer

import (
	"github.com/spf13/cobra"
)

type importerConfig struct {
	dataDir  string
	audioDir string
	videoDir string
	imageDir string
}

var cfg = importerConfig{}

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "db",
		Short: "commands related to to the search index",
	}

	root.PersistentFlags().StringVarP(&cfg.dataDir, "data-dir", "d", "./var/data/episodes", "Path to the raw data files")
	root.PersistentFlags().StringVarP(&cfg.audioDir, "audio-dir", "a", "", "Path to the audio files")
	root.PersistentFlags().StringVarP(&cfg.videoDir, "video-dir", "v", "", "Path to the video files")
	root.PersistentFlags().StringVarP(&cfg.imageDir, "image-dir", "", "./var/images", "Path to generated images")

	//index.AddCommand(LoadCmd())

	return root
}
