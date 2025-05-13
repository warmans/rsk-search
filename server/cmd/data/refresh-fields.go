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

				switch episode.Publication {
				case "office", "extras":
					episode.PublicationType = models.PublicationTypeTV
				case "guide", "podcast", "fame":
					episode.PublicationType = models.PublicationTypePodcast
				case "xfm", "nme", "bbc2":
					episode.PublicationType = models.PublicationTypeRadio
				case "preview":
					episode.PublicationType = models.PublicationTypePromo
				default:
					episode.PublicationType = models.PublicationTypeOther
				}

				// don't override reported incompleteness
				if len(episode.Transcript) == 0 {
					episode.Completion = models.CompletionStateEmpty
				} else {
					if episode.Completion != models.CompletionStateReported {
						hasGaps := len(episode.Transcript) == 0
						for k, v := range episode.Transcript {
							episode.Transcript[k].Position = int64(k + 1)
							if v.Type == models.DialogTypeGap || v.Placeholder {
								hasGaps = true
							}
						}
						if hasGaps {
							episode.Completion = models.CompletionStateGaps
						} else {
							episode.Completion = models.CompletionStateComplete
							// some episodes are locked by default, but they shouldn't be if they're complete
							episode.Locked = false
						}
					}
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
