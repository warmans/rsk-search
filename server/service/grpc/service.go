package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strings"
)

func NewSearchService(
	logger *zap.Logger,
	srvCfg config.SearchServiceConfig,
	searchBackend *search.Search,
	store *ro.Conn,
	persistentDB *rw.Conn,
	csrfCache *oauth.CSRFTokenCache,
	auth *jwt.Auth,
	oauthCfg *oauth.Config,
	pledgeClient *pledge.Client,
) *SearchService {
	return &SearchService{
		logger:        logger,
		srvCfg:        srvCfg,
		searchBackend: searchBackend,
		staticDB:      store,
		persistentDB:  persistentDB,
		csrfCache:     csrfCache,
		auth:          auth,
		oauthCfg:      oauthCfg,
		pledgeClient:  pledgeClient,
	}
}

type SearchService struct {
	logger        *zap.Logger
	srvCfg        config.SearchServiceConfig
	searchBackend *search.Search
	staticDB      *ro.Conn
	persistentDB  *rw.Conn
	csrfCache     *oauth.CSRFTokenCache
	auth          *jwt.Auth
	oauthCfg      *oauth.Config
	pledgeClient  *pledge.Client
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
