package grpc

import (
	"context"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (s *SearchService) GetSearchMetadata(ctx context.Context, empty *emptypb.Empty) (*api.SearchMetadata, error) {
	return meta.GetSearchMeta().Proto(), nil
}

func (s *SearchService) ListFieldValues(ctx context.Context, request *api.ListFieldValuesRequest) (*api.FieldValueList, error) {
	vals, err := s.searchBackend.ListTerms(request.Field, request.Prefix)
	if err != nil {
		return nil, err
	}
	return &api.FieldValueList{Values: vals.Proto()}, nil
}

func (s *SearchService) Search(ctx context.Context, request *api.SearchRequest) (*api.SearchResultList, error) {
	f, err := filter.Parse(request.Query)
	if err != nil {
		return nil, ErrInvalidRequestField("query", err.Error()).Err()
	}
	if err := checkWhy(f); err != nil {
		return nil, err
	}
	return s.searchBackend.Search(ctx, f, request.Page)
}

func (s *SearchService) GetEpisode(ctx context.Context, request *api.GetEpisodeRequest) (*api.Episode, error) {
	var ep *models.Episode
	err := s.staticDB.WithStore(func(s *ro.Store) error {
		var err error
		ep, err = s.GetEpisode(ctx, request.Id)
		if err != nil {
			return err
		}
		if ep == nil {
			return ErrNotFound(request.Id).Err()
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, request.Id).Err()
	}
	return ep.Proto(), nil
}

func (s *SearchService) ListEpisodes(ctx context.Context, request *api.ListEpisodesRequest) (*api.EpisodeList, error) {
	el := &api.EpisodeList{
		Episodes: []*api.ShortEpisode{},
	}
	err := s.staticDB.WithStore(func(s *ro.Store) error {
		eps, err := s.ListEpisodes(ctx)
		if err != nil {
			return err
		}
		for _, e := range eps {
			el.Episodes = append(el.Episodes, e.ShortProto())
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return el, nil
}

func checkWhy(f filter.Filter) error {
	visitor := filter.NewExtractFilterVisitor(f)
	filters, err := visitor.ExtractCompFilters("content")
	if err != nil {
		return nil // don't fail because of this stupid feature
	}
	for _, v := range filters {
		if strings.TrimSpace(strings.Trim(v.Value.String(), `"?`)) == "why" {
			return ErrServerConfused().Err()
		}
	}
	return nil
}
