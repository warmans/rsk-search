package db

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
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

func MergeTranscriptTagsCmd() *cobra.Command {
	var outputDir string
	var migrationsPath string
	var dryRun bool
	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "merge-transcript-tags",
		Short: "Merge tags to flat files from DB",
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
			defer func(conn *rw.Conn) {
				_ = conn.Close()
			}(conn)

			return mergeTags(outputDir, migrationsPath, conn, dryRun, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "./var/data/episodes", "Path to output the data")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", true, "Instead of saving the data print it out standard out")
	cmd.Flags().StringVarP(&migrationsPath, "migrations-path", "m", "pkg/store/rw/migrations", "Migrations are written to mark changes as merged")

	return cmd
}

func mergeTags(outputDataPath string, migrationsPath string, conn *rw.Conn, dryRun bool, logger *zap.Logger) error {

	logger.Info("Reading DB...")

	ctx := context.Background()

	mergedTags := models.Tags{}

	// IDs are a composite key of username + episode ID
	err := conn.WithStore(func(s *rw.Store) error {

		tags, err := s.ListTranscriptTags(
			ctx,
		)
		if err != nil {
			return err
		}

		for _, v := range tags {

			logger.Info(fmt.Sprintf("Processing rating %s/%s...", v.Name, v.Timestamp.String()))

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

			episodeOnDisk.Tags = append(episodeOnDisk.Tags, v)

			if dryRun {
				if err := stdoutPrinter.Encode(v); err != nil {
					return err
				}
			} else {
				if err := data.ReplaceEpisodeFile(outputDataPath, episodeOnDisk); err != nil {
					return err
				}
				mergedTags = append(mergedTags, v)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(mergedTags) > 0 {
		return generateTagMigrations(migrationsPath, mergedTags)
	}

	return nil
}

func generateTagMigrations(migrationsPath string, tags models.Tags) error {
	idsPairs := make([]string, len(tags))
	for k, v := range tags {
		idsPairs[k] = fmt.Sprintf(`('%s', '%s', '%d')`, v.EpisodeID, v.Name, v.Timestamp)
	}
	return os.WriteFile(
		path.Join(migrationsPath, fmt.Sprintf("%d_merge_tags.sql", time.Now().Unix())),
		[]byte(fmt.Sprintf(
			`DELETE FROM transcript_tag WHERE (episode_id, tag_name, tag_timestamp) IN (%s);`,
			strings.Join(idsPairs, ", "),
		)),
		0666,
	)
}
