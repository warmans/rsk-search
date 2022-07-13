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
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/util"
	"math"
	"strings"
	"time"
)

type ChunkActivity string

const ChunkContributionPoints = 3

const (
	ChunkActivityAccessed  = "accessed"  // chunk fetched
	ChunkActivitySubmitted = "submitted" // contribution submitted
	ChunkActivityApproved  = "approved"  // contribution approved
	ChunkActivityRejected  = "rejected"  // contribution rejected
)

func ActivityFromState(state models.ContributionState) ChunkActivity {
	switch state {
	case models.ContributionStateApproved:
		return ChunkActivityApproved
	case models.ContributionStateRejected:
		return ChunkActivityRejected
	case models.ContributionStateApprovalRequested:
		return ChunkActivitySubmitted
	}
	return ChunkActivityAccessed
}

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
		`
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
		`,
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

func (s *Store) GetChunk(ctx context.Context, chunkId string) (*models.Chunk, error) {
	ch, err := s.ListChunks(ctx, &common.QueryModifier{
		Filter: filter.Eq("id", filter.String(chunkId)),
	})
	if err != nil {
		return nil, err
	}
	if len(ch) == 0 {
		return nil, sql.ErrNoRows
	}
	return ch[0], s.UpdateChunkActivity(ctx, ch[0].ID, ChunkActivityAccessed)
}

