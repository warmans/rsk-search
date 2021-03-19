package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func NewSearchService(searchBackend *search.Search, store *ro.Conn, persistentDB *rw.Conn) *SearchService {
	return &SearchService{
		searchBackend: searchBackend,
		staticDB:      store,
		persistentDB:  persistentDB,
	}
}

type SearchService struct {
	searchBackend *search.Search
	staticDB      *ro.Conn
	persistentDB  *rw.Conn
}

func (s *SearchService) RegisterGRPC(server *grpc.Server) {
	api.RegisterSearchServiceServer(server, s)
}

func (s *SearchService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

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
		return nil, err
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
		return nil, err
	}
	return el, nil
}

func (s *SearchService) GetTscriptChunkStats(ctx context.Context, empty *emptypb.Empty) (*api.ChunkStats, error) {
	var stats *models.ChunkStats
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		stats, err = s.GetChunkStats(ctx)
		if stats == nil {
			// empty result
			stats = &models.ChunkStats{}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return stats.Proto(), nil
}

func (s *SearchService) GetTscriptChunk(ctx context.Context, request *api.GetTscriptChunkRequest) (*api.TscriptChunk, error) {
	var chunk *models.Chunk
	var tscriptID string
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		chunk, tscriptID, err = s.GetChunk(ctx, request.Id)
		if err != nil {
			return err
		}
		if chunk == nil {
			return ErrNotFound(request.Id).Err()
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return chunk.Proto(tscriptID), nil
}

func (s *SearchService) ListTscriptChunkSubmissions(ctx context.Context, request *api.ListTscriptChunkSubmissionsRequest) (*api.ChunkSubmissionList, error) {
	panic("implement me")
}

func (s *SearchService) SubmitTscriptChunk(ctx context.Context, request *api.TscriptChunkSubmissionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (s *SearchService) SubmitDialogCorrection(ctx context.Context, request *api.SubmitDialogCorrectionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
