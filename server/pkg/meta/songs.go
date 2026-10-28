package meta

import (
	"embed"
	"encoding/json"
	"github.com/warmans/rsk-search/pkg/spotify"
	"io/fs"
	"slices"
	"time"
)

//go:embed data/songs.json
var songs embed.FS

var songMeta = SongMetaMap{}

func init() {
	f, err := songs.Open("data/songs.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	dec := json.NewDecoder(f)
	if err := dec.Decode(&songMeta); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
}

func GetSongMeta() SongMetaMap {
	return songMeta
}

type Song struct {
	// the song: line from the transcript
	Transcribed []string
	// the search terms generated from the raw transcript lines
	Terms      []string
	EpisodeIDs []string
	Track      *spotify.Track
	// release date of the earliest episode the song appears in. Derived on extraction, not persisted.
	FirstOccurrence *time.Time `json:"-"`
}

type SongMetaMap struct {
	Songs map[string]*Song
}

type Songs []Song

func (s Songs) FindKeyByTerm(term string) (string, bool) {
	for _, v := range s {
		for _, t := range v.Terms {
			if term == t {
				return v.Track.TrackURI, true
			}
		}
	}
	return "", false
}

func (s SongMetaMap) ExtractSorted() Songs {
	songSlice := Songs{}
	for _, v := range s.Songs {
		if v == nil || v.Track == nil {
			continue
		}
		v.FirstOccurrence = firstOccurrence(v.EpisodeIDs)
		songSlice = append(songSlice, *v)
	}
	slices.SortFunc(songSlice, func(a, b Song) int {
		if a.Track.Name > b.Track.Name {
			return 1
		}
		return -1
	})
	return songSlice
}

// firstOccurrence returns the release date of the earliest known episode the song appears in, or nil if none are known.
func firstOccurrence(episodeIDs []string) *time.Time {
	var first *time.Time
	for _, id := range episodeIDs {
		date, ok := episodeDatesByID[id]
		if !ok {
			continue
		}
		if first == nil || date.Before(*first) {
			first = &date
		}
	}
	return first
}
