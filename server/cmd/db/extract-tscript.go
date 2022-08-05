package db

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/transcript"
	"go.uber.org/zap"
	"os"
	"sort"
)

var stdoutPrinter = json.NewEncoder(os.Stdout)

func init() {
	stdoutPrinter.SetIndent("  ", "  ")
}

func ExtractTscriptCmd() *cobra.Command {
	var outputDir string
	var dryRun bool
	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "extract-tscript",
		Short: "Extract the completed/partial transcripts and store them as regular transcripts.",
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

			return extract(outputDir, conn, dryRun, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "var/data/episodes", "Path to output the data")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", true, "Instead of saving the data print it out standard out")

	return cmd
}

func extract(outputDataPath string, conn *rw.Conn, dryRun bool, logger *zap.Logger) error {

	logger.Info("Reading DB...")

	ctx := context.Background()

	err := conn.WithStore(func(s *rw.Store) error {

		allTscripts, err := s.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, v := range allTscripts {

			logger.Info(fmt.Sprintf("Processing tscript %s-%s...", v.Publication, models.FormatStandardEpisodeName(v.Series, v.Episode)))

			episodeOnDisk, err := data.LoadEpisodeByName(outputDataPath, v.Publication, models.FormatStandardEpisodeName(v.Series, v.Episode))
			if err != nil {
				return err
			}
			if episodeOnDisk == nil {
				episodeOnDisk = &models.Transcript{
					Publication:    v.Publication,
					Series:         v.Series,
					Episode:        v.Episode,
					Name:           v.Name,
					Incomplete:     true,
					OffsetAccuracy: 0,
					Contributors:   []string{},
				}
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

				// all chunks need to be processed.
				currentPos := int64(1)
				if len(episodeOnDisk.Transcript) > 0 {
					currentPos = episodeOnDisk.Transcript[len(episodeOnDisk.Transcript)-1].Position
				}

				// if the transcript is missing insert a placeholder
				if chContribution == nil {
					currentPos += transcript.PosSpacing
					episodeOnDisk.Transcript = append(
						episodeOnDisk.Transcript,
						models.Dialog{
							ID:        models.DialogID(episodeOnDisk.ID(), currentPos),
							Position:  currentPos,
							OffsetSec: 0,
							Type:      "gap",
							Actor:     "",
							Meta:      nil,
							Content:   "[~3 mins of missing transcription]",
						},
					)
					episodeOnDisk.Incomplete = true
					continue
				}

				dialog, synopsis, trivia, err := transcript.Import(bufio.NewScanner(bytes.NewBufferString(chContribution.Transcription)), episodeOnDisk.ID(), currentPos)
				if err != nil {
					return err
				}

				// process contributors for this chunk of audio
				author, err := s.GetAuthor(context.Background(), chContribution.Author.ID)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to get author with ID %s", chContribution.Author.ID))
					continue
				} else {
					uniqueContributors[author.Name] = struct{}{}
				}
				episodeOnDisk.Transcript = append(episodeOnDisk.Transcript, dialog...)
				episodeOnDisk.Synopsis = append(episodeOnDisk.Synopsis, synopsis...)
				episodeOnDisk.Trivia = append(episodeOnDisk.Trivia, trivia...)
			}

			contributors := []string{}
			for c := range uniqueContributors {
				contributors = append(contributors, c)
			}
			sort.Strings(contributors)
			episodeOnDisk.Contributors = contributors

			if dryRun {
				if err := stdoutPrinter.Encode(episodeOnDisk); err != nil {
					return err
				}
			} else {
				if err := data.ReplaceEpisodeFile(outputDataPath, episodeOnDisk); err != nil {
					return errors.Wrap(err, "failed to write episode to disk")
				}
			}
		}
		return nil
	})
	return err
}
