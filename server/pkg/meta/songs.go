package meta

import (
	"embed"
	"encoding/json"
	"github.com/warmans/rsk-search/pkg/spotify"
)

//go:embed data/tags.json
var songs embed.FS

var songMeta = Tags{}

func init() {
	f, err := tagIndex.Open("data/songs.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&songMeta); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
}

type Song struct {
	Terms      []string
	EpisodeIDs []string
	Track      *spotify.Track
}

type Songs struct {
	Songs map[string]*Song
}

func (s *Songs) FindKeyByTerm(term string) (string, bool) {
	for k, v := range s.Songs {
		for _, t := range v.Terms {
			if term == t {
				return k, true
			}
		}
	}
	return "", false
}

