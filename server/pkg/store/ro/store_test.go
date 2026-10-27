package ro

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/util"
)

func TestInsertAndListChangelog(t *testing.T) {
	conn, err := NewConn(&common.Config{DSN: "file:rostoretest?mode=memory&cache=shared"})
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Migrate(); err != nil {
		t.Fatal(err)
	}

	in := &models.Changelog{
		Date:    time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		Content: "something changed",
	}

	if err := conn.WithStore(func(s *Store) error {
		return s.InsertChangelog(context.Background(), in)
	}); err != nil {
		t.Fatal(err)
	}

	var out []*models.Changelog
	if err := conn.WithStore(func(s *Store) error {
		var err error
		out, err = s.ListChangelogs(context.Background(), &common.QueryModifier{})
		return err
	}); err != nil {
		t.Fatal(err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 changelog entry, got %d", len(out))
	}
	if out[0].Content != in.Content {
		t.Errorf("content mismatch: got %q, expected %q", out[0].Content, in.Content)
	}
	if got, want := out[0].Date.Format(util.SQLDateFormat), in.Date.Format(util.SQLDateFormat); got != want {
		t.Errorf("date mismatch: got %q, expected %q", got, want)
	}
}
