package db

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/changelog"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"go.uber.org/zap"
)

func LoadChangelogs() *cobra.Command {

	var inputDir string
	var dbDSN string

	cmd := &cobra.Command{
		Use:   "load-changelogs",
		Short: "load changelog data into readonly-db",
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
			return populateChangelog(inputDir, conn, logger)
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/changelogs", "Path to raw markdown files")
	cmd.Flags().StringVarP(&dbDSN, "db-dsn", "d", "./var/ro.sqlite3", "readonly database DSN")

	return cmd
}

func populateChangelog(inputDataPath string, conn *ro.Conn, logger *zap.Logger) error {

	logger.Info("Populating DB...")

	changeLogs, err := changelog.List(inputDataPath)
	if err != nil {
		return err
	}
	for _, log := range changeLogs {
		if err := conn.WithStore(func(s *ro.Store) error {
			return s.InsertChangelog(context.Background(), log)
		}); err != nil {
			return err
		}
	}

	return nil
}
