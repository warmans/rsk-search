package models

import (
	"fmt"
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
	MetadataTypePublication = MetadataType("publication")
	MetadataTypeSeries      = MetadataType("series")
	MetadataTypeDate        = MetadataType("date")
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
	Source     string     `json:"source"`
	Meta       []Metadata `json:"metadata"`
	Transcript []Dialog   `json:"transcript"`
}

func (e Episode) MetaValue(t MetadataType) string {
	for _, m := range e.Meta {
		if m.Type == t {
			return m.Value
		}
	}
	return "na"
}

func (e Episode) CanonicalName() string {
	date := "na"
	if rawDate := e.MetaValue(MetadataTypeDate); rawDate != "" {
		t, err := time.Parse(time.RFC3339, rawDate)
		if err == nil {
			date = t.Format("Jan-02-2006")
		}
	}

	return fmt.Sprintf("%s-%s-%s", e.MetaValue(MetadataTypePublication), e.MetaValue(MetadataTypeSeries), date)
}
