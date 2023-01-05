package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
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

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw data files")
	cmd.Flags().StringVarP(&dbDSN, "db-dsn", "d", "./var/gen/ro.sqlite3", "databsae DSN")

	return cmd
}

func populateDB(inputDataPath string, conn *ro.Conn, logger *zap.Logger) error {

	logger.Info("Populating DB...")

	dirEntries, err := ioutil.ReadDir(inputDataPath)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		episode := &models.Transcript{}
		if err := util.WithReadJSONFileDecoder(path.Join(inputDataPath, dirEntry.Name()), func(dec *json.Decoder) error {
			return dec.Decode(episode)
		}); err != nil {
			return err
		}

		if err := conn.WithStore(func(s *ro.Store) error {
			return s.InsertEpisodeWithTranscript(context.Background(), episode)
		}); err != nil {
			return err
		}
	}

	return nil
}
