package ro

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (s *Store) InsertEpisodeWithTranscript(ctx context.Context, ep *models.Episode) error {

	epMeta, err := metaToString(ep.Meta)
	if err != nil {
		return err
	}
	epTags, err := tagListToString(ep.Tags)
	if err != nil {
		return err
	}
	epContributors, err := contributorsToString(ep.Contributors)
	if err != nil {
		return err
	}
	_, err = s.tx.ExecContext(
		ctx,
		`INSERT INTO episode (id, publication, series, episode, release_date, metadata, tags, contributors) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`,
		ep.ID(),
		ep.Publication,
		ep.Series,
		ep.Episode,
		util.SqlDate(ep.ReleaseDate),
		epMeta,
		epTags,
		epContributors,
	)

	for _, v := range ep.Transcript {
		diaMeta, err := metaToString(v.Meta)
		if err != nil {
			return err
		}
		diaTags, err := tagMapToString(v.ContentTags)
		if err != nil {
			return err
		}
		_, err = s.tx.ExecContext(ctx,
			`INSERT INTO dialog (id, episode_id, pos, type, actor, content, metadata, content_tags, notable, contributor) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			v.ID,
			ep.ID(),
			v.Position,
			string(v.Type),
			v.Actor,
			v.Content,
			diaMeta,
			diaTags,
			v.Notable,
			v.Contributor,
		)
		if err != nil {
			return err
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

func (s *Store) GetShortEpisode(ctx context.Context, id string) (*models.Episode, error) {
	ep := &models.Episode{}
	var metadata string
	var tags string
	var contributors string

	err := s.tx.
		QueryRowxContext(ctx, "SELECT publication, series, episode, release_date, metadata, tags, contributors FROM episode WHERE id = $1", id).
		Scan(&ep.Publication, &ep.Series, &ep.Episode, &ep.ReleaseDate, &metadata, &tags, &contributors)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if metadata != "" {
		if err := json.Unmarshal([]byte(metadata), &ep.Meta); err != nil {
			return nil, errors.Wrap(err, "failed to decode metadata")
		}
	}
	if tags != "" {
		if err := json.Unmarshal([]byte(tags), &ep.Tags); err != nil {
			return nil, errors.Wrap(err, "failed to decode tags")
		}
	}
	if contributors != "" {
		if err := json.Unmarshal([]byte(contributors), &ep.Contributors); err != nil {
			return nil, errors.Wrap(err, "failed to decode contributors")
		}
	}
	return ep, nil
}

func (s *Store) GetEpisode(ctx context.Context, id string) (*models.Episode, error) {
	ep, err := s.GetShortEpisode(ctx, id)
	if err != nil {
		return nil, err
	}
	if ep == nil {
		return nil, nil
	}
	ep.Transcript, _, err = s.getTranscriptForQuery(ctx, "SELECT * FROM dialog WHERE episode_id=$1 ORDER BY pos ASC", id)
	if err != nil {
		return nil, err
	}
	return ep, nil
}

func (s *Store) ListEpisodes(ctx context.Context) ([]*models.ShortEpisode, error) {

	results, err := s.tx.QueryxContext(ctx, "SELECT e.id, e.publication, e.series, e.episode, e.release_date, (SELECT COUNT(*) FROM dialog WHERE episode_id = e.id LIMIT 1) > 0 AS transcript_available FROM episode e ORDER BY e.series ASC, e.episode ASC")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	eps := []*models.ShortEpisode{}
	for results.Next() {
		ep := &models.ShortEpisode{}
		if err := results.Scan(&ep.ID, &ep.Publication, &ep.Series, &ep.Episode, &ep.ReleaseDate, &ep.TranscriptAvailable); err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}
	return eps, nil
}

func (s *Store) getTranscriptForQuery(ctx context.Context, query string, params ...interface{}) ([]models.Dialog, string, error) {

	res, err := s.tx.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, "", err
	}
	defer res.Close()

	var epID string

	results := make([]models.Dialog, 0)
	for res.Next() {

		result := models.Dialog{
			Meta: make(models.Metadata),
		}
		var meta string
		var tags string

		if err := res.Scan(&result.ID, &epID, &result.Position, &result.Type, &result.Actor, &result.Content, &meta, &tags, &result.Notable, &result.Contributor); err != nil {
			return nil, "", err
		}
		if meta != "" {
			if err := json.Unmarshal([]byte(meta), &result.Meta); err != nil {
				return nil, "", errors.Wrap(err, "failed to unmarshal meta")
			}
		}
		if tags != "" {
			if err := json.Unmarshal([]byte(tags), &result.ContentTags); err != nil {
				return nil, "", errors.Wrap(err, "failed to unmarshal tags")
			}
		}
		results = append(results, result)
	}
	return results, epID, nil
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

func tagListToString(tags []string) (string, error) {
	if tags == nil {
		return "", nil
	}
	bs, err := json.Marshal(tags)
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

func tagMapToString(tags map[string]string) (string, error) {
	if tags == nil {
		return "", nil
	}
	bs, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
