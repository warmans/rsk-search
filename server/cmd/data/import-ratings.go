package data

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
)

func ImportRatingsCmd() *cobra.Command {

	dbCfg := &common.Config{}
	var csvPath string
	var externalAuthorName string

	cmd := &cobra.Command{
		Use:   "import-ratings",
		Short: "Load reviews from a CV with a particular format",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			f, err := os.Open(csvPath)
			if err != nil {
				return fmt.Errorf("failed to open CSV: %w", err)
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					logger.Error("failed to close file", zap.Error(err))
				}
			}(f)

			return processFile(logger, externalAuthorName, csv.NewReader(f))
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	flag.StringVarEnv(cmd.Flags(), &csvPath, "", "csv-path", "reviews.csv", "Path to reviews csv")
	flag.StringVarEnv(cmd.Flags(), &externalAuthorName, "", "author-identifier", "discord:anon", "Attribution for reviews in file")

	return cmd
}

func processFile(logger *zap.Logger, authorName string, csvFile *csv.Reader) error {

	for {
		record, err := csvFile.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		episodeID := fmt.Sprintf("%s-%s", record[0], record[1])

		// validate
		publication, series, episode, err := models.ParseEpID(episodeID)
		if err != nil {
			logger.Error("failed to parse ID", zap.String("epid", episodeID), zap.Error(err))
			continue
		}
		// ensure superfluous chars are removed
		episodeID = models.ShortEpID(publication, series, episode)

		ep, err := data.LoadEpisodeByShortID(cfg.dataDir, episodeID)
		if err != nil {
			return err
		}
		if ep == nil {
			logger.Warn("No episode found, skipping", zap.String("epid", episodeID))
			continue
		}

		if record[2] == "" {
			continue
		}
		ratingFloat, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			logger.Error("failed to parse score", zap.String("epid", episodeID), zap.String("score", record[2]), zap.Error(err))
			continue
		}
		logger.Info(
			"Processing",
			zap.String("publication", publication),
			zap.Int32("series", series),
			zap.Int32("episode", episode),
			zap.Float32("score", float32(ratingFloat)),
		)

		if ep.Ratings.Scores == nil {
			ep.Ratings.Scores = make(map[string]float32)
		}

		ep.Ratings.Scores[authorName] = float32(ratingFloat)

		if err := data.ReplaceEpisodeFile(cfg.dataDir, ep); err != nil {
			return err
		}
	}

	return nil
}
