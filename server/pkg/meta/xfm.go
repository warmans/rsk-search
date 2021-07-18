package meta

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"
)

//go:embed data/xfm-episode-map.json
var xfmEpisodeIndex embed.FS

var parsedIndex = map[string]string{}
var episodeIDMap = map[string]struct{}{}

func init() {
	f, err := xfmEpisodeIndex.Open("data/xfm-episode-map.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	if err := dec.Decode(&parsedIndex); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
	for _, v := range parsedIndex {
		episodeIDMap[fmt.Sprintf("ep-%s-%s", PublicationXFM, v)] = struct{}{}
	}
}

func XfmEpisodeNames() map[string]string {
	cpy := map[string]string{}
	for k, v := range parsedIndex {
		cpy[k] = v
	}
	return cpy
}

// IsValidEpisodeID returns true if the value is a known episode ID according to the episode-map.json
// The id is in the format ep-[publication]-S[season]E[episode]
func IsValidEpisodeID(id string) bool {
	_, ok := episodeIDMap[id]
	return ok
}

func XfmEpisodeDateToName(date time.Time) (string, error) {
	dateStr := date.Format(time.RFC3339)
	name, ok := parsedIndex[dateStr]
	if !ok {
		return "", fmt.Errorf("date %s did not map to any known episode", dateStr)

	}
	return name, nil
}
