package ro

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/util"
)

//go:embed migrations
var migrations embed.FS

func NewConn(cfg *common.Config) (*Conn, error) {
	innerConn, err := common.NewConn("sqlite3", cfg)
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

func (s *Store) InsertEpisodeWithTranscript(ctx context.Context, ep *models.Transcript) error {

	epMeta, err := metaToString(ep.Meta)
	if err != nil {
		return err
	}
	epContributors, err := contributorsToString(ep.Contributors)
	if err != nil {
		return err
	}
	_, err = s.tx.ExecContext(
		ctx,
		`INSERT INTO episode (id, publication_type, publication, series, episode, release_date, metadata, contributors) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`,
		ep.ID(),
		ep.PublicationType,
		ep.Publication,
		ep.Series,
		ep.Episode,
		util.SqlDate(ep.ReleaseDate),
		epMeta,
		epContributors,
	)

	for _, v := range ep.Transcript {
		diaMeta, err := metaToString(v.Meta)
		if err != nil {
			return err
		}
		_, err = s.tx.ExecContext(ctx,
			`INSERT INTO dialog (id, episode_id, pos, offset, offset_inferred, type, actor, content, metadata, notable) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			v.ID,
			ep.ID(),
			v.Position,
			v.Timestamp,
			v.TimestampInferred,
			string(v.Type),
			v.Actor,
			v.Content,
			diaMeta,
			v.Notable,
		)
		if err != nil {
			return errors.Wrapf(err, "failed at line: %s", v.ID)
		}
	}
	return err
}

func (s *Store) GetDialogWithContext(ctx context.Context, dialogID string, withContext int32) ([]models.Dialog, string, error) {
	query := fmt.Sprintf(`
		WITH target AS (SELECT * FROM dialog WHERE id = $1 LIMIT 1)
		SELECT * FROM (SELECT * FROM dialog WHERE pos < (SELECT pos FROM target) AND episode_id = (SELECT episode_id FROM target) ORDER BY pos DESC LIMIT %d)
		UNION 
		SELECT * FROM target
		UNION
		SELECT * FROM (SELECT * FROM dialog WHERE pos > (SELECT pos FROM target) AND episode_id = (SELECT episode_id FROM target) ORDER BY pos ASC LIMIT %d)
		ORDER BY pos ASC`, withContext, withContext)

	return s.getTranscriptForQuery(ctx, query, dialogID)
}

func (s *Store) InsertChangelog(ctx context.Context, ep *models.Changelog) error {
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO changelog ("date", "content") VALUES ($2, $3)`,
		util.SqlDate(&ep.Date),
		ep.Content,
	)
	return err
}

func (s *Store) ListEpisodes(ctx context.Context, q *common.QueryModifier) ([]*models.EpisodeMeta, error) {
	fieldMap := map[string]string{
		"publication_type": "publication_type",
		"publication":      "publication",
		"series":           "series",
		"episode":          "episode",
		"release_date":     "release_date",
	}

	q.Apply(common.WithDefaultSorting("release_date", common.SortAsc))

	where, params, order, paging, err := q.ToSQL(fieldMap, true)
	if err != nil {
		return nil, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`SELECT publication_type, publication, series, episode, release_date  FROM episode %s %s %s`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	result := make([]*models.EpisodeMeta, 0)
	for rows.Next() {
		row := &models.EpisodeMeta{}
		if err := rows.Scan(&row.PublicationType, &row.Publication, &row.Series, &row.Episode, &row.ReleaseDate); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s *Store) ListChangelogs(ctx context.Context, q *common.QueryModifier) ([]*models.Changelog, error) {
	fieldMap := map[string]string{
		"date":    "date",
		"content": "content",
	}

	q.Apply(common.WithDefaultSorting("date", common.SortDesc))

	where, params, order, paging, err := q.ToSQL(fieldMap, true)
	if err != nil {
		return nil, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`SELECT date, content FROM changelog %s %s %s`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	out := make([]*models.Changelog, 0)
	for rows.Next() {
		row := &models.Changelog{}
		if err := rows.Scan(&row.Date, &row.Content); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

func (s *Store) getTranscriptForQuery(ctx context.Context, query string, params ...interface{}) ([]models.Dialog, string, error) {

	res, err := s.tx.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, "", err
	}
	defer func(res *sqlx.Rows) {
		_ = res.Close()
	}(res)

	var epID string

	results := make([]models.Dialog, 0)
	for res.Next() {

		result := models.Dialog{
			Meta: make(models.Metadata),
		}
		var meta string

		if err := res.Scan(&result.ID, &epID, &result.Position, &result.Timestamp, &result.TimestampInferred, &result.Type, &result.Actor, &result.Content, &meta, &result.Notable); err != nil {
			return nil, "", err
		}
		if meta != "" {
			if err := json.Unmarshal([]byte(meta), &result.Meta); err != nil {
				return nil, "", errors.Wrap(err, "failed to unmarshal meta")
			}
		}
		results = append(results, result)
	}
	return results, epID, nil
}

func (s *Store) InsertSong(ctx context.Context, song meta.Song) error {
	episodeIds, err := json.Marshal(song.EpisodeIDs)
	if err != nil {
		return err
	}
	transcribed, err := json.Marshal(song.Transcribed)
	if err != nil {
		return err
	}
	_, err = s.tx.ExecContext(
		ctx,
		`INSERT INTO song ("spotify_uri", "artist", "title", "album", "episode_ids", "album_image_url", "transcribed") VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		song.Track.TrackURI,
		song.Track.Artist(),
		song.Track.Name,
		song.Track.AlbumName,
		string(episodeIds),
		song.Track.AlbumImageUrl,
		string(transcribed),
	)
	return err
}

func (s *Store) ListSongs(ctx context.Context, q *common.QueryModifier) ([]*models.Song, int64, error) {
	fieldMap := map[string]string{
		"artist":      "artist",
		"title":       "title",
		"album":       "album",
		"episode_ids": "episode_ids",
	}

	q.Apply(common.WithDefaultSorting("title", common.SortAsc))

	where, params, order, paging, err := q.ToSQL(fieldMap, true)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`SELECT spotify_uri, artist, title, album, episode_ids, album_image_url, transcribed, COUNT() OVER() as total_rows FROM song %s %s %s`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var totalRows int64 //identical for every row
	out := make([]*models.Song, 0)
	for rows.Next() {
		row := &models.Song{}
		var episodeIDRaw []byte
		var transcribedRaw []byte
		if err := rows.Scan(&row.SpotifyURI, &row.Artist, &row.Title, &row.Album, &episodeIDRaw, &row.AlbumImageURL, &transcribedRaw, &totalRows); err != nil {
			return nil, 0, err
		}
		if len(episodeIDRaw) > 0 {
			if err := json.Unmarshal(episodeIDRaw, &row.EpisodeIDs); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal episode IDs: %w", err)
			}
		}
		if len(transcribedRaw) > 0 {
			if err := json.Unmarshal(transcribedRaw, &row.Transcribed); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal transcription lines: %w", err)
			}
		}
		out = append(out, row)
	}
	return out, totalRows, nil
}

func (s *Store) InsertCommunityProject(ctx context.Context, proj models.CommunityProject) error {
	_, err := s.tx.ExecContext(
		ctx,
		`INSERT INTO community_project ("id", "name", "summary", "content", "url", "created_at") VALUES ($1, $2, $3, $4, $5, $6)`,
		proj.ID,
		proj.Name,
		proj.Summary,
		proj.Content,
		proj.URL,
		proj.CreatedAt,
	)
	return err
}

func (s *Store) ListCommunityProjects(ctx context.Context, q *common.QueryModifier) (models.CommunityProjects, int64, error) {
	fieldMap := map[string]string{
		"id":         "id",
		"name":       "name",
		"created_at": "created_at",
	}
	q.Apply(common.WithDefaultSorting("created_at", common.SortAsc))

	where, params, order, paging, err := q.ToSQL(fieldMap, true)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.tx.QueryxContext(
		ctx,
		fmt.Sprintf(`SELECT id, name, summary, content, url, created_at, COUNT() OVER() as total_rows FROM community_project %s %s %s`, where, order, paging),
		params...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var totalRows int64 //identical for every row
	out := make([]models.CommunityProject, 0)
	for rows.Next() {
		row := models.CommunityProject{}
		if err := rows.Scan(&row.ID, &row.Name, &row.Summary, &row.Content, &row.URL, &row.CreatedAt, &totalRows); err != nil {
			return nil, 0, err
		}
		out = append(out, row)
	}
	return out, totalRows, nil
}

func metaToString(metadata models.Metadata) (string, error) {
	if metadata == nil {
		return "", nil
	}
	bs, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func contributorsToString(contributors []string) (string, error) {
	if contributors == nil {
		return "", nil
	}
	bs, err := json.Marshal(contributors)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
