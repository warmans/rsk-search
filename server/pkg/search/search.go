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
	ListTerms(fieldName string, prefix string) (models.FieldValues, error)
}

type DialogDocument struct {
	ID          string   `json:"id"`
	Mapping     string   `json:"mapping"`
	Publication string   `json:"publication"`
	Series      int32    `json:"series"`
	Episode     int32    `json:"episode"`
	Date        time.Time   `json:"date"`
	Actor       string   `json:"actor"`
	Position    int64    `json:"pos"`
	Content     string   `json:"content"`
	ContentType string   `json:"type"`
	Tags        []string `json:"tags"`
}
