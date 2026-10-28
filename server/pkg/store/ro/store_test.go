package ro

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/spotify"
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

func TestInsertAndListSongsSortedByFirstOccurrence(t *testing.T) {
	conn, err := NewConn(&common.Config{DSN: "file:songstoretest?mode=memory&cache=shared"})
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Migrate(); err != nil {
		t.Fatal(err)
	}

	mustDate := func(date string) *time.Time {
		parsed, err := time.Parse(util.ShortDateFormat, date)
		if err != nil {
			t.Fatal(err)
		}
		return &parsed
	}

	songs := []meta.Song{
		{
			Transcribed: []string{"song: later song"},
			EpisodeIDs:  []string{"xfm-S2E10"},
			Track: &spotify.Track{
				Artists:   []spotify.Artist{{Name: "Artist C"}},
				Name:      "Later Song",
				TrackURI:  "spotify:track:3",
				AlbumName: "Album C",
			},
			FirstOccurrence: mustDate("2001-11-01"),
		},
		{
			Transcribed: []string{"song: earliest song"},
			EpisodeIDs:  []string{"xfm-S0E02", "xfm-S1E05"},
			Track: &spotify.Track{
				Artists:   []spotify.Artist{{Name: "Artist B"}},
				Name:      "Earliest Song",
				TrackURI:  "spotify:track:2",
				AlbumName: "Album B",
			},
			FirstOccurrence: mustDate("1998-05-31"),
		},
		{
			Transcribed: []string{"song: middle song"},
			EpisodeIDs:  []string{"xfm-S1E01"},
			Track: &spotify.Track{
				Artists:   []spotify.Artist{{Name: "Artist A"}},
				Name:      "Middle Song",
				TrackURI:  "spotify:track:1",
				AlbumName: "Album A",
			},
			FirstOccurrence: mustDate("1999-08-17"),
		},
		{
			Transcribed: []string{"song: unknown date"},
			EpisodeIDs:  []string{"xfm-S9E99"},
			Track: &spotify.Track{
				Artists:   []spotify.Artist{{Name: "Artist D"}},
				Name:      "Unknown Date Song",
				TrackURI:  "spotify:track:4",
				AlbumName: "Album D",
			},
		},
	}

	if err := conn.WithStore(func(s *Store) error {
		for _, v := range songs {
			if err := s.InsertSong(context.Background(), v); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	listSorted := func(direction common.SortDirection) []*models.Song {
		var out []*models.Song
		if err := conn.WithStore(func(s *Store) error {
			var err error
			out, _, err = s.ListSongs(context.Background(), common.Q(common.WithSorting("first_occurrence", direction)))
			return err
		}); err != nil {
			t.Fatal(err)
		}
		return out
	}

	asc := listSorted(common.SortAsc)
	if len(asc) != 4 {
		t.Fatalf("expected 4 songs, got %d", len(asc))
	}
	gotAsc := []string{asc[0].Title, asc[1].Title, asc[2].Title, asc[3].Title}
	wantAsc := []string{"Unknown Date Song", "Earliest Song", "Middle Song", "Later Song"}
	for k := range wantAsc {
		if gotAsc[k] != wantAsc[k] {
			t.Errorf("asc mismatch: got %v, expected %v", gotAsc, wantAsc)
			break
		}
	}

	desc := listSorted(common.SortDesc)
	gotDesc := []string{desc[0].Title, desc[1].Title, desc[2].Title, desc[3].Title}
	wantDesc := []string{"Later Song", "Middle Song", "Earliest Song", "Unknown Date Song"}
	for k := range wantDesc {
		if gotDesc[k] != wantDesc[k] {
			t.Errorf("desc mismatch: got %v, expected %v", gotDesc, wantDesc)
			break
		}
	}

	if got, want := asc[1].FirstOccurrence.Format(util.ShortDateFormat), "1998-05-31"; got != want {
		t.Errorf("firstOccurrence mismatch: got %q, expected %q", got, want)
	}
	if asc[0].FirstOccurrence != nil {
		t.Errorf("expected nil firstOccurrence for song with no known episode dates, got %q", asc[0].FirstOccurrence.Format(util.ShortDateFormat))
	}
	if got, want := asc[1].EpisodeIDs[0], "xfm-S0E02"; got != want {
		t.Errorf("episodeIDs mismatch: got %q, expected %q", got, want)
	}
}
