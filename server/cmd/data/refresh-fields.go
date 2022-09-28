package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

func RefreshCmd() *cobra.Command {

	var inputDir string
	var singleEpisode string

	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "mark any episodes without gaps as complete",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			dirEntries, err := ioutil.ReadDir(inputDir)
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

				// set an initial version if missing.
				if episode.Version == "" {
					episode.Version = "0.0.0"
				}

				if singleEpisode != "" {
					if episode.ShortID() != singleEpisode {
						continue
					}
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				// identify gaps
				hasGaps := false
				for _, v := range episode.Transcript {
					if v.Type == models.DialogTypeGap {
						hasGaps = true
					}
				}
				if hasGaps {
					episode.Incomplete = true
				} else {
					episode.Incomplete = false
				}
				
				// ensure IDs are correct
				for k := range episode.Transcript {
					episode.Transcript[k].ID = models.DialogID(episode.ID(), episode.Transcript[k].Position)
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
