package v2

import (
	"context"
	"fmt"
	"github.com/blugelabs/bluge"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/bluge_query"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/ro"
)

const ResultContextLines = 3
const PageSize = 10

func NewSearch(index *bluge.Reader, readOnlyDB *ro.Conn, episodeCache *data.EpisodeCache) search.Searcher {
	return &Search{index: index, readOnlyDB: readOnlyDB, episodeCache: episodeCache}
}

type Search struct {
	index        *bluge.Reader
	readOnlyDB   *ro.Conn
	episodeCache *data.EpisodeCache
}

func (s *Search) Search(ctx context.Context, f filter.Filter, page int32) (*api.SearchResultList, error) {

	query, err := bluge_query.FilterToQuery(f)
	if err != nil {
		return nil, err
	}

	req := bluge.NewTopNSearch(PageSize, query).SetFrom(PageSize * int(page)).WithStandardAggregations()

	dmi, err := s.index.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &api.SearchResultList{
		ResultCount: int32(dmi.Aggregations().Count()),
		Results:     []*api.SearchResult{},
	}

	res.Results = []*api.SearchResult{}

	next, err := dmi.Next()
	for err == nil && next != nil {
		result := &api.SearchResult{
			Dialogs: []*api.DialogResult{},
		}
		err := s.readOnlyDB.WithStore(func(store *ro.Store) error {

			var innerErr error
			err = next.VisitStoredFields(func(field string, value []byte) bool {
				if innerErr != nil {
					return false
				}
				if field == "_id" {
					dialogs, episodeID, err := store.GetDialogWithContext(ctx, string(value), ResultContextLines)
					if err != nil {
						innerErr = err
						return false
					}
					ep, err := s.episodeCache.GetEpisode(episodeID)
					if err != nil {
						innerErr = err
						return false
					}
					result.Episode = ep.ShortProto()

					// dialogs
					lines := make([]*api.Dialog, len(dialogs))
					for k, d := range dialogs {
						lines[k] = d.Proto(string(value) == d.ID)
					}
					result.Dialogs = append(result.Dialogs, &api.DialogResult{Lines: lines, Score: float32(next.Score)})
				}
				return true
			})
			if err != nil {
				return fmt.Errorf("error accessing stored fields: %v", err)
			}
			if innerErr != nil {
				return fmt.Errorf("error processing result: %v", err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		res.Results = append(res.Results, result)

		next, err = dmi.Next()
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Search) ListTerms(fieldName string, prefix string) (models.FieldValues, error) {

	terms := models.FieldValues{}

	fieldDict, err := s.index.DictionaryIterator(fieldName, nil, []byte(prefix), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := fieldDict.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	tfd, err := fieldDict.Next()
	for err == nil && tfd != nil {
		terms = append(terms, models.FieldValue{Value: tfd.Term(), Count: int32(tfd.Count())})
		if len(terms) > 500 {
			return nil, fmt.Errorf("too many terms for field '%s' returned (prefix: %s)", fieldName, prefix)
		}
		tfd, err = fieldDict.Next()
	}

	return terms, nil
}
