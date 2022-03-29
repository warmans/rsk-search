package meta

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/episode-date-map.json
var episodeDateIndex embed.FS

//go:embed data/episode-order.json
var orderedEpisodes embed.FS

var parsedIndex = map[string]string{}
var episodeIDMap = map[string]struct{}{}
var episodeOrder = []string{}

func init() {
	epIndexFile, err := episodeDateIndex.Open("data/episode-date-map.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer epIndexFile.Close()
	dec := json.NewDecoder(epIndexFile)
	if err := dec.Decode(&parsedIndex); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
	for _, v := range parsedIndex {
		episodeIDMap[fmt.Sprintf("ep-%s", v)] = struct{}{}
	}

	orderedEpsFile, err := orderedEpisodes.Open("data/episode-order.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer orderedEpsFile.Close()
	if err := json.NewDecoder(orderedEpsFile).Decode(&episodeOrder); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
}

func EpisodeDates() map[string]string {
	cpy := map[string]string{}
	for k, v := range parsedIndex {
		cpy[k] = v
	}
	return cpy
}

// EpisodeList returns an index of xfm episodes in the "shortId" format e.g. xfm-S1E01
func EpisodeList() []string {
	cpy := make([]string, len(episodeOrder))
	for k, v := range episodeOrder {
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
