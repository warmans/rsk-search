package meta

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"
)

//go:embed data/xfm-episode-date-map.json
var xfmEpisodeDateIndex embed.FS

//go:embed data/xfm-episode-order.json
var xfmOrderedEpisodes embed.FS

var parsedIndex = map[string]string{}
var episodeIDMap = map[string]struct{}{}
var episodeOrder = []string{}

func init() {
	epIndexFile, err := xfmEpisodeDateIndex.Open("data/xfm-episode-date-map.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer epIndexFile.Close()
	dec := json.NewDecoder(epIndexFile)
	if err := dec.Decode(&parsedIndex); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
	for _, v := range parsedIndex {
		episodeIDMap[fmt.Sprintf("ep-%s-%s", PublicationXFM, v)] = struct{}{}
	}

	orderedEpsFile, err := xfmOrderedEpisodes.Open("data/xfm-episode-order.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer orderedEpsFile.Close()
	if err := json.NewDecoder(orderedEpsFile).Decode(&episodeOrder); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
}

func XfmEpisodeDates() map[string]string {
	cpy := map[string]string{}
	for k, v := range parsedIndex {
		cpy[k] = v
	}
	return cpy
}

func XfmEpisodeList() []string {
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

func XfmEpisodeDateToName(date time.Time) (string, error) {
	dateStr := date.Format(time.RFC3339)
	name, ok := parsedIndex[dateStr]
	if !ok {
		return "", fmt.Errorf("date %s did not map to any known episode", dateStr)

	}
	return name, nil
}
