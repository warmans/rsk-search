package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/bleve_query"
	"github.com/warmans/rsk-search/pkg/store"
)

func NewSearch(index bleve.Index, db *store.Conn) *Search {
	return &Search{index: index, db: db}
}

type Search struct {
	index bleve.Index
	db    *store.Conn
}

func (s *Search) Search(f filter.Filter) (*api.SearchResultList, error) {

	query, err := bleve_query.FilterToQuery(f)
	if err != nil {
		return nil, err
	}

	req := bleve.NewSearchRequest(query)
	req.Highlight = bleve.NewHighlightWithStyle("ansi")

	result, err := s.index.Search(req)
	if err != nil {
		return nil, err
	}

	res := &api.SearchResultList{
		ResultCount: int32(result.Total),
		Results:     []*api.SearchResult{},
	}

	if len(result.Hits) > 0 {
		res.Results = []*api.SearchResult{}
		for _, v := range result.Hits {

			result := &api.SearchResult{}
			//todo: fetch context

			lines := []*api.Dialog{{Id: v.ID}}
			result.Dialogs = append(result.Dialogs, &api.DialogResult{Lines: lines, Episode: nil}) //todo

			res.Results = append(res.Results, result)
		}
	}

	return res, nil
}
