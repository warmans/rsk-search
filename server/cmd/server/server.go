package server

import (
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/server"
	"github.com/warmans/rsk-search/pkg/service/grpc"
	"go.uber.org/zap"
)

const ServicePrefix = "RSK_SEARCH"

func ServerCmd() *cobra.Command {

	flag.Parse()

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the grpc/http server",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			grpcCfg := server.GrpcServerConfig{}
			grpcCfg.RegisterFlags(ServicePrefix)

			grpcServices := []server.GRPCService{
				grpc.NewSearchService(),
			}

			srv, err := server.NewServer(logger, grpcCfg, grpcServices, []server.HTTPService{})
			if err != nil {
				logger.Fatal("failed to create server", zap.Error(err))
			}
			go func() {
				if err := srv.StartGRPC(); err != nil {
					logger.Fatal("GRPC Failed", zap.Error(err))
				}
			}()
			if err := srv.StartHTTP(); err != nil {
				logger.Fatal("HTTP Failed", zap.Error(err))
			}
			return nil
		},
	}

	return cmd
}
