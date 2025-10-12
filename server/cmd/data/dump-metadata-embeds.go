package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
	"slices"
	"time"
)

type episodeDate struct {
	epid    string
	date    *time.Time
	episode int32
}

func DumpMetadataEmbedsCmd() *cobra.Command {

	var inputDir string
	var outputDir string

	publicationWhitelist := map[string]struct{}{
		"xfm":     {},
		"bbc2":    {},
		"fame":    {},
		"name":    {},
		"guide":   {},
		"podcast": {},
	}

	cmd := &cobra.Command{
		Use:   "dump-metadata-embeds",
		Short: "refresh data in meta/data directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			episodeIDs := make([]episodeDate, 0)

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

				if _, ok := publicationWhitelist[episode.Publication]; !ok {
					continue
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				episodeIDs = append(episodeIDs, episodeDate{epid: episode.ShortID(), date: episode.ReleaseDate, episode: episode.Episode})
			}

			slices.SortFunc(episodeIDs, func(a, b episodeDate) int {
				if a.date == nil {
					return -1
				}
				if b.date == nil {
					return 1
				}
				if a.date.Before(*b.date) {
					return -1
				}
				if a.date.Equal(*b.date) {
					if a.episode < b.episode {
						return -1
					}
				}
				return 1
			})

			dateMap := make(map[string]string)
			episodeOrder := make([]string, 0)

			for _, v := range episodeIDs {
				dateMap[v.date.Format(time.RFC3339)] = v.epid
				episodeOrder = append(episodeOrder, v.epid)
			}

			if err := util.WithReplaceJSONFileEncoder(path.Join(outputDir, "episode-date-map.json"), func(enc *json.Encoder) error {
				return enc.Encode(dateMap)
			}); err != nil {
				return err
			}

			if err := util.WithReplaceJSONFileEncoder(path.Join(outputDir, "episode-order.json"), func(enc *json.Encoder) error {
				return enc.Encode(episodeOrder)
			}); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw scraped files")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "./pkg/meta/data", "Path to output JSON")

	return cmd
}
