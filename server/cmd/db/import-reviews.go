package db

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
)

func ImportReviewsCmd() *cobra.Command {

	dbCfg := &common.Config{}
	var csvPath string
	var externalAuthorName string

	cmd := &cobra.Command{
		Use:   "import-reviews",
		Short: "Load reviews from a CV with a particular format",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			if dbCfg.DSN == "" {
				panic("dsn not set")
			}
			conn, err := rw.NewConn(dbCfg)
			if err != nil {
				return err
			}

			f, err := os.Open(csvPath)
			if err != nil {
				return fmt.Errorf("failed to open CSV: %w", err)
			}
			defer f.Close()

			return processFile(conn, logger, externalAuthorName, csv.NewReader(f))
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	flag.StringVarEnv(cmd.Flags(), &csvPath, "", "csv-path", "reviews.csv", "Path to reviews csv")
	flag.StringVarEnv(cmd.Flags(), &externalAuthorName, "", "author-name", "Anonymous", "Attribution for reviews in file")

	return cmd
}

func processFile(conn *rw.Conn, logger *zap.Logger, authorName string, csvFile *csv.Reader) error {

	var authorID string
	if err := conn.WithStore(func(s *rw.Store) error {
		var err error
		authorID, err = s.UpsertAuthorPlaceholder(context.Background(), authorName)
		return err
	}); err != nil {
		return err
	}

	logger.Info("Upserted author", zap.String("author_id", authorID))

	for {
		record, err := csvFile.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
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

		if err := conn.WithStore(func(s *rw.Store) error {
			return s.AddAuthorReview(
				context.Background(),
				authorID,
				episodeID,
				float32(ratingFloat),
				nil,
			)
		}); err != nil {
			logger.Error("failed to insert review ", zap.String("epid", episodeID), zap.Error(err))
			continue
		}
	}

	return nil
}
