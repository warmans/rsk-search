package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/util"
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

	// content tokens mapped to tags
	// e.g. Foo! (text) => foo (tag)
	ContentTags map[string]string `json:"content_tags"`
}

func (d Dialog) Proto(bestMatch bool) *api.Dialog {
	dialog := &api.Dialog{
		Id:           d.ID,
		Pos:          d.Position,
		Type:         string(d.Type),
		Actor:        d.Actor,
		Content:      d.Content,
		Metadata:     d.Meta.Proto(),
		IsMatchedRow: bestMatch,
		ContentTags:  make(map[string]*api.Tag),
	}
	for text, tagName := range d.ContentTags {
		tag := meta.GetTag(tagName)
		if tag != nil {
			dialog.ContentTags[text] = &api.Tag{Name: tagName, Kind: tag.Kind}
		}
	}
	return dialog
}

type Episode struct {
	Publication string    `json:"publication"`
	Series      int32     `json:"series"`
	Episode     int32     `json:"episode"`
	ReleaseDate time.Time `json:"release_date"`

	// additional optional data
	Meta       Metadata `json:"metadata"`
	Transcript []Dialog `json:"transcript"`
	Tags       []string `json:"tags"`
}

func (e *Episode) ShortProto() *api.ShortEpisode {
	if e == nil {
		return nil
	}
	ep := &api.ShortEpisode{
		Id:          e.ID(),
		Publication: e.Publication,
		Series:      e.Series,
		Episode:     e.Episode,
		Metadata:    e.Meta.Proto(),
	}
	for _, tn := range e.Tags {
		tag := meta.GetTag(tn)
		if tag != nil {
			ep.Tags = append(ep.Tags, &api.Tag{Name: tn, Kind: tag.Kind})
		}
	}
	return ep
}

func (e *Episode) Proto() *api.Episode {
	if e == nil {
		return nil
	}
	ep := &api.Episode{
		Id:          e.ID(),
		Publication: e.Publication,
		Series:      e.Series,
		Episode:     e.Episode,
		Metadata:    e.Meta.Proto(),
		ReleaseDate: e.ReleaseDate.Format(util.ShortDateFormat),
	}
	for _, tn := range e.Tags {
		tag := meta.GetTag(tn)
		if tag != nil {
			ep.Tags = append(ep.Tags, &api.Tag{Name: tn, Kind: tag.Kind})
		}
	}
	for _, d := range e.Transcript {
		ep.Transcript = append(ep.Transcript, d.Proto(false))
	}
	return ep
}

func (e *Episode) ID() string {
	return EpisodeID(e)
}

type DialogTags struct {
	DialogID string            `json:"dialog_id"`
	Tags     map[string]string `json:"tags"`
}

type FieldValues []FieldValue

func (f FieldValues) Proto() []*api.FieldValue {
	if f == nil {
		return make([]*api.FieldValue, 0)
	}
	out := make([]*api.FieldValue, len(f))
	for k, v := range f {
		out[k] = v.Proto()
	}
	return out
}

type FieldValue struct {
	Value string
	Count int32
}

func (f FieldValue) Proto() *api.FieldValue {
	return &api.FieldValue{Count: f.Count, Value: f.Value}
}

