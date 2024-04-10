package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewCommunityService(
	logger *zap.Logger,
	staticDB *ro.Conn,
) *CommunityService {
	return &CommunityService{
		logger:   logger,
		staticDB: staticDB,
	}
}

type CommunityService struct {
	logger   *zap.Logger
	staticDB *ro.Conn
}

func (s *CommunityService) RegisterGRPC(server *grpc.Server) {
	api.RegisterCommunityServiceServer(server, s)
}

func (s *CommunityService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterCommunityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *CommunityService) ListProjects(ctx context.Context, request *api.ListCommunityProjectsRequest) (*api.CommunityProjectList, error) {
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}
	var projectList *api.CommunityProjectList
	err = s.staticDB.WithStore(func(s *ro.Store) error {
		projects, totalCount, err := s.ListCommunityProjects(ctx, qm)
		if err != nil {
			return err
		}
		projectList = projects.Proto(totalCount)
		return nil
	})
	if err != nil {
		return nil, ErrInternal(err)
	}
	return projectList, nil
}
