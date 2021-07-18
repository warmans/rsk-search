package server

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/reward"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/server"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"github.com/warmans/rsk-search/service/grpc"
	"github.com/warmans/rsk-search/service/http"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ServicePrefix = "RSK_SEARCH"

func ServerCmd() *cobra.Command {

	grpcCfg := server.GrpcServerConfig{}
	srvCfg := config.SearchServiceConfig{}
	roDbCfg := &common.Config{}
	rwDbCfg := &common.Config{}
	oauthCfg := &oauth.Config{}
	jwtConfig := &jwt.Config{}
	rewardCfg := reward.Config{}
	pledgeCfg := pledge.Config{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the grpc/http server",
		RunE: func(cmd *cobra.Command, args []string) error {

			flag.Parse()

			var logger *zap.Logger
			var loggerErr error
			if os.Getenv("DEBUG") == "false" {
				conf := zap.NewDevelopmentConfig()
				conf.DisableStacktrace = true
				logger, loggerErr = conf.Build()
			} else {
				conf := zap.NewProductionConfig()
				logger, loggerErr = conf.Build()
			}
			if loggerErr != nil {
				panic(loggerErr)
			}
			defer logger.Sync() // flushes buffer, if any

			logger.Info("Init index...")
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

			/// setup rewards worker
			worker := reward.NewWorker(persistentDBConn, logger, rewardCfg)
			go func() {
				if err := worker.Start(); err != nil {
					logger.Fatal("worker failed", zap.Error(err))
				}
			}()
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				if err := worker.Stop(ctx); err != nil {
					logger.Error("worker stop failed", zap.Error(err))
				}
			}()

			// setup oauth
			tokenCache := oauth.NewCSRFCache()
			auth := jwt.NewAuth(jwtConfig)

			// validate pledge config
			if pledgeCfg.APIKey == "" {
				return fmt.Errorf("pledge API key was missing")
			}

			grpcServices := []server.GRPCService{
				grpc.NewSearchService(
					logger,
					srvCfg,
					search.NewSearch(rskIndex, readOnlyStoreConn),
					readOnlyStoreConn,
					auth,
				),
				grpc.NewTscriptService(
					logger,
					srvCfg,
					persistentDBConn,
					auth,
					pledge.NewClient(pledgeCfg),
				),
				grpc.NewOauthService(
					logger,
					tokenCache,
					oauthCfg,
				),
			}

			httpServices := []server.HTTPService{http.NewDownloadService(logger, srvCfg)}
			if oauthCfg.Secret != "" {
				httpServices = append(httpServices, http.NewOauthService(logger, tokenCache, persistentDBConn, auth, oauthCfg, srvCfg))
			} else {
				logger.Info("OAUTH SECRET WAS MISSING - OAUTH ENDPOINTS WILL NOT BE REGISTERED!")
			}

			srv, err := server.NewServer(logger, grpcCfg, grpcServices, httpServices)
			if err != nil {
				logger.Fatal("failed to create server", zap.Error(err))
			}

			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				srv.Stop()
			}()
			go func() {
				if err := srv.StartHTTP(); err != nil {
					logger.Fatal("HTTP Failed", zap.Error(err))
				}
			}()
			if err := srv.StartGRPC(); err != nil {
				logger.Fatal("GRPC Failed", zap.Error(err))
			}
			return nil
		},
	}

	grpcCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	srvCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	roDbCfg.RegisterFlags(cmd.Flags(), ServicePrefix, "ro")
	rwDbCfg.RegisterFlags(cmd.Flags(), ServicePrefix, "rw")
	oauthCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	jwtConfig.RegisterFlags(cmd.Flags(), ServicePrefix)
	rewardCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	pledgeCfg.RegisterFlags(cmd.Flags(), ServicePrefix)

	return cmd
}
