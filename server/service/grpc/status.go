package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/quota"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewStatusService(
	logger *zap.Logger,
	persistentDB *rw.Conn,
) *StatusService {
	return &StatusService{
		logger:       logger,
		persistentDB: persistentDB,
	}
}

type StatusService struct {
	logger       *zap.Logger
	persistentDB *rw.Conn
}

func (s *StatusService) RegisterGRPC(server *grpc.Server) {
	api.RegisterStatusServiceServer(server, s)
}

func (s *StatusService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterStatusServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *StatusService) GetQuotaSummary(ctx context.Context, empty *emptypb.Empty) (*api.Quotas, error) {

	res := &api.Quotas{}
	err := s.persistentDB.WithStore(func(s *rw.Store) error {
		_, bytesDownloaded, err := s.GetMediaStatsForCurrentMonth(ctx)
		if err != nil {
			return err
		}
		res.BandwidthTotalMib = quota.BandwidthQuotaInMiB
		res.BandwidthRemainingMib = quota.BandwidthQuotaInMiB - quota.BytesAsMib(bytesDownloaded)
		return nil
	})
	if err != nil {
		return nil, ErrInternal(err).Err()
	}
	return res, nil
}

func (s *StatusService) Health(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	// if this endpoint is even available the API is somewhat functional.
	return &emptypb.Empty{}, nil
}
