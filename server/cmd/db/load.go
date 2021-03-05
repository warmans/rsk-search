package db

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"strconv"
)

func LoadCmd() *cobra.Command {

	var inputDir string
	var dbDSN string

	cmd := &cobra.Command{
		Use:   "load",
		Short: "refresh the search index from the given directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			if dbDSN == "" {
				panic("dsn not set")
			}
			conn, err := store.NewConn(&store.Config{
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

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw data files")
	cmd.Flags().StringVarP(&dbDSN, "db-dsn", "d", "./var/rsk.sqlite3", "databsae DSN")

	return cmd
}

func populateDB(inputDataPath string, conn *store.Conn, logger *zap.Logger) error {

	logger.Info("Populating DB...")

	dirEntries, err := ioutil.ReadDir(inputDataPath)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {

		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		episode := &models.Episode{}
		if err := util.WithJSONFileDecoder(path.Join(inputDataPath, dirEntry.Name()), func(dec *json.Decoder) error {
			return dec.Decode(episode)
		}); err != nil {
			return err
		}

		if err := conn.WithStore(func(s *store.Store) error {
			return s.InsertEpisodeWithTranscript(context.Background(), episode)
		}); err != nil {
			return err
		}
	}

	return nil
}

func stringToIntOrZero(str string) int32 {
	i, _ := strconv.Atoi(str)
	return int32(i)
}
