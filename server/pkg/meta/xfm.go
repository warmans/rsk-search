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
}

func XfmEpisodeNames() map[string]string {
	cpy := map[string]string{}
	for k, v := range parsedIndex {
		cpy[k] = v
	}
	return cpy
}

func XfmEpisodeDateToName(date time.Time) (string, error) {
	dateStr := date.Format(time.RFC3339)
	name, ok := parsedIndex[dateStr]
	if !ok {
		return "", fmt.Errorf("date %s did not map to any known episode", dateStr)

	}
	return name, nil
}
