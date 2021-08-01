package db

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"math/rand"
)

func CreateRwTestdataCmd() *cobra.Command {

	dbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "create-rw-testdata",
		Short: "Load some random authors and contributions",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			if dbCfg.DSN == "" {
				panic("dsn not set")
			}
			conn, err := rw.NewConn(dbCfg)
			if err != nil {
				return err
			}
			if err := conn.Migrate(); err != nil {
				return err
			}
			return addTestData(conn, logger)
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")

	return cmd
}

func addTestData(conn *rw.Conn, logger *zap.Logger) error {

	logger.Info("Populating DB...")

	for i := 0; i <= 5; i++ {

		ath := &models.Author{
			Name:     fmt.Sprintf("author-%d", i),
			Identity: "{}",
			Banned:   false,
			Approver: false,
		}

		// add author
		err := conn.WithStore(func(s *rw.Store) error {
			if err := s.UpsertAuthor(context.Background(), ath); err != nil {
				return err
			}

			logger.Info(fmt.Sprintf("Added author %s..", ath.Name))

			for ii := 0; ii < rand.Intn(10); ii++ {
				ch := randomChunk(s)
				_, err := s.CreateChunkContribution(context.Background(), &models.ContributionCreate{
					AuthorID:      ath.ID,
					ChunkID:       ch.ID,
					Transcription: ch.Raw,
					State:         randomState(),
				})
				if err != nil {
					return err
				}
				logger.Info(fmt.Sprintf("Added contribution for chunk %s..", ch.ID))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

var states = []models.ContributionState{
	models.ContributionStatePending,
	models.ContributionStateApproved,
	models.ContributionStateApprovalRequested,
	models.ContributionStateRejected,
}

func randomState() models.ContributionState {
	return states[rand.Intn(len(states)-1)]
}

var chunks []*models.Chunk

func randomChunk(s *rw.Store) *models.Chunk {
	if chunks == nil {
		ch, err := s.ListChunks(context.Background(), common.Q(common.WithPaging(25, 0)))
		if err != nil {
			panic(err)
		}
		chunks = ch
	}
	return chunks[rand.Intn(len(chunks)-1)]
}
