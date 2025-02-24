package db

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
	"time"
)

func init() {
	stdoutPrinter.SetIndent("  ", "  ")
}

func MergeTranscriptRatingsCmd() *cobra.Command {
	var outputDir string
	var migrationsPath string
	var dryRun bool
	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "merge-transcript-ratings",
		Short: "Merge ratings to flat files from DB",
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

			return mergeRatings(outputDir, migrationsPath, conn, dryRun, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "./var/data/episodes", "Path to output the data")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", true, "Instead of saving the data print it out standard out")
	cmd.Flags().StringVarP(&migrationsPath, "migrations-path", "m", "pkg/store/rw/migrations", "Migrations are written to mark changes as merged")

	return cmd
}

func mergeRatings(outputDataPath string, migrationsPath string, conn *rw.Conn, dryRun bool, logger *zap.Logger) error {

	logger.Info("Reading DB...")

	ctx := context.Background()

	// IDs are a composite key of username + episode ID
	mergedRatingIDs := [][2]string{}
	err := conn.WithStore(func(s *rw.Store) error {

		ratings, err := s.ListTranscriptRatingScores(
			ctx,
		)
		if err != nil {
			return err
		}

		for _, v := range ratings {

			logger.Info(fmt.Sprintf("Processing rating %s/%s...", v.AuthorID, v.EpisodeID))

			episodeOnDisk, err := data.LoadEpisdeByEpisodeID(outputDataPath, fmt.Sprintf("ep-%s", v.EpisodeID))
			if err != nil {
				return err
			}
			if episodeOnDisk == nil {
				panic("nil episode encountered: " + v.EpisodeID)
			}

			// clear old data
			if episodeOnDisk.Ratings.Scores == nil {
				episodeOnDisk.Ratings.Scores = make(map[string]float32)
			}

			episodeOnDisk.Ratings.Scores[v.AuthorIdentifier] = v.Score

			if dryRun {
				if err := stdoutPrinter.Encode(episodeOnDisk); err != nil {
					return err
				}
			} else {
				if err := data.ReplaceEpisodeFile(outputDataPath, episodeOnDisk); err != nil {
					return err
				}
				mergedRatingIDs = append(mergedRatingIDs, [2]string{v.AuthorID, v.EpisodeID})
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(mergedRatingIDs) > 0 {
		return generateRatingMigrations(migrationsPath, mergedRatingIDs)
	}

	return nil
}

// The change is only technically merged once it goes live so the merged flag must be set using a migration.
func generateRatingMigrations(migrationsPath string, mergedRatingIDs [][2]string) error {
	idsPairs := make([]string, len(mergedRatingIDs))
	for k, v := range mergedRatingIDs {
		idsPairs[k] = fmt.Sprintf(`('%s', '%s')`, v[0], v[1])
	}
	return os.WriteFile(
		path.Join(migrationsPath, fmt.Sprintf("%d_merge_ratings.sql", time.Now().Unix())),
		[]byte(fmt.Sprintf(
			`DELETE FROM transcript_rating_score WHERE (author_id, episode_id) IN (%s);`,
			strings.Join(idsPairs, ", "),
		)),
		0666,
	)
}
