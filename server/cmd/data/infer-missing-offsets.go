package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/transcript"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
)

func InferMissingOffsetsCmd() *cobra.Command {

	var inputDir string
	var singleEpisode string

	cmd := &cobra.Command{
		Use:   "infer-missing-offsets",
		Short: "adds time offsets to dialog based on approximate words per min spoken",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

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
				if len(episode.Transcript) == 0 {
					continue
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				if episode = transcript.InferOffsets(episode); episode == nil {
					continue
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

	return cmd
}
