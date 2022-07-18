package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/queue"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewAdminService(
	logger *zap.Logger,
	taskQueue queue.ImportPipeline,
	auth *jwt.Auth,
	persistentDB *rw.Conn,
) *AdminService {
	return &AdminService{
		logger:       logger,
		taskQueue:    taskQueue,
		auth:         auth,
		persistentDB: persistentDB,
	}
}

type AdminService struct {
	logger       *zap.Logger
	taskQueue    queue.ImportPipeline
	auth         *jwt.Auth
	persistentDB *rw.Conn
}

func (s *AdminService) RegisterGRPC(server *grpc.Server) {
	api.RegisterAdminServiceServer(server, s)
}

func (s *AdminService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *AdminService) DeleteTscript(ctx context.Context, request *api.DeleteTscriptRequest) (*emptypb.Empty, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.Approver {
		return nil, ErrUnauthorized("Only approvers may delete incomplete transcripts").Err()
	}

	//todo: check there are no outstanding contributions

	if err := s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.DeleteTscript(ctx, request.Id)
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdminService) CreateTscriptImport(ctx context.Context, request *api.CreateTscriptImportRequest) (*api.TscriptImport, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.Approver {
		return nil, ErrUnauthorized("Only approvers may create new incomplete transcripts").Err()
	}
	var tscriptImport *models.TscriptImport

	err = s.persistentDB.WithStore(func(store *rw.Store) error {
		var err error
		tscriptImport, err = store.CreateTscriptImport(ctx, &models.TscriptImportCreate{
			EpID:   request.Epid,
			Mp3URI: request.Mp3Uri,
		})
		if err != nil {
			return err
		}

		// enqueue the import
		if err := s.taskQueue.StartNewImport(ctx, tscriptImport); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrFromStore(err, "").Err()
	}
	return tscriptImport.Proto(), nil
}

func (s *AdminService) getClaims(ctx context.Context) (*jwt.Claims, error) {
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
