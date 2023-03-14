package grpc

import (
	"context"
	"fmt"
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
		return nil, ErrUnauthorized("Only approvers may delete incomplete transcripts")
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
		return nil, ErrUnauthorized("Only approvers may create new transcript imports")
	}
	var tscriptImport *models.TscriptImport

	err = s.persistentDB.WithStore(func(store *rw.Store) error {
		var err error
		existingTscripts, err := store.ListTscripts(ctx)
		if err != nil {
			return err
		}
		for _, v := range existingTscripts {
			if request.Epid == v.AsEpisode().ID() {
				return ErrFailedPrecondition(fmt.Sprintf("Epid %s already exists as a tscript", request.Epid))
			}
		}

		tscriptImport, err = store.CreateTscriptImport(ctx, &models.TscriptImportCreate{
			EpID:   request.Epid,
			EpName: request.Epname,
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
		return nil, ErrFromStore(err, "")
	}
	return tscriptImport.Proto(), nil
}

func (s *AdminService) ListTscriptImports(ctx context.Context, request *api.ListTscriptImportsRequest) (*api.TscriptImportList, error) {
	claims, err := s.getClaims(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.Approver {
		return nil, ErrUnauthorized("Only approvers may view transcript imports")
	}
	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}

	var res []*models.TscriptImport
	if err := s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		res, err = s.ListTscriptImports(ctx, qm)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	out := make([]*api.TscriptImport, len(res))
	for k, v := range res {
		out[k] = v.Proto()
	}

	return &api.TscriptImportList{Imports: out}, nil
}

func (s *AdminService) getClaims(ctx context.Context) (*jwt.Claims, error) {
	token := jwt.ExtractTokenFromRequestContext(ctx)
	if token == "" {
		return nil, ErrUnauthorized("no token provided")
	}
	claims, err := s.auth.VerifyToken(token)
	if err != nil {
		return nil, ErrUnauthorized(err.Error())
	}
	return claims, nil
}
