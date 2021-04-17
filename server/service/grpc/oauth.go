package grpc

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/oauth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

func NewOauthService(logger *zap.Logger, csrfCache *oauth.CSRFTokenCache, oauthCfg *oauth.Config) *OauthService {
	return &OauthService{
		logger:    logger,
		csrfCache: csrfCache,
		oauthCfg:  oauthCfg,
	}
}

type OauthService struct {
	logger    *zap.Logger
	csrfCache *oauth.CSRFTokenCache
	oauthCfg  *oauth.Config
}

func (s *OauthService) RegisterGRPC(server *grpc.Server) {
	api.RegisterOauthServiceServer(server, s)
}

func (s *OauthService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterOauthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *OauthService) GetRedditAuthURL(ctx context.Context, empty *emptypb.Empty) (*api.RedditAuthURL, error) {

	returnURL := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md["grpcgateway-referer"]) > 0 {
		// we don't want to keep the query or fragment
		if parsed, err := url.Parse(md["grpcgateway-referer"][0]); err == nil {
			parsed.RawQuery = ""
			parsed.RawFragment = ""
			returnURL = parsed.String()
		}
	}
	return &api.RedditAuthURL{
		Url: fmt.Sprintf(
			"https://www.reddit.com/api/v1/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=%s&duration=temporary&scope=identity",
			s.oauthCfg.AppID,
			s.csrfCache.NewCSRFToken(returnURL),
			s.oauthCfg.ReturnURL,
		),
	}, nil
}
