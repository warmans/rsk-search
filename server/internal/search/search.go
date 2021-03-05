package search

import (
	"context"
	"github.com/blevesearch/bleve/v2"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/bleve_query"
	"github.com/warmans/rsk-search/pkg/store"
)

const SearchResultContextLines = 3

func NewSearch(index bleve.Index, db *store.Conn) *Search {
	return &Search{index: index, db: db}
}

type Search struct {
	index bleve.Index
	db    *store.Conn
}

func (s *Search) Search(ctx context.Context, f filter.Filter) (*api.SearchResultList, error) {

	query, err := bleve_query.FilterToQuery(f)
	if err != nil {
		return nil, err
	}

	req := bleve.NewSearchRequest(query)

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
			err := s.db.WithStore(func(s *store.Store) error {
				dialogs, episodeID, err := s.GetDialog(ctx, searchResult.ID, SearchResultContextLines)
				if err != nil {
					return err
				}
				episode, err := s.GetShortEpisode(ctx, episodeID)
				if err != nil {
					return err
				}

				// ep
				result.Episode = episode.Proto()

				// dialogs
				lines := make([]*api.Dialog, len(dialogs))
				for k, d := range dialogs {
					lines[k] = d.Proto(searchResult.ID == d.ID)
				}
				result.Dialogs = append(result.Dialogs, &api.DialogResult{Lines: lines, Score: float32(searchResult.Score)})
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
