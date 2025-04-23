package common

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/util"
	"io"
	"path"
	"strings"
)

type Config struct {
	DSN string
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string, dbName string) {
	flag.StringVarEnv(
		fs,
		&c.DSN,
		prefix,
		fmt.Sprintf("%s-db-dsn", dbName),
		fmt.Sprintf("./var/gen/%s.sqlite3", dbName),
		"DB connection string",
	)
}

func NewConn(driver string, cfg *Config) (*Conn, error) {
	db, err := sqlx.Connect(driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	return &Conn{db: db}, nil
}

type Conn struct {
	db *sqlx.DB
}

func (c *Conn) Migrate(migrations embed.FS) error {

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
		defer func(rows *sql.Rows) {
			_ = rows.Close()
		}(rows)

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

			bytes, err := io.ReadAll(f)
			_ = f.Close()
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

func (c *Conn) Close() error {
	return c.db.Close()
}