func (s *Store) ListChunks(ctx context.Context, q *common.QueryModifier) ([]*models.Chunk, error) {

	fieldMap := map[string]string{
		"id":           "ch.id",
		"tscript_id":   "ch.tscript_id",
		"start_second": "ch.start_second",
		"end_second":   "ch.end_second",
	}

	where, params, order, paging, err := q.ToSQL(fieldMap, true)
	if err != nil {
		return nil, err
	}
	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
			SELECT 
				ch.id, 
				ch.tscript_id,
				ch.raw,
				ch.start_second,
				ch.end_second,
				(SELECT COUNT(*) FROM tscript_contribution WHERE tscript_chunk_id  = ch.id AND state != 'pending' AND state != 'rejected') AS num_contributions
			FROM tscript_chunk ch
			%s 
			%s 
			%s`,
			where,
			order,
			paging,
		),
		params...,
	)
	if err != nil {
		return nil, err
	}
	chunks := make([]*models.Chunk, 0)
	for rows.Next() {
		ch := &models.Chunk{}
		if err := rows.StructScan(ch); err != nil {
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

func (s *Store) CreateChunkContribution(ctx context.Context, c *models.ContributionCreate) (*models.ChunkContribution, error) {
	if c.State == "" {
		c.State = models.ContributionStatePending
	}
	if banned, err := s.AuthorIsBanned(ctx, c.AuthorID); err != nil || banned {
		if err != nil {
			return nil, err
		}
		return nil, ErrNotPermitted
	}
	contribution := &models.ChunkContribution{
		ID: shortuuid.New(),
		Author: &models.ShortAuthor{
			ID: c.AuthorID,
		},
		ChunkID:       c.ChunkID,
		Transcription: c.Transcription,
		State:         c.State,
	}
	row := s.tx.QueryRowxContext(
		ctx,
		`
		INSERT INTO tscript_contribution (id, author_id, tscript_chunk_id, transcription, state, created_at) VALUES ($1, $2, $3, $4, $5, NOW())
	 	RETURNING (SELECT name FROM author WHERE id=$2) AS author_name, created_at
		`,
		contribution.ID,
		contribution.Author.ID,
		contribution.ChunkID,
		contribution.Transcription,
		contribution.State,
	)
	if err := row.Scan(&contribution.Author.Name, &contribution.CreatedAt); err != nil {
		return nil, err
	}
	return contribution, s.UpdateChunkActivity(ctx, c.ChunkID, ChunkActivitySubmitted)
}

func (s *Store) UpdateChunkContribution(ctx context.Context, c *models.ContributionUpdate) error {
	if c.ID == "" {
		return fmt.Errorf("no identifier was provided")
	}

	oldCon, err := s.GetChunkContribution(ctx, c.ID)
	if err != nil {
		return err
	}
	_, err = s.tx.ExecContext(
		ctx,
		`UPDATE tscript_contribution SET transcription=$1, state=$2 WHERE id=$3`,
		c.Transcription,
		c.State,
		c.ID,
	)
	if oldCon.State != models.ContributionStateApproved && c.State == models.ContributionStateApproved {
		contributionID, err := s.CreateAuthorContribution(ctx, models.AuthorContributionCreate{
			AuthorID:         oldCon.Author.ID,
			EpID:             strings.Replace(oldCon.TscriptID, "ts-", "ep-", 1),
			ContributionType: models.ContributionTypeChunk,
			Points:           ChunkContributionPoints,
		})
		if err != nil {
			return err
		}
		if err := s.LinkAuthorContributionToTscriptContribution(ctx, contributionID, c.ID); err != nil {
			return err
		}
	}
	return err
}

func (s *Store) UpdateChunkContributionState(ctx context.Context, id string, state models.ContributionState, comment string) error {

	con, err := s.GetChunkContribution(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.tx.ExecContext(
		ctx,
		`UPDATE tscript_contribution SET state=$1, state_comment=NULLIF($2, '') WHERE id=$3`,
		state,
		comment,
		id,
	)
	if err != nil {
		return err
	}
	if con.State != models.ContributionStateApproved && state == models.ContributionStateApproved {
		contributionID, err := s.CreateAuthorContribution(ctx, models.AuthorContributionCreate{
			AuthorID:         con.Author.ID,
			EpID:             strings.Replace(con.TscriptID, "ts-", "ep-", 1),
			ContributionType: models.ContributionTypeChunk,
			Points:           ChunkContributionPoints,
		})
		if err != nil {
			return err
		}
		if err := s.LinkAuthorContributionToTscriptContribution(ctx, contributionID, con.ID); err != nil {
			return err
		}
	}
	return err
}

func (s *Store) ListChunkContributions(ctx context.Context, q *common.QueryModifier) ([]*models.ChunkContribution, error) {

	fieldMap := map[string]string{
		"id":            "c.id",
		"tscript_id":    "ch.tscript_id",
		"chunk_id":      "ch.id",
		"author_id":     "c.author_id",
		"author_name":   "a.name",
		"transcription": "c.transcription",
		"state":         "c.state",
		"created_at":    "c.created_at",
	}

	where, params, order, paging, err := q.ToSQL(fieldMap, false)
	if err != nil {
		return nil, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
		SELECT 
			c.id,
			ch.tscript_id,
			c.tscript_chunk_id,
       		c.author_id, 
       		COALESCE(a.name, 'unknown'), 
       		c.transcription, 
       		COALESCE(c.state, 'unknown'),
       		COALESCE(c.state_comment, ''),
       		c.created_at
		FROM tscript_contribution c
		LEFT JOIN tscript_chunk ch ON c.tscript_chunk_id = ch.id 
		LEFT JOIN author a ON c.author_id = a.id
		WHERE a.banned = false 
		AND %s
		%s
		%s
		`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*models.ChunkContribution, 0)
	for rows.Next() {
		cur := &models.ChunkContribution{Author: &models.ShortAuthor{}}
		if err := rows.Scan(
			&cur.ID,
			&cur.TscriptID,
			&cur.ChunkID,
			&cur.Author.ID,
			&cur.Author.Name,
			&cur.Transcription,
			&cur.State,
			&cur.StateComment,
			&cur.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, cur)
	}
	return out, nil
}

func (s *Store) GetChunkContribution(ctx context.Context, id string) (*models.ChunkContribution, error) {
	contributions, err := s.ListChunkContributions(ctx, &common.QueryModifier{Filter: filter.Eq("id", filter.String(id))})
	if err != nil {
		return nil, err
	}
	if len(contributions) == 0 {
		return nil, sql.ErrNoRows
	}
	return contributions[0], nil
}

func (s *Store) DeleteContribution(ctx context.Context, id string) error {
	_, err := s.tx.ExecContext(
		ctx,
		`DELETE FROM tscript_contribution WHERE id=$1`,
		id,
	)
	return err
}

