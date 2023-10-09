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

func (s *OauthService) GetAuthURL(ctx context.Context, request *api.GetAuthURLRequest) (*api.AuthURL, error) {

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
	switch request.Provider {
	case "reddit":
		return &api.AuthURL{
			Url: fmt.Sprintf(
				"https://www.reddit.com/api/v1/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=%s&duration=temporary&scope=identity",
				s.oauthCfg.RedditAppID,
				s.csrfCache.NewCSRFToken(returnURL),
				url.QueryEscape(fmt.Sprintf(s.oauthCfg.ReturnURL, request.Provider)),
			),
		}, nil
	case "discord":
		return &api.AuthURL{
			Url: fmt.Sprintf(
				"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify",
				s.oauthCfg.DiscordAppID,
				url.QueryEscape(fmt.Sprintf(s.oauthCfg.ReturnURL, request.Provider)),
			),
		}, nil
	}
	return nil, ErrInvalidRequestField("provider", fmt.Errorf("unknown provider"))
}
