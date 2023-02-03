package search

import (
	"context"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"time"
)

type Searcher interface {
	Search(ctx context.Context, f filter.Filter, page int32) (*api.SearchResultList, error)
	// PredictSearchTerms supports auto-complete for the search bar.
	PredictSearchTerms(ctx context.Context, prefix string, exact bool, numPredictions int32, f filter.Filter) (*api.SearchTermPredictions, error)
	ListTerms(fieldName string, prefix string) (models.FieldValues, error)
}

type DialogDocument struct {
	ID           string     `json:"id"`
	TranscriptID string     `json:"transcript_id"`
	Mapping      string     `json:"mapping"`
	Publication  string     `json:"publication"`
	Series       int64      `json:"series"`
	Episode      int64      `json:"episode"`
	Date         *time.Time `json:"date"`
	Actor        string     `json:"actor"`
	Position     int64      `json:"pos"`
	Content      string     `json:"content"`
	ContentType  string     `json:"type"`
}

func (d DialogDocument) GetNamedField(name string) interface{} {
	switch name {
	case "transcript_id":
		return d.TranscriptID
	case "publication":
		return d.Publication
	case "series":
		return d.Series
	case "episode":
		return d.Episode
	case "date":
		return d.Date
	case "actor":
		return d.Actor
	case "pos":
		return d.Position
	case "content":
		return d.Content
	case "autocomplete":
		return d.Content
	case "type":
		return d.ContentType
	}
	return ""
}
