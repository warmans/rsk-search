package server

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/server"
	"github.com/warmans/rsk-search/pkg/service/config"
	"github.com/warmans/rsk-search/pkg/service/grpc"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
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

			srvCfg := config.SearchServiceConfig{}
			srvCfg.RegisterFlags(ServicePrefix)

			roDbCfg := &common.Config{}
			roDbCfg.RegisterFlags(ServicePrefix, "ro")

			rwDbCfg := &common.Config{}
			rwDbCfg.RegisterFlags(ServicePrefix, "rw")

			flag.Parse()

			rskIndex, err := bleve.Open(srvCfg.BleveIndexPath)
			if err != nil {
				return err
			}

			// DB is volatile and will be recreated with each deployment
			readOnlyStoreConn, err := ro.NewConn(roDbCfg)
			if err != nil {
				return err
			}

			// DB is persistent and will retain data between deployments
			logger.Info("Init persistent DB", zap.String("path", rwDbCfg.DSN))
			persistentDBConn, err := rw.NewConn(rwDbCfg)
			if err != nil {
				return err
			}
			logger.Info("Running persistent DB migrations")
			if err := persistentDBConn.Migrate(); err != nil {
				return err
			}

			grpcServices := []server.GRPCService{
				grpc.NewSearchService(search.NewSearch(rskIndex, readOnlyStoreConn, persistentDBConn), readOnlyStoreConn),
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
