package db

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
)

func ExtractTscriptRawCmd() *cobra.Command {
	var outputDir string
	var dryRun bool
	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "extract-tscript-raw",
		Short: "Extract all transcripts in their original plaintext format",
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

			return extractRaw(outputDir, conn, dryRun, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "var/data/raw", "Path to output the data")

	return cmd
}

func extractRaw(outputDataPath string, conn *rw.Conn, dryRun bool, logger *zap.Logger) error {

	logger.Info("Reading DB...")

	ctx := context.Background()

	err := conn.WithStore(func(s *rw.Store) error {

		allTscripts, err := s.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, v := range allTscripts {

			logger.Info(fmt.Sprintf("Processing tscript %s-%s...", v.Publication, models.FormatStandardEpisodeName(v.Series, v.Episode)))

			outputFile, err := os.Create(path.Join(outputDataPath, fmt.Sprintf("%s.raw.txt", models.FormatStandardEpisodeName(v.Series, v.Episode))))
			if err != nil {
				return errors.Wrap(err, "failed to create output file")
			}

			allChunks, err := s.ListChunks(
				ctx,
				&common.QueryModifier{
					Filter: filter.Eq("tscript_id", filter.String(v.ID)),
					Sorting: &common.Sorting{
						Field:     "start_second",
						Direction: common.SortAsc,
					},
				},
			)
			if err != nil {
				return err
			}

			approved, err := s.ListChunkContributions(
				ctx,
				&common.QueryModifier{
					Filter: filter.And(
						filter.Eq("tscript_id", filter.String(v.ID)),
						filter.Eq("state", filter.String("approved")),
					),
				},
			)
			if err != nil {
				return err
			}
			if len(approved) == 0 {
				logger.Info("Nothing to do - none approved")
				continue
			}
			for _, ch := range allChunks {

				var chContribution *models.ChunkContribution
				for _, co := range approved {
					if co.ChunkID == ch.ID {
						chContribution = co
					}
				}

				// if the transcript is missing insert a placeholder
				if chContribution == nil {
					if _, err := outputFile.WriteString("none: [missing transcript section]\n"); err != nil {
						return err
					}
					continue
				}

				scn := bufio.NewScanner(bytes.NewBuffer([]byte(chContribution.Transcription)))
				for scn.Scan() {
					line := strings.TrimSpace(strings.ReplaceAll(scn.Text(), "\u00a0", " "))
					if line == "" {
						continue
					}
					if _, err := outputFile.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
						return err
					}
				}
				if err := scn.Err(); err != nil {
					return err
				}
			}
			if err := outputFile.Close(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
