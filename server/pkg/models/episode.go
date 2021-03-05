package models

import (
	"github.com/warmans/rsk-search/gen/api"
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

	MetadataSongArtist = MetadataType("song_artist")
	MetadataSongTrack  = MetadataType("song_track")
	MetadataSongAlbum  = MetadataType("song_album")
)

type Metadata map[MetadataType]string

func (m Metadata) Proto() map[string]string {
	p := map[string]string{}
	for k, v := range m {
		p[string(k)] = v
	}
	return p

}

type Dialog struct {
	ID       string     `json:"id"`
	Position int64      `json:"pos"`
	Type     DialogType `json:"type"`
	Actor    string     `json:"actor"`
	Meta     Metadata   `json:"metadata"`
	Content  string     `json:"content"`
}

func (d Dialog) Proto(bestMatch bool) *api.Dialog {
	return &api.Dialog{
		Id:       d.ID,
		Pos:      d.Position,
		Type:     string(d.Type),
		Actor:    d.Actor,
		Content:  d.Content,
		Metadata: d.Meta.Proto(),
		IsMatchedRow: bestMatch,
	}
}

type Episode struct {
	Publication string    `json:"publication"`
	Series      int32     `json:"series"`
	Episode     int32     `json:"episode"`
	ReleaseDate time.Time `json:"release_date"`

	// additional optional data
	Meta Metadata `json:"metadata"`

	Transcript []Dialog `json:"transcript"`
}

func (e *Episode) Proto() *api.Episode {
	if e == nil {
		return nil
	}
	return &api.Episode{
		Publication: e.Publication,
		Series:      e.Series,
		Episode:     e.Episode,
		Metadata:    e.Meta.Proto(),
	}
}