package grpc

import (
	"context"
	"fmt"
	"github.com/warmans/rsk-search/gen/api"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

func (s *SearchService) GetRedditAuthURL(ctx context.Context, empty *emptypb.Empty) (*api.RedditAuthURL, error) {

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