func (s *Store) ListNonPendingTscriptContributions(ctx context.Context, tscriptID string, page int32) ([]*models.ChunkContribution, error) {

	out := make([]*models.ChunkContribution, 0)

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
			SELECT 
				COALESCE(co.id, ''), 
				COALESCE(co.author_id, ''), 
				COALESCE(a.name, ''),
				ch.id, 
				COALESCE(co.transcription, ''), 
				COALESCE(co.state, 'unknown') 
			FROM tscript_chunk ch 
			LEFT JOIN tscript_contribution co ON ch.id = co.tscript_chunk_id AND co.state != $1
			LEFT JOIN author a ON co.author_id = a.id
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
		cur := &models.ChunkContribution{Author: &models.ShortAuthor{}}
		if err := rows.Scan(&cur.ID, &cur.Author.ID, &cur.Author.Name, &cur.ChunkID, &cur.Transcription, &cur.State); err != nil {
			return nil, err
		}
		out = append(out, cur)
	}
	return out, nil
}

func (s *Store) AuthorLeaderboard(ctx context.Context) (*models.AuthorLeaderboard, error) {

	query := `
        SELECT
		   ranks.name,
		   ranks.identity,
           ranks.approver,
           ranks.num_approved,
		   (SELECT COALESCE(SUM(r.claim_value), 0) FROM author_reward r WHERE r.author_id = ranks.id AND r.claimed = TRUE)
        FROM (
            SELECT 
            	a.id,
                a.name,
			    COALESCE(a.identity, '{}') as identity,
                a.approver,
                COALESCE(SUM(CASE WHEN c.state = 'approved' THEN 1 ELSE 0 END), 0) as num_approved
            FROM author a
            LEFT JOIN tscript_contribution c ON c.author_id = a.id
            GROUP BY a.id, a.name, a.approver) ranks
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
		ranking := &models.AuthorRanking{Author: &models.ShortAuthor{}}
		var identity string
		err := rows.Scan(
			&ranking.Author.Name,
			&identity,
			&ranking.Approver,
			&ranking.AcceptedContributions,
			&ranking.AwardValue,
		)
		if err != nil {
			return nil, err
		}

		ident := &oauth.Identity{}
		if err := json.Unmarshal([]byte(identity), &ident); err == nil {
			ranking.Author.IdentityIconImg = ident.Icon
		}

		authors = append(authors, ranking)
	}

	return &models.AuthorLeaderboard{Authors: authors}, nil
}

func (s *Store) ListRanks(ctx context.Context) ([]*models.Rank, error) {
	rows, err := s.tx.QueryxContext(ctx, "SELECT id, name, points FROM rank")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []*models.Rank{}
	for rows.Next() {
		cur := &models.Rank{}
		if err := rows.StructScan(cur); err != nil {
			return nil, err
		}
		res = append(res, cur)
	}
	return res, nil
}

func (s *Store) ListAuthorRankings(ctx context.Context, qm *common.QueryModifier) ([]*models.AuthorRank, error) {
	fieldMap := map[string]string{
		"author_id":   "a.id",
		"author_name": "a.name",
		"points":      "points",
		"rank":        "rank",
	}

	where, params, order, paging, err := qm.ToSQL(fieldMap, true)
	if err != nil {
		return nil, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
		SELECT 
			a.id,
			a.name,
 			COALESCE(a.identity, '{}'),
			COALESCE((SELECT SUM(points) AS points FROM author_contribution where author_id=a.id), 0) as points,
			COALESCE((SELECT id FROM rank WHERE points <= (SELECT COALESCE(SUM(points), 0) AS points FROM author_contribution where author_id = a.id) order by points desc limit 1), '') as current_rank,
			COALESCE((SELECT id FROM rank WHERE points > (SELECT COALESCE(SUM(points), 0) AS points FROM author_contribution where author_id = a.id) order by points asc limit 1), '') as next_rank,
			COALESCE((SELECT SUM(claim_value) FROM author_reward WHERE claimed = true AND author_id = a.id), 0) as reward_value_claimed,
			COALESCE(c.approved_chunks, 0),
			COALESCE(c.approved_changes, 0)
		FROM author a
		LEFT JOIN (
			SELECT 
				ac.author_id,
				COALESCE(SUM(CASE WHEN ac.contribution_type = 'chunk' then 1 ELSE 0 END), 0) as approved_chunks,
				COALESCE(SUM(CASE WHEN ac.contribution_type = 'change' then 1 ELSE 0 END), 0) as approved_changes
			FROM author_contribution ac
			group by ac.author_id
		) c ON a.id = c.author_id
		%s
		%s
		%s
		`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*models.AuthorRank, 0)
	for rows.Next() {
		cur := &models.AuthorRank{Author: &models.ShortAuthor{}}
		var ident string
		if err := rows.Scan(&cur.Author.ID, &cur.Author.Name, &ident, &cur.Points, &cur.CurrentRankID, &cur.NextRankID, &cur.RewardValueUSD, &cur.ApprovedChunks, &cur.ApprovedChanges); err != nil {
			return nil, err
		}
		decodedIdent := &oauth.Identity{}
		if err := json.Unmarshal([]byte(ident), &decodedIdent); err == nil {
			cur.Author.IdentityIconImg = decodedIdent.Icon
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

func (s *Store) ListRequiredAuthorRewardsV2(ctx context.Context, pointsForReward float32) ([]*models.RequiredReward, error) {
	rows, err := s.tx.QueryxContext(ctx, `
		SELECT rr.*
		FROM (
			SELECT author_id, SUM(points - points_spent) AS points_spendable 
			FROM author_contribution 
			GROUP BY author_id
		) AS rr 
		WHERE rr.points_spendable > $1
	`, pointsForReward)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*models.RequiredReward, 0)
	for rows.Next() {
		var pointsSpendable float32
		var authorID string
		if err := rows.Scan(&authorID, &pointsSpendable); err != nil {
			return nil, err
		}
		for i := 0; i < int(math.Floor(float64(pointsSpendable/pointsForReward))); i++ {
			cur := &models.RequiredReward{AuthorID: authorID, Points: pointsForReward}
			out = append(out, cur)
		}
	}
	return out, nil
}

func (s *Store) ListPendingRewards(ctx context.Context, authorID string) ([]*models.AuthorReward, error) {
	return s.listRewards(ctx, authorID, false)
}

func (s *Store) ListClaimedRewards(ctx context.Context, authorID string) ([]*models.AuthorReward, error) {
	return s.listRewards(ctx, authorID, true)
}

func (s *Store) listRewards(ctx context.Context, authorID string, claimed bool) ([]*models.AuthorReward, error) {
	rows, err := s.tx.QueryxContext(ctx, `SELECT * from author_reward WHERE author_id = $1 AND claimed = $2 AND "error" IS NULL`, authorID, claimed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := []*models.AuthorReward{}
	for rows.Next() {
		reward := &models.AuthorReward{}
		err := rows.StructScan(reward)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, reward)
	}
	return rewards, nil
}

func (s *Store) CreatePendingReward(ctx context.Context, authorID string, pointsSpent float32) (string, error) {
	id := shortuuid.New()
	if _, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO author_reward (id, author_id, points_spent, created_at) VALUES ($1, $2, $3, NOW())`,
		id,
		authorID,
		pointsSpent); err != nil {
		return "", err
	}
	return id, nil
}

// GetReward gets a reward without any locking. If the reward is going to be updated, use GetRewardForUpdate.
func (s *Store) GetReward(ctx context.Context, id string) (*models.AuthorReward, error) {
	reward := &models.AuthorReward{}
	err := s.tx.QueryRowxContext(ctx, `SELECT * from author_reward WHERE claimed = FALSE AND error IS NULL AND id = $1`, id).StructScan(reward)
	if err != nil {
		return nil, err
	}
	return reward, nil
}

func (s *Store) GetRewardForUpdate(ctx context.Context, id string) (*models.AuthorReward, error) {
	reward := &models.AuthorReward{}
	err := s.tx.QueryRowxContext(ctx, `SELECT * from author_reward WHERE claimed = FALSE AND error IS NULL AND id = $1 FOR UPDATE`, id).StructScan(reward)
	if err != nil {
		return nil, err
	}
	return reward, nil
}

func (s *Store) ClaimReward(ctx context.Context, id string, kind string, value float32, currency string, confirmationCode string, description string) error {
	if _, err := s.tx.ExecContext(
		ctx,
		`UPDATE author_reward SET claimed=true, claim_kind=$1, claim_value=$2, claim_value_currency=$3, claim_confirmation_code=$4, claim_description=$5, claim_at=NOW() WHERE id=$6`,
		kind,
		value,
		currency,
		confirmationCode,
		description,
		id); err != nil {
		return err
	}
	return nil
}

func (s *Store) FailReward(ctx context.Context, id string, reason string) error {
	if _, err := s.tx.ExecContext(ctx, `UPDATE author_reward SET error=$1 WHERE id=$2`, reason, id); err != nil {
		return err
	}
	return nil
}

func (s *Store) BatchGetAuthor(ctx context.Context, authorIDs ...string) ([]*models.Author, error) {

	placeholders := make([]string, len(authorIDs))
	params := make([]interface{}, len(authorIDs))

	for k, id := range authorIDs {
		placeholders[k] = fmt.Sprintf("$%d", k+1)
		params[k] = id
	}

	rows, err := s.tx.QueryxContext(ctx, fmt.Sprintf(`SELECT * from author WHERE id IN (%s)`, strings.Join(placeholders, ", ")), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := []*models.Author{}
	for rows.Next() {
		author := &models.Author{}
		if err := rows.StructScan(author); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (s *Store) BatchGetAuthorMap(ctx context.Context, authorIDs ...string) (map[string]*models.Author, error) {
	authors, err := s.BatchGetAuthor(ctx, authorIDs...)
	if err != nil {
		return nil, err
	}
	authorMap := make(map[string]*models.Author, 0)
	for _, v := range authors {
		authorMap[v.ID] = v
	}
	return authorMap, nil
}

func (s *Store) GetAuthor(ctx context.Context, id string) (*models.Author, error) {
	authors, err := s.BatchGetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(authors) == 0 {
		return nil, sql.ErrNoRows
	}
	return authors[0], nil
}

func (s *Store) GetTranscriptChange(ctx context.Context, id string) (*models.TranscriptChange, error) {

	change := &models.TranscriptChange{}
	var authorID string

	err := s.tx.
		QueryRowxContext(ctx, `SELECT id, author_id, epid, summary, transcription, state, created_at, merged FROM transcript_change WHERE id=$1`, id).
		Scan(&change.ID, &authorID, &change.EpID, &change.Summary, &change.Transcription, &change.State, &change.CreatedAt, &change.Merged)
	if err != nil {
		return nil, err
	}
	author, err := s.GetAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}
	change.Author = author

	return change, nil
}

func (s *Store) CreateTranscriptChange(ctx context.Context, c *models.TranscriptChangeCreate) (*models.TranscriptChange, error) {

	author, err := s.GetAuthor(ctx, c.AuthorID)
	if err != nil {
		return nil, err
	}
	if author.Banned {
		return nil, ErrNotPermitted
	}
	change := &models.TranscriptChange{
		ID:            shortuuid.New(),
		EpID:          c.EpID,
		Author:        author,
		Summary:       c.Summary,
		Transcription: c.Transcription,
		State:         models.ContributionStatePending,
		CreatedAt:     time.Now(),
	}
	_, err = s.tx.ExecContext(
		ctx,
		`INSERT INTO transcript_change (id, author_id, epid, summary, transcription, state, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		change.ID,
		change.Author.ID,
		change.EpID,
		change.Summary,
		change.Transcription,
		models.ContributionStatePending,
		change.CreatedAt.Format(util.SQLDateFormat),
	)
	return change, err
}

func (s *Store) UpdateTranscriptChange(ctx context.Context, c *models.TranscriptChangeUpdate, ppointsOnApprove float32) (*models.TranscriptChange, error) {
	if c.ID == "" {
		return nil, fmt.Errorf("no identifier was provided")
	}
	change, err := s.GetTranscriptChange(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	if change.Author.Banned {
		return nil, ErrNotPermitted
	}
	_, err = s.tx.ExecContext(
		ctx,
		`UPDATE transcript_change SET transcription=$1, summary=$2, state=$3 WHERE id=$4`,
		c.Transcription,
		c.Summary,
		c.State,
		c.ID,
	)
	if change.State != models.ContributionStateApproved && c.State == models.ContributionStateApproved {
		var contributionID string
		var err error
		if contributionID, err = s.CreateAuthorContribution(ctx, models.AuthorContributionCreate{
			AuthorID:         change.Author.ID,
			EpID:             change.EpID,
			ContributionType: models.ContributionTypeChange,
			Points:           ppointsOnApprove,
		}); err != nil {
			return nil, err
		}
		if err := s.LinkAuthorContributionToChange(ctx, contributionID, c.ID); err != nil {
			return nil, err
		}
	}

	change.Transcription = c.Transcription
	change.Summary = c.Summary
	change.State = c.State

	return change, err
}

func (s *Store) DeleteTranscriptChange(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("no identifier was provided")
	}
	_, err := s.tx.ExecContext(
		ctx,
		`DELETE FROM transcript_change WHERE id=$1`,
		id,
	)
	return err
}

func (s *Store) UpdateTranscriptChangeState(ctx context.Context, id string, state models.ContributionState, pointsOnApprove float32) error {

	change, err := s.GetTranscriptChange(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.tx.ExecContext(
		ctx,
		`UPDATE transcript_change SET state=$1 WHERE id=$2`,
		state,
		id,
	)
	if err != nil {
		return err
	}
	if change.State != models.ContributionStateApproved && state == models.ContributionStateApproved {
		var contributionID string
		var err error
		if contributionID, err = s.CreateAuthorContribution(ctx, models.AuthorContributionCreate{
			AuthorID:         change.Author.ID,
			EpID:             change.EpID,
			ContributionType: models.ContributionTypeChange,
			Points:           pointsOnApprove,
		}); err != nil {
			return err
		}
		if err := s.LinkAuthorContributionToChange(ctx, contributionID, change.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListTranscriptChanges(ctx context.Context, q *common.QueryModifier) ([]*models.TranscriptChange, error) {
	fieldMap := map[string]string{
		"id":            "c.id",
		"author_id":     "c.author_id",
		"epid":          "c.epid",
		"summary":       "c.summary",
		"transcription": "c.transcription",
		"state":         "c.state",
		"created_at":    "c.created_at",
		"merged":        "c.merged",
	}

	where, params, order, paging, err := q.ToSQL(fieldMap, false)
	if err != nil {
		return nil, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
		SELECT c.id, c.author_id, c.epid, c.summary, c.transcription, c.state, c.created_at, c.merged, COALESCE(con.points, 0)
		FROM transcript_change c
		LEFT JOIN author a ON c.author_id = a.id
		LEFT JOIN author_contribution_transcript_change actc ON c.id = actc.transcript_change_id
		LEFT JOIN author_contribution con ON con.id = actc.author_contribution_id AND con.contribution_type='change'
		WHERE a.id IS NOT NULL AND a.banned = false
		AND %s
		%s
		%s
		`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authorIDs []string
	out := make([]*models.TranscriptChange, 0)
	for rows.Next() {
		cur := &models.TranscriptChange{Author: &models.Author{}}
		if err := rows.Scan(
			&cur.ID,
			&cur.Author.ID,
			&cur.EpID,
			&cur.Summary,
			&cur.Transcription,
			&cur.State,
			&cur.CreatedAt,
			&cur.Merged,
			&cur.PointsAwarded); err != nil {
			return nil, err
		}
		authorIDs = append(authorIDs, cur.Author.ID)
		out = append(out, cur)
	}
	if len(authorIDs) == 0 {
		return out, nil
	}

	authorMap, err := s.BatchGetAuthorMap(ctx, authorIDs...)
	if err != nil {
		return nil, err
	}
	for k := range out {
		out[k].Author = authorMap[out[k].Author.ID]
	}
	return out, nil
}

func (s *Store) CreateAuthorContribution(ctx context.Context, co models.AuthorContributionCreate) (string, error) {
	id := shortuuid.New()
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO author_contribution (id, author_id, epid, contribution_type, points, points_spent, created_at) VALUES ($1, $2, $3, $4, $5, 0, $6)`,
		id,
		co.AuthorID,
		co.EpID,
		co.ContributionType,
		co.Points,
		time.Now(),
	)
	return id, err
}

func (s *Store) LinkAuthorContributionToChange(ctx context.Context, contributionID string, changeID string) error {
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO author_contribution_transcript_change (author_contribution_id, transcript_change_id) VALUES ($1, $2)`,
		contributionID,
		changeID,
	)
	return err
}

func (s *Store) LinkAuthorContributionToTscriptContribution(ctx context.Context, contributionID string, tscriptContributionID string) error {
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO author_contribution_tscript_contribution (author_contribution_id, tscript_contribution_id) VALUES ($1, $2)`,
		contributionID,
		tscriptContributionID,
	)
	return err
}

func (s *Store) SpendPoints(ctx context.Context, authorId string, spendRequired float32) error {

	// use FOR UPDATE SKIP LOCKED just in case multiple worker processes run at once.
	// this should stop the same points being spent more than once.
	rows, err := s.tx.QueryxContext(
		ctx,
		`
		SELECT id, points - points_spent AS remainder 
		FROM author_contribution 
		WHERE author_id = $1
	  	AND points - points_spent > 0
   		FOR UPDATE SKIP LOCKED
		`,
		authorId,
	)
	if err != nil {
		return err
	}

	//
	// Buffer the spendable data in memory, since it's not possible to run 2 queries at once with the pq driver.
	//
	spendableRows := make([]struct {
		id        string
		remainder float32
	}, 0)
	if err := func() error {
		defer rows.Close()
		for rows.Next() {
			spendable := struct {
				id        string
				remainder float32
			}{}

			if err := rows.Scan(&spendable.id, &spendable.remainder); err != nil {
				return err
			}
			spendableRows = append(spendableRows, spendable)
		}
		return nil
	}(); err != nil {
		return err
	}

	//
	// now spend as many points as possible
	//
	for _, v := range spendableRows {
		var rowSpend float32
		if v.remainder > spendRequired {
			rowSpend = spendRequired
		} else {
			rowSpend = v.remainder
		}
		if _, err := s.tx.ExecContext(ctx, `UPDATE author_contribution SET points_spent = points_spent + $1 WHERE id=$2`, rowSpend, v.id); err != nil {
			return err
		}
		spendRequired = spendRequired - rowSpend
		if spendRequired == 0 {
			break
		}
	}

	if spendRequired > 0 {
		return fmt.Errorf("unable to spend spendRequired, %0.2f was not allocated", spendRequired)
	}
	return nil
}

func (s *Store) ListAuthorContributions(ctx context.Context, q *common.QueryModifier) ([]*models.AuthorContribution, error) {
	fieldMap := map[string]string{
		"id":                "c.id",
		"author_id":         "c.author_id",
		"epid":              "c.epid",
		"contribution_type": "c.contribution_type",
		"created_at":        "c.created_at",
	}

	where, params, order, paging, err := q.ToSQL(fieldMap, false)
	if err != nil {
		return nil, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`
		SELECT c.id, c.epid, c.contribution_type, COALESCE(c.points, 0), c.created_at, a.id, a.name
		FROM author_contribution c
		LEFT JOIN author a ON c.author_id = a.id
		WHERE a.id IS NOT NULL AND a.banned = false
		AND %s
		%s
		%s
		`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*models.AuthorContribution, 0)
	for rows.Next() {
		cur := &models.AuthorContribution{Author: &models.ShortAuthor{}}
		if err := rows.Scan(
			&cur.ID,
			&cur.EpID,
			&cur.ContributionType,
			&cur.Points,
			&cur.CreatedAt,
			&cur.Author.ID,
			&cur.Author.Name); err != nil {
			return nil, err
		}
		out = append(out, cur)
	}
	return out, nil
}

func (s *Store) CreateTscriptImport(ctx context.Context, tscriptImport *models.TscriptImportCreate) (*models.TscriptImport, error) {
	imp := &models.TscriptImport{
		ID:     shortuuid.New(),
		EpID:   tscriptImport.EpID,
		Mp3URI: tscriptImport.Mp3URI,
	}
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO tscript_import (id, epid, mp3_uri) VALUES ($1, $2, $3)`,
		imp.ID,
		imp.EpID,
		imp.Mp3URI,
	)
	if err != nil {
		return nil, err
	}
	return imp, err
}

func (s *Store) PushTscriptImportLog(ctx context.Context, tscriptImportID string, log *models.TscriptImportLog) error {
	logJSON, err := json.Marshal([]*models.TscriptImportLog{log})
	if err != nil {
		return fmt.Errorf("failed to marshal log to JSON: %s", err.Error())
	}
	_, err = s.tx.ExecContext(ctx, `UPDATE tscript_import SET log=((CASE WHEN log IS NULL THEN '[]'::JSONB ELSE log END) || $1::JSONB) WHERE id=$2 `, string(logJSON), tscriptImportID)
	return err
}
