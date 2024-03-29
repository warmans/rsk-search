package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewUserService(
	logger *zap.Logger,
	persistentDB *rw.Conn,
	auth *jwt.Auth,
) *UserService {
	return &UserService{
		logger:       logger,
		persistentDB: persistentDB,
		auth:         auth,
	}
}

type UserService struct {
	logger       *zap.Logger
	persistentDB *rw.Conn
	auth         *jwt.Auth
}

func (s *UserService) RegisterGRPC(server *grpc.Server) {
	api.RegisterUserServiceServer(server, s)
}

func (s *UserService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *UserService) ListNotifications(ctx context.Context, request *api.ListNotificationsRequest) (*api.NotificationsList, error) {

	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}

	qm, err := NewQueryModifiers(request)
	if err != nil {
		return nil, err
	}

	var notifications models.AuthorNotifications
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var storeErr error
		notifications, storeErr = s.ListAuthorNotifications(ctx, claims.AuthorID, qm)
		return storeErr
	})
	if err != nil {
		return nil, ErrFromStore(err, "")
	}
	return notifications.Proto(), err
}

func (s *UserService) MarkNotificationsRead(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.MarkAllAuthorNotificationsRead(ctx, claims.AuthorID)
	})
	if err != nil {
		return nil, ErrInternal(err)
	}
	return &emptypb.Empty{}, nil
}
