package rw

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
)

type ChunkActivity string

const (
	ChunkActivityAccessed  = "accessed"  // chunk fetched
	ChunkActivitySubmitted = "submitted" // contribution submitted
	ChunkActivityApproved  = "approved"  // contribution approved
)

//go:embed migrations
var migrations embed.FS

func NewConn(cfg *common.Config) (*Conn, error) {
	innerConn, err := common.NewConn("postgres", cfg)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: innerConn}, nil
}

type Conn struct {
	*common.Conn
}

func (c *Conn) Migrate() error {
	return c.Conn.Migrate(migrations)
}

func (c *Conn) WithStore(f func(s *Store) error) error {
	return c.WithTx(func(tx *sqlx.Tx) error {
		return f(&Store{tx: tx})
	})
}

type Store struct {
	tx *sqlx.Tx
}

func (s *Store) InsertTscript(ctx context.Context, tscript *models.Tscript) error {

	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO tscript (id, publication, series, episode) VALUES ($1, $2, $3, $4)`,
		tscript.ID(),
		tscript.Publication,
		tscript.Series,
		tscript.Episode,
	)
	if err != nil {
		return err
	}
	for _, v := range tscript.Chunks {
		if err != nil {
			return err
		}
		_, err = s.tx.ExecContext(ctx,
			`INSERT INTO tscript_chunk (id, tscript_id, raw, start_second, end_second) VALUES ($1, $2, $3, $4, $5)`,
			v.ID,
			tscript.ID(),
			v.Raw,
			v.StartSecond,
			v.EndSecond,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Store) GetChunk(ctx context.Context, chunkId string) (*models.Chunk, string, error) {
	ch := &models.Chunk{}
	var tscriptID string

	err := s.tx.
		QueryRowxContext(ctx, "SELECT id, tscript_id, raw, start_second, end_second FROM tscript_chunk WHERE id = $1", chunkId).
		Scan(&ch.ID, &tscriptID, &ch.Raw, &ch.StartSecond, &ch.EndSecond)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}
	return ch, tscriptID, s.UpdateChunkActivity(ctx, ch.ID, ChunkActivityAccessed)
}

func (s *Store) CreateContribution(ctx context.Context, c *models.Contribution) error {
	if c.ID == "" {
		c.ID = shortuuid.New()
	}
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO tscript_contribution (id, author_id, tscript_chunk_id, transcription) VALUES ($1, $2, $3, $4)`,
		c.ID,
		c.AuthorID,
		c.ChunkID,
		c.Transcription,
	)
	if err != nil {
		return err
	}
	return s.UpdateChunkActivity(ctx, c.ChunkID, ChunkActivitySubmitted)
}

func (s *Store) UpdateChunkActivity(ctx context.Context, id string, activity ChunkActivity) error {
	var col string
	switch activity {
	case ChunkActivityAccessed:
		col = "accessed_at"
	case ChunkActivitySubmitted:
		col = "submitted"
	case ChunkActivityApproved:
		col = "approved_at"
	default:
		return fmt.Errorf("unknown activity %s", activity)
	}
	_, err := s.tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO tscript_chunk_activity (tscript_chunk_id, %s) VALUES ($1, NOW()) ON CONFLICT(tscript_chunk_id) DO UPDATE SET %s=NOW() ", col, col), id)
	return err
}

func (s *Store) GetChunkStats(ctx context.Context) (*models.ChunkStats, error) {
	ch := &models.ChunkStats{}

	query := `
		SELECT 
			c.id as next_chunk, agg.*
		FROM tscript_chunk c 
		LEFT JOIN tscript_chunk_activity a ON c.id = a.tscript_chunk_id
		JOIN (
			SELECT 
				SUM(1) as total_chunks, 
				SUM(CASE WHEN aa.approved_at IS NOT NULL then 1 ELSE 0 END) as approved_chunks,
				SUM(CASE WHEN aa.submitted_at IS NOT NULL then 1 ELSE 0 END) as submitted_chunks
			FROM tscript_chunk cc 
			LEFT JOIN tscript_chunk_activity aa ON cc.id = aa.tscript_chunk_id
		) agg ON true
		WHERE a.approved_at IS NULL
		ORDER BY a.accessed_at ASC NULLS FIRST , a.submitted_at ASC NULLS FIRST 
		LIMIT 1
	`

	err := s.tx.QueryRowxContext(ctx, query).Scan(&ch.NextChunk, &ch.TotalChunks, &ch.ApprovedChunks, &ch.SubmittedChunks)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return ch, nil
}
