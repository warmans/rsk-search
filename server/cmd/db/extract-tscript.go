package db

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/transcript"
	"go.uber.org/zap"
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

			logger.Info("Connecting...")
			conn, err := rw.NewConn(dbCfg)
			if err != nil {
				return err
			}
			defer func(conn *rw.Conn) {
				_ = conn.Close()
			}(conn)

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
					OffsetAccuracy: 0,
					Contributors:   []string{},
				}
			}

			// the transcript is rebuilt from scratch but synopsis/trivia are
			// merged by position so entries that only exist on disk are kept
			episodeOnDisk.Transcript = nil

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

			currentPos := int64(0)
			for _, ch := range allChunks {

				var chContribution *models.ChunkContribution
				for _, co := range approved {
					if co.ChunkID == ch.ID {
						chContribution = co
					}
				}

				// all chunks need to be processed.
				if len(episodeOnDisk.Transcript) > 0 {
					currentPos = episodeOnDisk.Transcript[len(episodeOnDisk.Transcript)-1].Position
				}

				// if the transcript is missing insert a placeholder
				if chContribution == nil {
					ts, err := transcript.Import(bufio.NewScanner(bytes.NewBufferString(ch.Raw)), episodeOnDisk.ID(), currentPos+transcript.PosSpacing)
					if err == nil {
						for _, v := range ts.Transcript {
							v.Placeholder = true
							episodeOnDisk.Transcript = append(episodeOnDisk.Transcript, v)
						}
					} else {
						logger.Error("Failed to parse placeholder transcript", zap.Error(err))
						episodeOnDisk.Transcript = append(
							episodeOnDisk.Transcript,
							models.Dialog{
								ID:          models.DialogID(episodeOnDisk.ID(), currentPos+transcript.PosSpacing),
								Position:    currentPos + transcript.PosSpacing,
								Timestamp:   0,
								Type:        "gap",
								Actor:       "",
								Meta:        nil,
								Content:     "[~3 mins of missing transcription]",
								Placeholder: true,
							},
						)
					}
					episodeOnDisk.Completion = models.CompletionStateGaps
					continue
				}

				ts, err := transcript.Import(bufio.NewScanner(bytes.NewBufferString(chContribution.Transcription)), episodeOnDisk.ID(), currentPos+transcript.PosSpacing)
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
				episodeOnDisk.Transcript = append(episodeOnDisk.Transcript, ts.Transcript...)
				episodeOnDisk.Synopsis = upsertSynopsis(episodeOnDisk.Synopsis, ts.Synopsis)
				episodeOnDisk.Trivia = upsertTrivia(episodeOnDisk.Trivia, ts.Trivia)
			}

			// keep entries in position order; stable so on-disk entries without
			// a DB counterpart keep their relative order
			sort.SliceStable(episodeOnDisk.Synopsis, func(i, j int) bool {
				return episodeOnDisk.Synopsis[i].StartPos < episodeOnDisk.Synopsis[j].StartPos
			})
			sort.SliceStable(episodeOnDisk.Trivia, func(i, j int) bool {
				return episodeOnDisk.Trivia[i].StartPos < episodeOnDisk.Trivia[j].StartPos
			})

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

// upsertSynopsis adds incoming entries to existing, replacing any entry that
// covers the same start/end positions. Entries already in existing with no
// incoming counterpart are left untouched.
func upsertSynopsis(existing []models.Synopsis, incoming []models.Synopsis) []models.Synopsis {
	for _, in := range incoming {
		replaced := false
		for i := range existing {
			if existing[i].StartPos == in.StartPos && existing[i].EndPos == in.EndPos {
				existing[i] = in
				replaced = true
				break
			}
		}
		if !replaced {
			existing = append(existing, in)
		}
	}
	return existing
}

// upsertTrivia adds incoming entries to existing, replacing any entry that
// covers the same start/end positions. Entries already in existing with no
// incoming counterpart are left untouched.
func upsertTrivia(existing []models.Trivia, incoming []models.Trivia) []models.Trivia {
	for _, in := range incoming {
		replaced := false
		for i := range existing {
			if existing[i].StartPos == in.StartPos && existing[i].EndPos == in.EndPos {
				existing[i] = in
				replaced = true
				break
			}
		}
		if !replaced {
			existing = append(existing, in)
		}
	}
	return existing
}
