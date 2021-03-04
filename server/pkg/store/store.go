package store

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"io"
	"path"
	"strings"
)

//go:embed migrations
var migrations embed.FS

type Config struct {
	DSN string
}

func (c *Config) RegisterFlags(prefix string) {
	flag.StringVarEnv(&c.DSN, prefix, "db-dsn", "./var/rsk.sqlite3", "DB connection string")
}

func NewConn(cfg *Config) (*Conn, error) {
	db, err := sqlx.Connect("sqlite3", cfg.DSN)
	if err != nil {
		return nil, err
	}
	return &Conn{db: db}, nil
}

type Conn struct {
	db *sqlx.DB
}

func (c *Conn) Migrate() error {

	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS migration_log (
		  file_name TEXT PRIMARY KEY
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create metadata table: %w", err)
	}

	appliedMigrations := []string{}
	err = c.WithTx(func(tx *sqlx.Tx) error {
		rows, err := tx.Query("SELECT file_name FROM migration_log ORDER BY file_name DESC")
		if err != nil {
			return fmt.Errorf("failed to get migrations: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return err
			}
			appliedMigrations = append(appliedMigrations, name)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = c.WithTx(func(tx *sqlx.Tx) error {
		entries, err := migrations.ReadDir("migrations")
		if err != nil {
			return err
		}
		for _, dirEntry := range entries {
			if !strings.HasSuffix(dirEntry.Name(), ".sql") {
				continue
			}
			if util.InStrings(dirEntry.Name(), appliedMigrations...) {
				continue
			}

			migrationPath := path.Join("migrations", dirEntry.Name())
			f, err := migrations.Open(migrationPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", migrationPath, err)
			}
			defer f.Close()

			bytes, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			if _, err := tx.Exec(string(bytes)); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", dirEntry.Name(), err)
			}
			if _, err := tx.Exec("INSERT INTO migration_log (file_name) VALUES ($1)", dirEntry.Name()); err != nil {
				return fmt.Errorf("failed to update migration log: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to apply migraions: %w", err)
	}
	return nil
}

func (c *Conn) WithTx(f func(tx *sqlx.Tx) error) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	if err := f(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf("failed to rollback (%s) from error : %w)", err2.Error(), err)
		}
		return err
	}
	return tx.Commit()
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
	_, err = s.tx.ExecContext(
		ctx,
		`INSERT INTO episode (id, publication, series, episode, release_date, metadata) VALUES ($1, $2, $3, $4, $5, $6)`,
		util.EpisodeName(ep),
		ep.Publication,
		ep.Series,
		ep.Episode,
		util.SqlDate(ep.ReleaseDate),
		epMeta,
	)

	for _, v := range ep.Transcript {
		diaMeta, err := metaToString(v.Meta)
		if err != nil {
			return err
		}
		_, err = s.tx.ExecContext(ctx,
			`INSERT INTO dialog (id, episode_id, pos, type, actor, content, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			v.ID,
			util.EpisodeName(ep),
			v.Position,
			string(v.Type),
			v.Actor,
			v.Content,
			diaMeta,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Store) GetShortEpisode(ctx context.Context, id string) (*models.Episode, error) {
	ep := &models.Episode{}
	return ep, nil
}

func limitStmnt(pageSize int32, page int32) string {
	if pageSize < 1 {
		pageSize = 25
	}
	if page < 1 {
		page = 1
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, pageSize*(page-1))
}

func metaToString(metadata models.Metadata) (string, error) {
	if metadata == nil {
		return "", nil
	}
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(metaBytes), nil
}
