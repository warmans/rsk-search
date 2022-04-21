package db

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"strings"
)

func LoadTscriptCmd() *cobra.Command {

	var inputDir string
	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "load-tscript",
		Short: "load chunked, incomplete transcripts into persistent DB. Duplicate IDs are discarded.",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					panic("failed to sync logger: "+err.Error())
				}
			}()

			if dbCfg.DSN == "" {
				panic("dsn not set")
			}
			conn, err := rw.NewConn(dbCfg)
			if err != nil {
				return err
			}
			defer conn.Close()

			return populatePersistentDB(inputDir, conn, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "var/data/incomplete/chunked", "Path to chunks transcripts")

	return cmd
}

func populatePersistentDB(inputDataPath string, conn *rw.Conn, logger *zap.Logger) error {

	logger.Info("Populating DB...")

	incompleteEntries, err := ioutil.ReadDir(inputDataPath)
	if err != nil {
		return err
	}
	for _, dirEntry := range incompleteEntries {
		if dirEntry.IsDir() || strings.HasPrefix(dirEntry.Name(), ".") {
			continue
		}
		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		tscript := &models.Tscript{}
		if err := util.WithReadJSONFileDecoder(path.Join(inputDataPath, dirEntry.Name()), func(dec *json.Decoder) error {
			return dec.Decode(tscript)
		}); err != nil {
			return err
		}
		if err := conn.WithStore(func(s *rw.Store) error {
			return s.InsertOrIgnoreTscript(context.Background(), tscript)
		}); err != nil {
			return err
		}
	}

	return nil
}
