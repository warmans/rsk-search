package grpc

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func NewRadioService(
	logger *zap.Logger,
	persistentDB *rw.Conn,
	auth *jwt.Auth,
	episodeCache *data.EpisodeCache,
) *RadioService {
	return &RadioService{
		logger:       logger,
		persistentDB: persistentDB,
		auth:         auth,
		episodeCache: episodeCache,
	}
}

type RadioService struct {
	logger       *zap.Logger
	persistentDB *rw.Conn
	auth         *jwt.Auth
	episodeCache *data.EpisodeCache
}

func (s *RadioService) RegisterGRPC(server *grpc.Server) {
	api.RegisterRadioServiceServer(server, s)
}

func (s *RadioService) RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	if err := api.RegisterRadioServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		panic(err)
	}
}

func (s *RadioService) GetState(ctx context.Context, empty *emptypb.Empty) (*api.RadioState, error) {
	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	var state *models.RadioState
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		state, err = s.GetLatestRadioState(ctx, claims.AuthorID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return state.Proto(), nil
}

func (s *RadioService) PutState(ctx context.Context, request *api.PutStateRequest) (*emptypb.Empty, error) {
	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}

	startedAt, err := time.Parse(time.RFC3339, request.CurrentEpisode.StartedAt)
	if err != nil {
		return nil, ErrInvalidRequestField("current_episode.started_at", err)
	}

	state := &models.RadioState{
		AuthorID:         claims.AuthorID,
		EpisodeID:        request.CurrentEpisode.ShortId,
		StartedAt:        startedAt,
		CurrentTimestamp: time.Duration(request.CurrentTimestampMs) * time.Millisecond,
	}
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		return s.SetRadioState(ctx, state)
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *RadioService) GetNext(ctx context.Context, empty *emptypb.Empty) (*api.NextEpisode, error) {
	claims, err := GetClaims(ctx, s.auth)
	if err != nil {
		return nil, err
	}
	var nextEpisodeID string
	err = s.persistentDB.WithStore(func(s *rw.Store) error {
		var err error
		nextEpisodeID, err = s.GetRadioNext(ctx, claims.AuthorID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &api.NextEpisode{
		ShortId: nextEpisodeID,
	}, nil
}
