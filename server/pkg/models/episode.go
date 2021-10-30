package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/util"
	"strings"
	"time"
)

type DialogType string

const (
	DialogTypeUnkown = DialogType("unknown")
	DialogTypeSong   = DialogType("song")
	DialogTypeChat   = DialogType("chat")
	DialogTypeNone   = DialogType("none")
	DialogTypeGap    = DialogType("gap")
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
	ID             string     `json:"id"`
	Position       int64      `json:"pos"`
	OffsetSec      int64      `json:"offset_sec"` // second offset from start of episode
	OffsetInferred bool       `json:"offset_inferred"`
	Type           DialogType `json:"type"`
	Actor          string     `json:"actor"`
	Meta           Metadata   `json:"metadata"`
	Content        string     `json:"content"`
	Notable        bool       `json:"notable"` // note-worthy line of dialog.
}

func (d Dialog) Proto(bestMatch bool) *api.Dialog {
	dialog := &api.Dialog{
		Id:             d.ID,
		Pos:            d.Position,
		OffsetSec:      d.OffsetSec,
		OffsetInferred: d.OffsetInferred,
		Type:           string(d.Type),
		Actor:          d.Actor,
		Content:        d.Content,
		Metadata:       d.Meta.Proto(),
		IsMatchedRow:   bestMatch,
		Notable:        d.Notable,
	}
	return dialog
}

type Transcript struct {
	Publication string    `json:"publication"`
	Series      int32     `json:"series"`
	Episode     int32     `json:"episode"`
	ReleaseDate time.Time `json:"release_date"`
	Incomplete  bool      `json:"incomplete"`

	// additional optional data
	Meta         Metadata   `json:"metadata"`
	Transcript   []Dialog   `json:"transcript"`
	Synopsis     []Synopsis `json:"synopsis"`
	Contributors []string   `json:"contributors"`
}

func (e *Transcript) ID() string {
	return EpisodeID(e)
}

func (e *Transcript) ShortID() string {
	return strings.TrimPrefix(EpisodeID(e), "ep-")
}

func (e *Transcript) ShortProto() *api.ShortTranscript {
	if e == nil {
		return nil
	}
	ep := &api.ShortTranscript{
		Id:          e.ID(),
		Publication: e.Publication,
		Series:      e.Series,
		Episode:     e.Episode,
	}
	return ep
}

func (e *Transcript) Proto(withRawTranscript string, audioURI string) *api.Transcript {
	if e == nil {
		return nil
	}
	ep := &api.Transcript{
		Id:            e.ID(),
		ShortId:       e.ShortID(),
		Publication:   e.Publication,
		Series:        e.Series,
		Episode:       e.Episode,
		Metadata:      e.Meta.Proto(),
		ReleaseDate:   e.ReleaseDate.Format(util.ShortDateFormat),
		Contributors:  e.Contributors,
		Incomplete:    e.Incomplete,
		RawTranscript: withRawTranscript,
		AudioUri:      audioURI,
	}
	for _, d := range e.Transcript {
		ep.Transcript = append(ep.Transcript, d.Proto(false))
	}
	for _, s := range e.Synopsis {
		ep.Synopses = append(ep.Synopses, s.Proto())
	}
	return ep
}

type ShortTranscript struct {
	ID                  string
	Publication         string
	Series              int32
	Episode             int32
	ReleaseDate         time.Time
	TranscriptAvailable bool
}

func (e *ShortTranscript) ShortProto() *api.ShortTranscript {
	if e == nil {
		return nil
	}
	return &api.ShortTranscript{
		Id:                  e.ID,
		Publication:         e.Publication,
		Series:              e.Series,
		Episode:             e.Episode,
		TranscriptAvailable: e.TranscriptAvailable,
	}
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

type Synopsis struct {
	Description string `json:"description"`
	StartPos    int64  `json:"start_pos"`
	EndPos      int64  `json:"end_pos"`
}

func (f Synopsis) Proto() *api.Synopsis {
	return &api.Synopsis{
		Description: f.Description,
		StartPos:    f.StartPos,
		EndPos:      f.EndPos,
	}
}
