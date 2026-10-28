package meta

import (
	"testing"
	"time"
)

func TestExtractSortedPopulatesFirstOccurrence(t *testing.T) {
	songs := GetSongMeta().ExtractSorted()
	if len(songs) == 0 {
		t.Fatal("no songs in embedded metadata")
	}
	populated := 0
	for _, song := range songs {
		if song.FirstOccurrence == nil {
			continue
		}
		populated++
		// the first occurrence must not be later than any known episode date for the song
		for _, id := range song.EpisodeIDs {
			if date, ok := episodeDatesByID[id]; ok && date.Before(*song.FirstOccurrence) {
				t.Errorf("first occurrence %s for %s is later than episode %s on %s", song.FirstOccurrence.Format(time.DateOnly), song.Track.Name, id, date.Format(time.DateOnly))
			}
		}
	}
	if populated == 0 {
		t.Fatal("no songs had a first occurrence populated")
	}
}
