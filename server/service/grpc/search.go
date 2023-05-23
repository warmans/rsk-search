package grpc

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func NewSearchService(
	logger *zap.Logger,
	srvCfg config.SearchServiceConfig,
	searchBackend search.Searcher,
	staticDB *ro.Conn,
	persistentDB *rw.Conn,
	auth *jwt.Auth,
	episodeCache *data.EpisodeCache,

) *SearchService {
	return &SearchService{
		logger:        logger,
		srvCfg:        srvCfg,
		searchBackend: searchBackend,
		staticDB:      staticDB,
		persistentDB:  persistentDB,
		auth:          auth,
		episodeCache:  episodeCache,
	}
}

type SearchService struct {
	logger        *zap.Logger
	srvCfg        config.SearchServiceConfig
	searchBackend search.Searcher
	staticDB      *ro.Conn
	persistentDB  *rw.Conn
	auth          *jwt.Auth
	episodeCache  *data.EpisodeCache
}

func (s *SearchService) RegisterGRPC(server *grpc.Server) {
	api.RegisterSearchServiceServer(server, s)
}

func (s *SearchService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *SearchService) GetMetadata(ctx context.Context, empty *emptypb.Empty) (*api.Metadata, error) {
	return &api.Metadata{
		SearchFields:    meta.GetSearchMeta().Proto(),
		EpisodeShortIds: meta.EpisodeList(),
	}, nil
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
		return nil, ErrInvalidRequestField("query", err, fmt.Sprintf("query: %s", request.Query))
	}

	if err := checkWhy(f); err != nil {
		return nil, err
	}
	return s.searchBackend.Search(ctx, f, request.Page)
}

func (s *SearchService) PredictSearchTerm(ctx context.Context, request *api.PredictSearchTermRequest) (*api.SearchTermPredictions, error) {

	var f filter.Filter
	if request.Query != "" {
		var err error
		f, err = filter.Parse(request.Query)
		if err != nil {
			return nil, ErrInvalidRequestField("query", err)
		}
	}

	maxPredictions := int32(100)
	if request.MaxPredictions < maxPredictions {
		maxPredictions = request.MaxPredictions
	}
	if strings.TrimSpace(request.Prefix) == "" {
		return &api.SearchTermPredictions{}, nil
	}
	return s.searchBackend.PredictSearchTerms(ctx, request.Prefix, request.Exact, maxPredictions, f, request.Regexp)
}

func (s *SearchService) ListChangelogs(ctx context.Context, request *api.ListChangelogsRequest) (*api.ChangelogList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}
	result := &api.ChangelogList{
		Changelogs: make([]*api.Changelog, 0),
	}
	err = s.staticDB.WithStore(func(s *ro.Store) error {
		changelogs, err := s.ListChangelogs(ctx, qm)
		if err != nil {
			return err
		}
		for _, v := range changelogs {
			result.Changelogs = append(result.Changelogs, v.Proto())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func checkWhy(f filter.Filter) error {
	if f == nil {
		return nil
	}
	visitor := filter.NewExtractFilterVisitor(f)
	filters, err := visitor.ExtractCompFilters("content")
	if err != nil {
		return nil // don't fail because of this stupid feature
	}
	for _, v := range filters {
		if strings.TrimSpace(strings.Trim(v.Value.String(), `"?`)) == "why" {
			return ErrServerConfused()
		}
	}
	return nil
}
