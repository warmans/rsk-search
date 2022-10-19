package db

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/transcript"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

func init() {
	stdoutPrinter.SetIndent("  ", "  ")
}

func MergeTranscriptChangesCmd() *cobra.Command {
	var outputDir string
	var migrationsPath string
	var dryRun bool
	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "merge-transcript-changes",
		Short: "Replace flat file transcripts from transcript change records and mark as merged",
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
			defer conn.Close()

			return mergeAll(outputDir, migrationsPath, conn, dryRun, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "var/data/episodes", "Path to output the data")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", true, "Instead of saving the data print it out standard out")
	cmd.Flags().StringVarP(&migrationsPath, "migrations-path", "m", "pkg/store/rw/migrations", "Migrations are written to mark changes as merged")

	return cmd
}

func mergeAll(outputDataPath string, migrationsPath string, conn *rw.Conn, dryRun bool, logger *zap.Logger) error {

	logger.Info("Reading DB...")

	ctx := context.Background()

	approvedChangeIDs := []string{}
	err := conn.WithStore(func(s *rw.Store) error {

		approvedChanges, err := s.ListTranscriptChanges(
			ctx,
			common.Q(common.WithFilter(
				filter.And(
					filter.Eq("state", filter.String(string(models.ContributionStateApproved))),
					filter.Eq("merged", filter.Bool(false)),
				)),
			),
		)
		if err != nil {
			return err
		}

		for _, v := range approvedChanges {

			logger.Info(fmt.Sprintf("Processing change %s (%s)...", v.ID, v.EpID))

			episodeOnDisk, err := data.LoadEpisdeByEpisodeID(outputDataPath, v.EpID)
			if err != nil {
				return err
			}
			if episodeOnDisk == nil {
				panic("nil episode encountered: " + v.EpID)
			}

			// clear old data
			episodeOnDisk.Synopsis = nil
			episodeOnDisk.Transcript = nil
			episodeOnDisk.Trivia = nil

			// contributors should be merged with whatever is on disk
			uniqueContributors := map[string]struct{}{}
			for _, v := range episodeOnDisk.Contributors {
				uniqueContributors[v] = struct{}{}
			}

			dialog, synopsis, trivia, err := transcript.Import(bufio.NewScanner(bytes.NewBufferString(v.Transcription)), episodeOnDisk.ID(), 1)
			if err != nil {
				return err
			}

			// process contributors for this chunk of audio
			uniqueContributors[v.Author.Name] = struct{}{}
			episodeOnDisk.Transcript = append(episodeOnDisk.Transcript, dialog...)
			episodeOnDisk.Synopsis = append(episodeOnDisk.Synopsis, synopsis...)
			episodeOnDisk.Trivia = append(episodeOnDisk.Trivia, trivia...)

			// metadata
			episodeOnDisk.Summary = v.Summary

			contributors := []string{}
			for c := range uniqueContributors {
				contributors = append(contributors, c)
			}
			sort.Strings(contributors)
			episodeOnDisk.Contributors = contributors

			// update version
			switch true {
			case v.PointsAwarded > 0 && v.PointsAwarded <= 1:
				episodeOnDisk.Version, err = util.NextVersion(episodeOnDisk.Version, util.PatchVersion)
			case v.PointsAwarded > 1 && v.PointsAwarded <= 2:
				episodeOnDisk.Version, err = util.NextVersion(episodeOnDisk.Version, util.MinorVersion)
			case v.PointsAwarded > 2:
				episodeOnDisk.Version, err = util.NextVersion(episodeOnDisk.Version, util.MajorVersion)
			}
			if err != nil {
				return errors.Wrap(err, "failed to update version")
			}

			if dryRun {
				if err := stdoutPrinter.Encode(episodeOnDisk); err != nil {
					return err
				}
			} else {
				if err := data.ReplaceEpisodeFile(outputDataPath, episodeOnDisk); err != nil {
					return err
				}
				approvedChangeIDs = append(approvedChangeIDs, v.ID)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(approvedChangeIDs) > 0 {
		return generateMigration(migrationsPath, approvedChangeIDs)
	}

	return nil
}

// The change is only technically merged once it goes live so the merged flag must be set using a migration.
func generateMigration(migrationsPath string, approvedChangeIDs []string) error {
	for k, v := range approvedChangeIDs {
		approvedChangeIDs[k] = fmt.Sprintf(`'%s'`, v)
	}
	return os.WriteFile(
		path.Join(migrationsPath, fmt.Sprintf("%d_merge_changes.sql", time.Now().Unix())),
		[]byte(fmt.Sprintf("UPDATE transcript_change SET merged=true WHERE id IN (%s);", strings.Join(approvedChangeIDs, ", "))),
		0666,
	)
}
