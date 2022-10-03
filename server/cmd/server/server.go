package server

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/blugelabs/bluge"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/assemblyai"
	"github.com/warmans/rsk-search/pkg/coffee"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/pledge"
	"github.com/warmans/rsk-search/pkg/reward"
	v2 "github.com/warmans/rsk-search/pkg/search/v2"
	"github.com/warmans/rsk-search/pkg/server"
	speech2text2 "github.com/warmans/rsk-search/pkg/speech2text"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"github.com/warmans/rsk-search/service/grpc"
	httpsrv "github.com/warmans/rsk-search/service/http"
	"github.com/warmans/rsk-search/service/queue"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"path"
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
	importQueueConfig := &queue.ImportQueueConfig{}
	speech2TextCfg := &speech2text2.GcloudConfig{}
	coffeeCfg := &coffee.Config{}
	assemblyAiCfg := &assemblyai.Config{}

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
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			episodeCache, err := data.NewEpisodeStore(path.Join(srvCfg.FilesBasePath, "data", "episodes"))
			if err != nil {
				logger.Fatal("failed to create episode cache", zap.Error(err))
			}

			// DB is volatile and will be recreated with each deployment
			logger.Info("Init read-only DB...", zap.String("path", roDbCfg.DSN))
			readOnlyStoreConn, err := ro.NewConn(roDbCfg)
			if err != nil {
				return err
			}
			defer func() {
				if err := readOnlyStoreConn.Close(); err != nil {
					logger.Error("failed to close RO db", zap.Error(err))
				}
			}()

			// search index
			blugeCfg := bluge.DefaultConfig(srvCfg.BlugeIndexPath)
			rskIndex, err := bluge.OpenReader(blugeCfg)
			if err != nil {
				return err
			}
			search := v2.NewSearch(rskIndex, readOnlyStoreConn, episodeCache, srvCfg.AudioUriPattern)

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

			//todo: need to mount credentials in prod env
			googleStorage, err := storage.NewClient(context.Background())
			if err != nil {
				logger.Fatal("Failed to create google storage client", zap.Error(err))
			}

			googleSpeech, err := speech.NewClient(context.Background())
			if err != nil {
				logger.Fatal("Failed to create google speech client", zap.Error(err))
			}

			// task queue
			taskQueue := queue.NewImportQueue(
				logger,
				afero.NewOsFs(),
				persistentDBConn,
				speech2text2.NewGcloud(logger, googleStorage, googleSpeech, speech2TextCfg),
				assemblyai.NewClient(logger, http.DefaultClient, assemblyAiCfg),
				googleStorage,
				importQueueConfig,
			)
			go func() {
				if err := taskQueue.Start(); err != nil {
					logger.Fatal("Task queue failed", zap.Error(err))
				}
			}()

			// buy-me-a-coffee client
			var coffeeClient *coffee.Client
			if coffeeCfg.AccessToken != "" {
				coffeeClient = coffee.NewClient(coffeeCfg)
				coffeeWorker := coffee.NewWorker(coffeeClient, persistentDBConn, logger, coffeeCfg)
				go func() {
					if err := coffeeWorker.Start(); err != nil {
						logger.Error("Coffee worker failed", zap.Error(err))
					}
					defer func() {
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						defer cancel()
						if err := coffeeWorker.Stop(ctx); err != nil {
							logger.Error("coffee worker stop failed", zap.Error(err))
						}
					}()
				}()
			} else {
				logger.Info("Coffee client disabled (no access token)")
			}

			grpcServices := []server.GRPCService{
				grpc.NewSearchService(
					logger,
					srvCfg,
					search,
					readOnlyStoreConn,
					persistentDBConn,
					auth,
					episodeCache,
				),
				grpc.NewContribService(
					logger,
					srvCfg,
					persistentDBConn,
					auth,
					pledge.NewClient(pledgeCfg),
					episodeCache,
					coffeeClient,
				),
				grpc.NewOauthService(
					logger,
					tokenCache,
					oauthCfg,
				),
				grpc.NewAdminService(
					logger,
					taskQueue,
					auth,
					persistentDBConn,
				),
			}

			httpServices := []server.HTTPService{
				httpsrv.NewMetricsService(),
				httpsrv.NewDownloadService(logger, srvCfg),
			}
			if oauthCfg.Secret != "" {
				httpServices = append(httpServices, httpsrv.NewOauthService(logger, tokenCache, persistentDBConn, auth, oauthCfg, srvCfg))
			} else {
				logger.Info("OAUTH SECRET WAS MISSING - OAUTH ENDPOINTS WILL NOT BE REGISTERED!")
			}

			srv, err := server.NewServer(logger, grpcCfg, grpcServices, httpServices)
			if err != nil {
				logger.Fatal("failed to create server", zap.Error(err))
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				srv.Stop()
				taskQueue.Stop()
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
	importQueueConfig.RegisterFlags(cmd.Flags(), ServicePrefix)
	speech2TextCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	coffeeCfg.RegisterFlags(cmd.Flags(), ServicePrefix)
	assemblyAiCfg.RegisterFlags(cmd.Flags(), ServicePrefix)

	return cmd
}
