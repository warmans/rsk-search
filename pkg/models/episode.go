package models

import (
	"time"
)

type DialogType string

const (
	DialogTypeUnkown = DialogType("unknown")
	DialogTypeSong   = DialogType("song")
	DialogTypeChat   = DialogType("chat")
)

type MetadataType string

const (
	MetadataTypePilkipediaURL     = MetadataType("pilkipedia_url")
	MetadataTypeSpotifyPreviewURL = MetadataType("spotify_player_url")
	MetadataTypeSpotifyURI        = MetadataType("spotify_uri")
	MetadataTypeDurationMs        = MetadataType("duration_ms")
)

type Dialog struct {
	ID       string     `json:"id"`
	Position int64      `json:"pos"`
	Type     DialogType `json:"type"`
	Actor    string     `json:"actor"`
	Content  string     `json:"content"`
}

type Metadata struct {
	Type  MetadataType `json:"type"`
	Value string       `json:"value"`
}

type Episode struct {
	Publication string    `json:"publication"`
	Series      int32     `json:"series"`
	Episode     int32     `json:"episode"`
	ReleaseDate time.Time `json:"release_date"`

	// additional optional data
	Meta map[MetadataType]string `json:"metadata"`

	Transcript []Dialog `json:"transcript"`
}
