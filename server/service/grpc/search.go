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
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/transcript"
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
	store *ro.Conn,
	auth *jwt.Auth,
	episodeCache *data.EpisodeCache,
) *SearchService {
	return &SearchService{
		logger:        logger,
		srvCfg:        srvCfg,
		searchBackend: searchBackend,
		staticDB:      store,
		auth:          auth,
		episodeCache:  episodeCache,
	}
}

type SearchService struct {
	logger        *zap.Logger
	srvCfg        config.SearchServiceConfig
	searchBackend search.Searcher
	staticDB      *ro.Conn
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

func (s *SearchService) getClaims(ctx context.Context) (*jwt.Claims, error) {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return nil, ErrUnauthorized("no token provided").Err()
	}
	claims, err := s.auth.VerifyToken(token)
	if err != nil {
		return nil, ErrUnauthorized(err.Error()).Err()
	}
	return claims, nil
}

func NewQueryModifiers(req interface{}) (*common.QueryModifier, error) {
	q := common.Q()
	if p, ok := req.(common.Pager); ok {
		q.Apply(common.WithPaging(p.GetPageSize(), p.GetPage()))
	}
	if p, ok := req.(common.Sorter); ok {
		if p.GetSortField() != "" {
			givenDirection := common.SortDirection(strings.ToUpper(p.GetSortDirection()))
			if givenDirection != common.SortAsc && givenDirection != common.SortDesc {
				return nil, ErrInvalidRequestField("sort_direction", "Must be 'asc' or 'desc'").Err()
			}
			q.Apply(common.WithSorting(p.GetSortField(), givenDirection))
		}
	}
	if p, ok := req.(common.Filterer); ok {
		if strings.TrimSpace(p.GetFilter()) != "" {
			fil, err := filter.Parse(p.GetFilter())
			if err != nil {
				return nil, ErrInvalidRequestField("filter", err.Error()).Err()
			}
			q.Apply(common.WithFilter(fil))
		}
	}
	return q, nil
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

func (s *SearchService) GetTranscript(_ context.Context, request *api.GetTranscriptRequest) (*api.Transcript, error) {
	ep, err := s.episodeCache.GetEpisode(request.Epid)
	if err == data.ErrNotFound || ep == nil {
		return nil, ErrNotFound(request.Epid).Err()
	}
	var rawTranscript string
	if request.WithRaw {
		var err error
		rawTranscript, err = transcript.Export(ep.Transcript, ep.Synopsis, ep.Trivia)
		if err != nil {
			return nil, ErrInternal(err).Err()
		}
	}
	return ep.Proto(rawTranscript, fmt.Sprintf(s.srvCfg.AudioUriPattern, ep.ShortID())), nil
}

func (s *SearchService) ListTranscripts(_ context.Context, _ *api.ListTranscriptsRequest) (*api.TranscriptList, error) {
	el := &api.TranscriptList{
		Episodes: []*api.ShortTranscript{},
	}
	for _, ep := range s.episodeCache.ListEpisodes() {
		el.Episodes = append(el.Episodes, ep.ShortProto())
	}
	return el, nil
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
