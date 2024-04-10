package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
)

func LoadCmd() *cobra.Command {

	var inputDir string
	var dbDSN string

	cmd := &cobra.Command{
		Use:   "load-ro",
		Short: "refresh the search index from the given directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			if dbDSN == "" {
				panic("dsn not set")
			}
			conn, err := ro.NewConn(&common.Config{
				DSN: dbDSN,
			})
			if err != nil {
				return err
			}
			if err := conn.Migrate(); err != nil {
				return err
			}
			return populateDB(inputDir, conn, logger)
		},
	}

	cmd.Flags().StringVarP(&inputDir, "data-dir", "i", "./var/data", "Path to data dir")
	cmd.Flags().StringVarP(&dbDSN, "db-dsn", "d", "./var/gen/ro.sqlite3", "databsae DSN")

	return cmd
}

func populateDB(dataDir string, conn *ro.Conn, logger *zap.Logger) error {

	logger.Info("Populating DB...")
	ctx := context.Background()

	dirEntries, err := os.ReadDir(path.Join(dataDir, "episodes"))
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		episode := &models.Transcript{}
		if err := util.WithReadJSONFileDecoder(path.Join(dataDir, "episodes", dirEntry.Name()), func(dec *json.Decoder) error {
			return dec.Decode(episode)
		}); err != nil {
			return err
		}

		if err := conn.WithStore(func(s *ro.Store) error {
			return s.InsertEpisodeWithTranscript(ctx, episode)
		}); err != nil {
			return err
		}
	}

	logger.Info("Loading songs...")
	err = conn.WithStore(func(s *ro.Store) error {
		for _, v := range meta.GetSongMeta().ExtractSorted() {
			if err := s.InsertSong(ctx, v); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load songs: %w", err)
	}

	logger.Info("Loading community projects...")
	err = conn.WithStore(func(s *ro.Store) error {
		projects := []models.CommunityProject{}
		if err := util.WithReadJSONFileDecoder(path.Join(dataDir, "community", "projects.json"), func(dec *json.Decoder) error {
			return dec.Decode(&projects)
		}); err != nil {
			return fmt.Errorf("failed to decode project JSON: %w", err)
		}
		for _, v := range projects {
			if err := s.InsertCommunityProject(ctx, v); err != nil {
				return fmt.Errorf("failed to insert project JSON: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load community projects: %w", err)
	}

	return nil
}
