package db

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
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
					fmt.Println("WARNING: failed to sync logger: "+err.Error())
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

	dirEntries, err := ioutil.ReadDir(inputDataPath)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		date, err := parseDateFromName(dirEntry.Name())
		if err != nil {
			return errors.Wrapf(err, "failed to parse filename to YYYY-MM-DD date %s", dirEntry.Name())
		}

		dat, err := os.ReadFile(path.Join(inputDataPath, dirEntry.Name()))
		if err != nil {
			return errors.Wrapf(err, "failed to read %s", dirEntry.Name())
		}

		if err := conn.WithStore(func(s *ro.Store) error {
			return s.InsertChangelog(context.Background(), &models.Changelog{Date: date, Content: string(dat)})
		}); err != nil {
			return err
		}
	}

	return nil
}

func parseDateFromName(name string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSuffix(name, ".md"))
}
