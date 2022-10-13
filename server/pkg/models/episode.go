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

type AudioQuality string

func (q AudioQuality) Proto() api.AudioQuality {
	switch q {
	case AudioQualityPoor:
		return api.AudioQuality_POOR
	case AudioQualityAverage:
		return api.AudioQuality_AVERAGE
	case AudioQualityGood:
		return api.AudioQuality_GOOD
	}
	return api.AudioQuality_UNKNOWN
}

const (
	AudioQualityUnknown = AudioQuality("")
	AudioQualityPoor    = AudioQuality("poor")
	AudioQualityAverage = AudioQuality("average")
	AudioQualityGood    = AudioQuality("good")
)

type MetadataType string

const (
	MetadataTypePilkipediaURL     = MetadataType("pilkipedia_url")
	MetadataTypeSpotifyPreviewURL = MetadataType("spotify_player_url")
	MetadataTypeSpotifyURI        = MetadataType("spotify_uri")
	MetadataTypeDurationMs        = MetadataType("duration_ms")
	CoverArtURL                   = MetadataType("cover_art_url")

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
		Pos:            int32(d.Position),
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
	Publication string `json:"publication"`
	Series      int32  `json:"series"`
	Episode     int32  `json:"episode"`
	// some episodes don't really have a proper series/episode and need to be identified by a name e.g. Radio 2 special
	Name        string     `json:"name"`
	Summary     string     `json:"summary"`
	Version     string     `json:"version"` // SemVer
	ReleaseDate *time.Time `json:"release_date"`
	// is the episode missing some sections of transcript?
	Incomplete bool `json:"incomplete"`
	// is the episode a "clip show"?
	Bestof bool `json:"bestof"`
	// is the episode a "one off" special type episode?
	Special bool `json:"special"`
	// Force an episode into the locked state. If false it will not override any other locking.
	Locked bool `json:"locked"`

	OffsetAccuracy int32 `json:"offset_accuracy"`

	AudioQuality AudioQuality `json:"audio_quality"`

	// additional optional data
	Meta         Metadata   `json:"metadata"`
	Transcript   []Dialog   `json:"transcript"`
	Synopsis     []Synopsis `json:"synopsis"`
	Contributors []string   `json:"contributors"`
	Trivia       []Trivia   `json:"trivia"`
}

func (e *Transcript) ID() string {
	return EpIDFromTranscript(e)
}

func (e *Transcript) ShortID() string {
	return strings.TrimPrefix(EpIDFromTranscript(e), "ep-")
}

func (e Transcript) Actors() []string {
	actorMap := map[string]struct{}{}
	for _, v := range e.Transcript {
		if strings.TrimSpace(v.Actor) == "" {
			continue
		}
		// this is almost always "ricky and steve"
		if strings.Contains(v.Actor, " and ") || strings.Contains(v.Actor, " & ") {
			continue
		}
		actorMap[strings.TrimSpace(v.Actor)] = struct{}{}
	}
	actorList := []string{}
	for k := range actorMap {
		actorList = append(actorList, k)
	}
	return actorList
}

func (e *Transcript) ShortProto(audioURI string) *api.ShortTranscript {
	if e == nil {
		return nil
	}
	ep := &api.ShortTranscript{
		Id:                  e.ID(),
		Publication:         e.Publication,
		Series:              e.Series,
		Episode:             e.Episode,
		TranscriptAvailable: len(e.Transcript) > 0,
		Incomplete:          e.Incomplete,
		ReleaseDate:         util.ShortDate(e.ReleaseDate),
		Summary:             e.Summary,
		Synopsis:            make([]*api.Synopsis, len(e.Synopsis)),
		TriviaAvailable:     len(e.Trivia) > 0,
		Actors:              e.Actors(),
		ShortId:             e.ShortID(),
		AudioUri:            audioURI,
		OffsetAccuracyPcnt:  e.OffsetAccuracy,
		Name:                e.Name,
		Version:             e.Version,
		Metadata:            e.Meta.Proto(),
		Bestof:              e.Bestof,
		Special:             e.Special,
		AudioQuality:        e.AudioQuality.Proto(),
	}
	for k, s := range e.Synopsis {
		ep.Synopsis[k] = s.Proto()
	}
	return ep
}

func (e *Transcript) Proto(withRawTranscript string, audioURI string, forceLockedOn bool) *api.Transcript {
	if e == nil {
		return nil
	}
	ep := &api.Transcript{
		Id:                 e.ID(),
		ShortId:            e.ShortID(),
		Publication:        e.Publication,
		Series:             e.Series,
		Episode:            e.Episode,
		Metadata:           e.Meta.Proto(),
		ReleaseDate:        util.ShortDate(e.ReleaseDate),
		Contributors:       e.Contributors,
		Incomplete:         e.Incomplete,
		RawTranscript:      withRawTranscript,
		AudioUri:           audioURI,
		Actors:             e.Actors(),
		OffsetAccuracyPcnt: e.OffsetAccuracy,
		Name:               e.Name,
		Version:            e.Version,
		Bestof:             e.Bestof,
		Special:            e.Special,
		Locked:             e.Locked || forceLockedOn,
		Summary:            e.Summary,
		AudioQuality:       e.AudioQuality.Proto(),
	}
	for _, d := range e.Transcript {
		ep.Transcript = append(ep.Transcript, d.Proto(false))
	}
	for _, s := range e.Synopsis {
		ep.Synopses = append(ep.Synopses, s.Proto())
	}
	for _, t := range e.Trivia {
		ep.Trivia = append(ep.Trivia, t.Proto())
	}
	return ep
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
		StartPos:    int32(f.StartPos),
		EndPos:      int32(f.EndPos),
	}
}

type Trivia struct {
	Description string `json:"description"`
	StartPos    int64  `json:"start_pos"`
	EndPos      int64  `json:"end_pos"`
}

func (f Trivia) Proto() *api.Trivia {
	return &api.Trivia{
		Description: f.Description,
		StartPos:    int32(f.StartPos),
		EndPos:      int32(f.EndPos),
	}
}
