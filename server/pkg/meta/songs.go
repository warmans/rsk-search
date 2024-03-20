package meta

import (
	"embed"
	"encoding/json"
	"github.com/warmans/rsk-search/pkg/spotify"
	"slices"
)

//go:embed data/songs.json
var songs embed.FS

var songMeta = Songs{}

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

func GetSongMeta() Songs {
	return songMeta
}

type Song struct {
	Terms      []string
	EpisodeIDs []string
	Track      *spotify.Track
}

type Songs struct {
	Songs map[string]*Song
}

func (s Songs) ExtractSorted() []Song {
	songSlice := []Song{}
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

func (s Songs) FindKeyByTerm(term string) (string, bool) {
	for k, v := range s.Songs {
		for _, t := range v.Terms {
			if term == t {
				return k, true
			}
		}
	}
	return "", false
}
