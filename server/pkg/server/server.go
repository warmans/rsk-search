package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/server/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

// GRPCService describes a gRPC GRPCService.
type GRPCService interface {
	RegisterGRPC(*grpc.Server)
	RegisterHTTP(ctx context.Context, router *mux.Router, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption)
}

type HTTPService interface {
	RegisterHTTP(ctx context.Context, router *mux.Router)
}

type GrpcServerConfig struct {
	GRPCAddr string
	HTTPAddr string
}

func (c *GrpcServerConfig) RegisterFlags(prefix string) {
	flag.StringVarEnv(&c.GRPCAddr, prefix, "grpc-addr", "0.0.0.0:9090", "GRPC bind address")
	flag.StringVarEnv(&c.HTTPAddr, prefix, "http-addr", ":8888", "HTTP bind address")
}

func NewServer(logger *zap.Logger, cfg GrpcServerConfig, grpcServices []GRPCService, httpServices []HTTPService) (*Server, error) {

	grpc := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithMessageProducer(middleware.LogMessageProducer())),
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithMessageProducer(middleware.LogMessageProducer())),
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	s := &Server{
		cfg:          cfg,
		logger:       logger,
		grpc:         grpc,
		grpcServices: grpcServices,
		httpServices: httpServices,
	}
	return s, nil
}

type Server struct {
	cfg          GrpcServerConfig
	logger       *zap.Logger
	grpc         *grpc.Server
	grpcServices []GRPCService
	httpServices []HTTPService
}

func (s *Server) StartGRPC() error {
	for _, srv := range s.grpcServices {
		srv.RegisterGRPC(s.grpc)
		s.logger.Info(fmt.Sprintf("Registered GRPC service %T", srv))
	}
	lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		return errors.Wrap(err, "failed to listen to address")
	}

	s.logger.Info("Starting gRPC server", zap.String("addr", lis.Addr().String()))
	defer s.logger.Info("gRPC server stopped")

	return s.grpc.Serve(lis)
}

func (s *Server) StartHTTP() error {
	defer s.logger.Info("HTTP server stopped")

	router := mux.NewRouter()
	gwmux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	for _, srv := range s.httpServices {
		srv.RegisterHTTP(ctx, router)
	}

	for _, srv := range s.grpcServices {
		srv.RegisterHTTP(ctx, router, gwmux, s.cfg.GRPCAddr, opts)
		s.logger.Info(fmt.Sprintf("Registered HTTP service %T", srv))
	}

	// this must be after the service registration for whatever reason.
	router.PathPrefix("/").Handler(gwmux)

	s.logger.Info("Starting HTTP server", zap.String("addr", s.cfg.HTTPAddr))
	return http.ListenAndServe(s.cfg.HTTPAddr, router)
}
