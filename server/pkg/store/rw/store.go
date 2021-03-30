package rw

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
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

func (s *Store) ListTscripts(ctx context.Context) ([]*models.TscriptStats, error) {
	out := make([]*models.TscriptStats, 0)

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
			SELECT 
				ts.id,
				ts.publication, 
				ts.series,
				ts.episode,
				json_object_agg(ch.id, contribution_states.states) AS contribution_states,
 				COUNT(DISTINCT ch.id) num_chunks,
 				COUNT(DISTINCT co.id) num_contributions,
 				SUM(CASE WHEN a.banned = false AND co.state = 'approved' THEN 1 ELSE 0 END) num_approved_contributions,
 				SUM(CASE WHEN a.banned = false AND co.state = 'pending' THEN 1 ELSE 0 END) num_pending_contributions,
 				SUM(CASE WHEN a.banned = false AND co.state = 'request_approval' THEN 1 ELSE 0 END) num_request_approval_contributions
			FROM tscript ts
			LEFT JOIN tscript_chunk ch ON ts.id = ch.tscript_id
			LEFT JOIN tscript_contribution co ON ch.id = co.tscript_chunk_id
			LEFT JOIN author a ON co.author_id = a.id
			LEFT JOIN (
                SELECT tscript_chunk_id, json_agg(DISTINCT state) AS states 
                FROM tscript_contribution 
                LEFT JOIN tscript_chunk ON tscript_contribution.tscript_chunk_id = tscript_chunk.id 
                GROUP BY tscript_chunk_id) as contribution_states ON ch.id = contribution_states.tscript_chunk_id
			GROUP BY ts.id
			ORDER BY ts.publication, ts.series, ts.episode ASC
		`),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cur := &models.TscriptStats{
			ChunkContributionStates: map[string][]models.ContributionState{},
		}
		var contribStates string

		if err := rows.Scan(
			&cur.ID,
			&cur.Publication,
			&cur.Series,
			&cur.Episode,
			&contribStates,
			&cur.NumChunks,
			&cur.NumContributions,
			&cur.NumApprovedContributions,
			&cur.NumPendingContributions,
			&cur.NumRequestApprovalContributions,

		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(contribStates), &cur.ChunkContributionStates); err != nil {
			return nil, err
		}
		out = append(out, cur)
	}
	return out, nil
}

func (s *Store) InsertOrIgnoreTscript(ctx context.Context, tscript *models.Tscript) error {

	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO tscript (id, publication, series, episode) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`,
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
			`INSERT INTO tscript_chunk (id, tscript_id, raw, start_second, end_second) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`,
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

func (s *Store) ListChunks(ctx context.Context, limit int32) ([]*models.Chunk, error) {

	rows, err := s.tx.QueryxContext(ctx, fmt.Sprintf(`SELECT id, raw, start_second, end_second FROM tscript_chunk LIMIT %d`, limit))
	if err != nil {
		return nil, err
	}
	chunks := []*models.Chunk{}
	for rows.Next() {
		ch := &models.Chunk{}
		if err := rows.Scan(&ch.ID, &ch.Raw, &ch.StartSecond, &ch.EndSecond); err != nil {
			return nil, err
		}
		chunks = append(chunks, ch)
	}
	return chunks, nil
}

func (s *Store) ListTscriptChunks(ctx context.Context, tscriptID string) ([]*models.Chunk, error) {
	rows, err := s.tx.QueryxContext(ctx, fmt.Sprintf(`SELECT id, raw, start_second, end_second FROM tscript_chunk WHERE tscript_id=$1`), tscriptID)
	if err != nil {
		return nil, err
	}
	chunks := []*models.Chunk{}
	for rows.Next() {
		ch := &models.Chunk{}
		if err := rows.Scan(&ch.ID, &ch.Raw, &ch.StartSecond, &ch.EndSecond); err != nil {
			return nil, err
		}
		chunks = append(chunks, ch)
	}
	return chunks, nil
}

func (s *Store) GetChunkContributionCount(ctx context.Context, chunkId string) (int32, error) {
	var count int32
	err := s.tx.
		QueryRowxContext(ctx, "SELECT COUNT(*) FROM tscript_contribution c LEFT JOIN author a ON c.author_id = a.id WHERE a.banned = false AND tscript_chunk_id = $1 AND c.state NOT IN ('pending', 'rejected')", chunkId).
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
			COALESCE(SUM(CASE WHEN c.state = 'pending' THEN 1 ELSE 0 END), 0) as num_pending,
			COALESCE(SUM(CASE WHEN c.state = 'request_approval' THEN 1 ELSE 0 END), 0) as num_request_approval,
			COALESCE(SUM(CASE WHEN c.state = 'approved' THEN 1 ELSE 0 END), 0) as num_approved,
			COALESCE(SUM(CASE WHEN c.state = 'rejected' THEN 1 ELSE 0 END), 0) as num_rejected,
			COALESCE(SUM(CASE WHEN c.created_at > NOW() - INTERVAL '1 HOUR' THEN 1 ELSE 0 END), 0) as total_in_last_hour
		FROM tscript_contribution c
		WHERE author_id = $1
	`
	stats := &models.AuthorStats{}
	err := s.tx.
		QueryRowxContext(ctx, query, authorID).
		Scan(
			&stats.PendingContributions,
			&stats.RequestApprovalContributions,
			&stats.ApprovedContributions,
			&stats.RejectedContributions,
			&stats.ContributionsInLastHour,
		)

	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *Store) CreateContribution(ctx context.Context, c *models.Contribution) error {
	if c.ID == "" {
		c.ID = shortuuid.New()
	}
	if c.State == "" {
		c.State = models.ContributionStatePending
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
		c.State,
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
		`UPDATE tscript_contribution SET transcription=$1, state=$2 WHERE id=$3`,
		c.Transcription,
		c.State,
		c.ID,
	)
	return err
}

func (s *Store) UpdateContributionState(ctx context.Context, id string, state models.ContributionState) error {
	_, err := s.tx.ExecContext(
		ctx,
		`UPDATE tscript_contribution SET state=$1 WHERE id=$2`,
		state,
		id,
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

func (s *Store) ListNonPendingTscriptContributions(ctx context.Context, tscriptID string, page int32) ([]*models.Contribution, error) {

	out := make([]*models.Contribution, 0)

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
			SELECT 
				COALESCE(co.id, ''), 
				COALESCE(co.author_id, ''), 
				ch.id, 
				COALESCE(co.transcription, ''), 
				COALESCE(co.state, 'unknown') 
			FROM tscript_chunk ch 
			LEFT JOIN tscript_contribution co ON ch.id = co.tscript_chunk_id AND co.state != $1
			WHERE ch.tscript_id = $2
			ORDER BY ch.start_second ASC
			LIMIT 25 OFFSET %d`, page),
		models.ContributionStatePending,
		tscriptID,
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

func (s *Store) AuthorLeaderboard(ctx context.Context) (*models.AuthorLeaderboard, error) {

	query := `
        SELECT * FROM (
            SELECT 
                a.name,
                a.approver,
                COALESCE(SUM(CASE WHEN c.state = 'approved' THEN 1 ELSE 0 END), 0) as num_approved
            FROM author a
            LEFT JOIN tscript_contribution c ON c.author_id = a.id
            GROUP BY a.name, a.approver) ranks
		WHERE ranks.num_approved > 0
		ORDER BY ranks.num_approved DESC
		LIMIT 25
	`
	rows, err := s.tx.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := []*models.AuthorRanking{}
	for rows.Next() {
		author := &models.AuthorRanking{}
		err := rows.Scan(
			&author.Name,
			&author.Approver,
			&author.AcceptedContributions,
		)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return &models.AuthorLeaderboard{Authors: authors}, nil
}

func (s *Store) ListAuthorContributions(ctx context.Context, authorName string, page int32) ([]*models.Contribution, error) {
	return s.listContributions(
		ctx,
		50,
		page,
		"co.author_id = $1",
		[]interface{}{authorName},
		"co.created_at ASC",
	)
}

func (s *Store) ListApprovedTscriptContributions(ctx context.Context, tscriptID string, numPerPage int32, page int32) ([]*models.Contribution, error) {
	return s.listContributions(
		ctx,
		numPerPage,
		page,
		"co.state = $1 and ch.tscript_id = $2",
		[]interface{}{models.ContributionStateApproved, tscriptID},
		"ch.start_second ASC",
	)
}

func (s *Store) listContributions(ctx context.Context, numPerPage int32, page int32, where string, params []interface{}, order string) ([]*models.Contribution, error) {
	out := make([]*models.Contribution, 0)

	if where != "" {
		where = fmt.Sprintf("WHERE %s", where)
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
			SELECT 
				COALESCE(co.id, ''), 
				COALESCE(co.author_id, ''), 
				ch.id, 
				COALESCE(co.transcription, ''), 
				COALESCE(co.state, 'unknown') 
			FROM tscript_chunk ch 
			LEFT JOIN tscript_contribution co ON ch.id = co.tscript_chunk_id
			%s 
			ORDER BY %s
			LIMIT %d OFFSET %d`, where, order, numPerPage, page),
		params...,
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
				COALESCE(SUM(CASE WHEN aa.approved_at IS NOT NULL then 1 ELSE 0 END), 0) as approved_chunks,
				COALESCE(SUM(CASE WHEN aa.submitted_at IS NOT NULL then 1 ELSE 0 END), 0) as submitted_chunks
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
