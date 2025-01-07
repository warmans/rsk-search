package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
)

func AddRemasterAudioCmd() *cobra.Command {

	var inputDir string
	var audioDir string
	var singleEpisode string

	cmd := &cobra.Command{
		Use:   "add-remaster-audio",
		Short: "Add the remastered audio file names to metadata",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			logger.Info("Loading transcript data from...", zap.String("path", inputDir), zap.String("audio_dir", audioDir))

			dirEntries, err := os.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {

				if dirEntry.IsDir() {
					continue
				}

				episode := &models.Transcript{}
				if err := util.WithReadJSONFileDecoder(path.Join(inputDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				if singleEpisode != "" {
					if episode.ShortID() != singleEpisode {
						continue
					}
				}

				if episode.ReleaseDate == nil || episode.Publication != "xfm" {
					continue
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				found, err := findAudioFile(path.Join(audioDir, "remaster"), episode)
				if err != nil {
					return err
				}
				if found == "" {
					fmt.Printf("NOTHING FOUND FOR: %s\n", episode.ID())
					continue
				}

				episode.Media.RemasteredAudioFileName = found

				durationMs, err := getAudioDurationMs(path.Join(audioDir, "remaster", found))
				if err != nil {
					fmt.Println("Failed to get duration: ", err.Error())
				} else {
					episode.Media.RemasteredAudioDurationMs = durationMs
				}

				if err := data.ReplaceEpisodeFile(inputDir, episode); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw scraped files")
	cmd.Flags().StringVarP(&singleEpisode, "single-episode", "s", "", "Only process the given episode e.g. xfm-S2E04")
	cmd.Flags().StringVarP(&audioDir, "audio-dir", "a", os.Getenv("AUDIO_DIR"), "Root audio dir")

	return cmd
}

func findAudioFile(remasteredAudioDir string, episode *models.Transcript) (string, error) {
	entries, err := os.ReadDir(remasteredAudioDir)
	if err != nil {
		return "", err
	}
	for _, f := range entries {
		if f.IsDir() || !strings.HasSuffix(f.Name(), "mp3") {
			continue
		}
		if strings.HasPrefix(f.Name(), episode.ReleaseDate.Format("2006-01-02")) {
			return f.Name(), nil
		}
	}
	return "", nil
}
