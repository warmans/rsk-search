package meta

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"io/fs"
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
	defer func(epIndexFile fs.File) {
		_ = epIndexFile.Close()
	}(epIndexFile)
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
	defer func(orderedEpsFile fs.File) {
		_ = orderedEpsFile.Close()
	}(orderedEpsFile)
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
	copy(cpy, episodeOrder)
	return cpy
}

func PreviousEpisode(epid string) (string, bool) {
	for k, v := range episodeOrder {
		if v == epid {
			if k == 0 {
				return "", false
			}
			return episodeOrder[k-1], true
		}
	}
	return "", false
}

// IsValidEpisodeID returns true if the ID appears to be valid based on the format, but doesn't mean it will exist.
func IsValidEpisodeID(id string) bool {
	if _, _, _, err := models.ParseEpID(id); err != nil {
		return false
	}
	return true
}

// IsKnownEpisodeID returns true if this episode is the data mapping metadata.
func IsKnownEpisodeID(id string) bool {
	_, ok := episodeIDMap[id]
	return ok
}
