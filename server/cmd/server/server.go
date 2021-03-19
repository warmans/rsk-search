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

	grpcCfg := server.GrpcServerConfig{}
	srvCfg := config.SearchServiceConfig{}
	roDbCfg := &common.Config{}
	rwDbCfg := &common.Config{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the grpc/http server",
		RunE: func(cmd *cobra.Command, args []string) error {

			flag.Parse()

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			rskIndex, err := bleve.Open(srvCfg.BleveIndexPath)
			if err != nil {
				return err
			}

			// DB is volatile and will be recreated with each deployment
			logger.Info("Init read-only DB...", zap.String("path", rwDbCfg.DSN))
			readOnlyStoreConn, err := ro.NewConn(roDbCfg)
			if err != nil {
				return err
			}
			defer func() {
				if err := readOnlyStoreConn.Close(); err != nil {
					logger.Error("failed to close RO db", zap.Error(err))
				}
			}()

			// DB is persistent and will retain data between deployments
			logger.Info("Init persistent DB...")
			persistentDBConn, err := rw.NewConn(rwDbCfg)
			if err != nil {
				return err
			}
			defer func() {
				if err := persistentDBConn.Close(); err != nil {
					logger.Error("failed to close persistent db", zap.Error(err))
				}
			}()
			logger.Info("Running persistent DB migrations")
			if err := persistentDBConn.Migrate(); err != nil {
				return err
			}

			grpcServices := []server.GRPCService{
				grpc.NewSearchService(search.NewSearch(rskIndex, readOnlyStoreConn), readOnlyStoreConn, persistentDBConn),
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

	grpcCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	srvCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	roDbCfg.RegisterFlags(cmd.Flags(), ServicePrefix, "ro")
	rwDbCfg.RegisterFlags(cmd.Flags(), ServicePrefix, "rw")

	return cmd
}

