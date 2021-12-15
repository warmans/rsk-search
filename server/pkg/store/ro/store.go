package ro

import (
	"context"
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
		`INSERT INTO episode (id, publication, series, episode, release_date, metadata, contributors) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`,
		ep.ID(),
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
			`INSERT INTO dialog (id, episode_id, pos, offset, type, actor, content, metadata, notable) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			v.ID,
			ep.ID(),
			v.Position,
			v.OffsetSec,
			string(v.Type),
			v.Actor,
			v.Content,
			diaMeta,
			v.Notable,
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

		if err := res.Scan(&result.ID, &epID, &result.Position, &result.OffsetSec, &result.Type, &result.Actor, &result.Content, &meta, &result.Notable); err != nil {
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
