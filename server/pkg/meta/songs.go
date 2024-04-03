package meta

import (
	"embed"
	"encoding/json"
	"github.com/warmans/rsk-search/pkg/spotify"
	"slices"
)

//go:embed data/songs.json
var songs embed.FS

var songMeta = SongMetaMap{}

func init() {
	f, err := songs.Open("data/songs.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&songMeta); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
}

func GetSongMeta() SongMetaMap {
	return songMeta
}

type Song struct {
	Terms      []string
	EpisodeIDs []string
	Track      *spotify.Track
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
