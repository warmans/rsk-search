package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"google.golang.org/grpc"
)

func NewSearchService() *SearchService {
	return &SearchService{}
}

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, request *api.SearchRequest) (*api.SearchResultList, error) {
	return nil, ErrNotImplemented().Err()
}

func (s *SearchService) RegisterGRPC(server *grpc.Server) {
	api.RegisterSearchServiceServer(server, s)
}

func (s *SearchService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}
