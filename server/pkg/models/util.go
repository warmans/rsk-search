package models

import (
	"fmt"
	"strconv"
	"strings"
)

// e.g. S01E02 becomes 1,2,nil
func ParseStandardEpisodeName(raw string) (int32, int32, error) {
	raw = strings.TrimPrefix(raw, "S")
	parts := strings.Split(raw, "E")
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

func EpisodeID(ep *Transcript) string {
	return fmt.Sprintf("ep-%s-%s", ep.Publication, FormatStandardEpisodeName(ep.Series, ep.Episode))
}

func DialogID(episodeID string, pos int64) string {
	return fmt.Sprintf("%s-%d", episodeID, pos)
}

func IncompleteTranscriptID(t Tscript) string {
	return fmt.Sprintf("ts-%s-%s", t.Publication, FormatStandardEpisodeName(t.Series, t.Episode))
}
