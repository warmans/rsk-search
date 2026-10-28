package models

import (
	"time"

	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/util"
)

type Song struct {
	SpotifyURI      string
	Artist          string
	Title           string
	Album           string
	EpisodeIDs      []string
	Transcribed     []string
	AlbumImageURL   string
	FirstOccurrence *time.Time
}

func (s *Song) Proto() *api.Song {
	return &api.Song{
		SpotifyUri:      s.SpotifyURI,
		Artist:          s.Artist,
		Title:           s.Title,
		Album:           s.Album,
		EpisodeIds:      s.EpisodeIDs,
		Transcribed:     s.Transcribed,
		AlbumImageUrl:   s.AlbumImageURL,
		FirstOccurrence: util.ShortDate(s.FirstOccurrence),
	}
}
