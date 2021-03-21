package rw

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
)

type ChunkActivity string

const (
	ChunkActivityAccessed  = "accessed"  // chunk fetched
	ChunkActivitySubmitted = "submitted" // contribution submitted
	ChunkActivityApproved  = "approved"  // contribution approved
	ChunkActivityRejected  = "rejected"  // contribution rejected
)

var ErrNotPermitted = errors.New("user not allowed to perform action")

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

func (s *Store) GetChunkContributionCount(ctx context.Context, chunkId string) (int32, error) {
	var count int32
	err := s.tx.
		QueryRowxContext(ctx, "SELECT COUNT(*) FROM tscript_contribution c LEFT JOIN author a ON c.author_id = a.id WHERE a.banned = false AND tscript_chunk_id = $1", chunkId).
		Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return count, nil
		}
		return 0, err
	}
	return count, nil
}

func (s *Store) GetAuthorStats(ctx context.Context, authorID string) (*models.AuthorStats, error) {

	query := `
		SELECT 
			SUM(CASE WHEN c.state = 'pending' THEN 1 ELSE 0 END) as num_pending,
			SUM(CASE WHEN c.state = 'approved' THEN 1 ELSE 0 END) as num_approved,
			SUM(CASE WHEN c.state = 'rejected' THEN 1 ELSE 0 END) as num_rejected,
			SUM(CASE WHEN c.created_at > NOW() - INTERVAL '1 HOUR' THEN 1 ELSE 0 END) as total_in_last_hour
		FROM tscript_contribution c
		WHERE author_id = $1
	`
	stats := &models.AuthorStats{}
	err := s.tx.
		QueryRowxContext(ctx, query, authorID).
		Scan(&stats.PendingContributions, &stats.ApprovedContributions, &stats.RejectedContributions, &stats.ContributionsInLastHour)

	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *Store) CreateContribution(ctx context.Context, c *models.Contribution) error {
	if c.ID == "" {
		c.ID = shortuuid.New()
	}
	if banned, err := s.AuthorIsBanned(ctx, c.AuthorID); err != nil || banned {
		if err != nil {
			return err
		}
		return ErrNotPermitted
	}
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO tscript_contribution (id, author_id, tscript_chunk_id, transcription, state, created_at) VALUES ($1, $2, $3, $4, $5, NOW())`,
		c.ID,
		c.AuthorID,
		c.ChunkID,
		c.Transcription,
		models.ContributionStatePending,
	)
	if err != nil {
		return err
	}
	return s.UpdateChunkActivity(ctx, c.ChunkID, ChunkActivitySubmitted)
}

func (s *Store) UpdateContribution(ctx context.Context, c *models.Contribution) error {
	if c.ID == "" {
		return fmt.Errorf("no identifier was provided")
	}
	if banned, err := s.AuthorIsBanned(ctx, c.AuthorID); err != nil || banned {
		if err != nil {
			return errors.Wrap(err, "failed to identity author")
		}
		return ErrNotPermitted
	}
	_, err := s.tx.ExecContext(
		ctx,
		`UPDATE tscript_contribution SET transcription=$1 WHERE id=$2`,
		c.ID,
		c.Transcription,
	)
	return err
}

func (s *Store) GetContribution(ctx context.Context, id string) (*models.Contribution, error) {
	out := &models.Contribution{}
	row := s.tx.QueryRowxContext(ctx, `SELECT id, author_id, tscript_chunk_id, transcription, COALESCE(state, 'unknown') FROM tscript_contribution WHERE id=$1`, id)
	if err := row.Scan(&out.ID, &out.AuthorID, &out.ChunkID, &out.Transcription, &out.State); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) ListAuthorContributions(ctx context.Context, authorName string, page int32) ([]*models.Contribution, error) {
	out := make([]*models.Contribution, 0)

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`SELECT id, author_id, tscript_chunk_id, transcription, COALESCE(state, 'unknown') FROM tscript_contribution WHERE author_id = $1 LIMIT 25 OFFSET %d`, page),
		authorName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cur := &models.Contribution{}
		if err := rows.Scan(&cur.ID, &cur.AuthorID, &cur.ChunkID, &cur.Transcription, &cur.State); err != nil {
			return nil, err
		}
		out = append(out, cur)
	}
	return out, nil
}

func (s *Store) UpdateChunkActivity(ctx context.Context, id string, activity ChunkActivity) error {
	var col string
	switch activity {
	case ChunkActivityAccessed:
		col = "accessed_at"
	case ChunkActivitySubmitted:
		col = "submitted_at"
	case ChunkActivityApproved:
		col = "approved_at"
	case ChunkActivityRejected:
		col = "rejected_at"
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
		LEFT JOIN (
			SELECT
				tscript_chunk_id AS chunk_id,
				SUM(1) as total_submitted
			FROM tscript_contribution 
			GROUP BY tscript_chunk_id
		) stats ON c.id = stats.chunk_id
		WHERE a.approved_at IS NULL
		ORDER BY stats.total_submitted DESC NULLS FIRST, a.accessed_at ASC NULLS FIRST , a.submitted_at ASC NULLS FIRST 
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

func (s *Store) UpsertAuthor(ctx context.Context, author *models.Author) error {
	if author.ID == "" {
		author.ID = shortuuid.New()
	}
	if author.Name == "" {
		return fmt.Errorf("author name cannot be empty")
	}
	row := s.tx.QueryRowxContext(
		ctx,
		"INSERT INTO author (id, name, identity, created_at) VALUES ($1, $2, $3, NOW()) ON CONFLICT(name) DO UPDATE SET identity=$3 RETURNING id, banned, approver",
		author.ID,
		author.Name,
		author.Identity,
	)
	return row.Scan(&author.ID, &author.Banned, &author.Approver)
}

func (s *Store) AuthorIsBanned(ctx context.Context, id string) (bool, error) {
	var banned bool
	row := s.tx.QueryRowxContext(ctx, "SELECT banned FROM author WHERE id=$1 ", id)
	return banned, row.Scan(&banned)
}
