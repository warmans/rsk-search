package models

import (
	"fmt"
	"github.com/warmans/rsk-search/pkg/util"
	"strconv"
	"strings"
)

func ParseEpID(raw string) (string, int32, int32, error) {

	// this shouldn't be in a episodeID but there are still some around.
	raw = strings.TrimPrefix(raw, "ep-")

	publicationAndSeries := strings.Split(raw, "-")
	if len(publicationAndSeries) != 2 {
		return "", 0, 0, fmt.Errorf("could not parse publication from filename: %s", raw)
	}
	series, episode, err := ExtractSeriesAndEpisode(publicationAndSeries[1])
	if err != nil {
		return "", 0, 0, fmt.Errorf("could not parse series/episode from filename: %s", publicationAndSeries[1])
	}

	return publicationAndSeries[0], series, episode, nil
}

// ExtractSeriesAndEpisode e.g. xfm-S01E02 becomes 1,2,nil
func ExtractSeriesAndEpisode(raw string) (int32, int32, error) {

	raw = strings.ToLower(raw)

	// remove ep-?xfm-
	raw = util.LastSegment(raw, "-")

	raw = strings.TrimPrefix(raw, "s")
	parts := strings.Split(raw, "e")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("name was in wrong format: %s", raw)
	}
	series, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("series %s was not parsable: %w", parts[0], err)
	}
	episode, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("episode %s was not parsable: %w", parts[1], err)
	}
	return int32(series), int32(episode), nil
}

func FormatStandardEpisodeName(series, episode int32) string {
	return fmt.Sprintf("S%dE%02d", series, episode)
}

func ShortEpID(publication string, series, episode int32) string {
	return fmt.Sprintf("%s-%s", publication, FormatStandardEpisodeName(series, episode))
}

// EpID is an episode ID that is prefixed with `ep-`.
func EpID(publication string, series int32, episode int32) string {
	return fmt.Sprintf("ep-%s", ShortEpID(publication, series, episode))
}

func EpIDFromTranscript(ep *Transcript) string {
	return EpID(ep.Publication, ep.Series, ep.Episode)
}

func DialogID(episodeID string, pos int64) string {
	return fmt.Sprintf("%s-%d", episodeID, pos)
}

func IncompleteTranscriptID(t ChunkedTranscript) string {
	return fmt.Sprintf("ts-%s-%s", t.Publication, FormatStandardEpisodeName(t.Series, t.Episode))
}
