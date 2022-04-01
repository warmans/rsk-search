package v1

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/bleve_query"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/ro"
)

const ResultContextLines = 3
const PageSize = 10

func NewSearch(index bleve.Index, readOnlyDB *ro.Conn, episodeCache *data.EpisodeCache, audioUriPattern string) search.Searcher {
	return &Search{index: index, readOnlyDB: readOnlyDB, episodeCache: episodeCache, audioUriPattern: audioUriPattern}
}

type Search struct {
	index        bleve.Index
	readOnlyDB   *ro.Conn
	episodeCache *data.EpisodeCache
	audioUriPattern string
}

func (s *Search) Search(ctx context.Context, f filter.Filter, page int32) (*api.SearchResultList, error) {

	query, err := bleve_query.FilterToQuery(f)
	if err != nil {
		return nil, err
	}

	req := bleve.NewSearchRequest(query)
	req.Size = PageSize
	req.From = PageSize * int(page)

	result, err := s.index.Search(req)
	if err != nil {
		return nil, err
	}

	res := &api.SearchResultList{
		ResultCount: int32(result.Total),
		Results:     []*api.SearchResult{},
	}

	//todo: need to change this to allow multiple results from the same episode to go under
	// an single result.
	// the processing is a bit tricky though as it will be harder to return the correct limit

	if len(result.Hits) > 0 {
		res.Results = []*api.SearchResult{}
		for _, searchResult := range result.Hits {
			result := &api.SearchResult{
				Dialogs: []*api.DialogResult{},
			}
			err := s.readOnlyDB.WithStore(func(store *ro.Store) error {
				dialogs, episodeID, err := store.GetDialogWithContext(ctx, searchResult.ID, ResultContextLines)
				if err != nil {
					return err
				}

				ep, err := s.episodeCache.GetEpisode(episodeID)
				if err != nil {
					return err
				}
				result.Episode = ep.ShortProto(fmt.Sprintf(s.audioUriPattern, ep.ShortID()))

				// dialogs
				lines := make([]*api.Dialog, len(dialogs))
				for k, d := range dialogs {
					lines[k] = d.Proto(searchResult.ID == d.ID)
				}
				result.Dialogs = append(result.Dialogs, &api.DialogResult{Transcript: lines, Score: float32(searchResult.Score)})
				return nil
			})
			if err != nil {
				return nil, err
			}
			res.Results = append(res.Results, result)
		}
	}

	return res, nil
}

func (s *Search) ListTerms(fieldName string, prefix string) (models.FieldValues, error) {
	dct, err := s.index.FieldDictPrefix(fieldName, []byte(prefix))
	if err != nil {
		return nil, err
	}
	defer dct.Close()

	terms := models.FieldValues{}
	for {
		entry, err := dct.Next()
		if err != nil {
			return nil, err
		}
		if entry == nil {
			break
		}
		terms = append(terms, models.FieldValue{Value: entry.Term, Count: int32(entry.Count)})
	}
	return terms, nil
}
