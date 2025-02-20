package bot

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/archive"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"os"
	"os/signal"
)

func RootCommand() *cobra.Command {

	var botToken string
	var guildID string
	var apiTarget string
	var webUrl string
	var archiveDir string
	var tlsCertPath string

	jwtConfig := &jwt.Config{}

	cmd := &cobra.Command{
		Use:   "discord-bot",
		Short: "starts a discord-bot",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := createLogger()
			if err != nil {
				return err
			}
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			auth := jwt.NewAuth(jwtConfig)
			systemToken, err := auth.NewSystemJWT()
			if err != nil {
				return fmt.Errorf("failed to create system token: %w", err)
			}

			logger.Info("Creating discord session...")
			if botToken == "" {
				return fmt.Errorf("discord token is required")
			}
			//if guildID == "" {
			//	return fmt.Errorf("guildID is required")
			//}
			session, err := discordgo.New("Bot " + botToken)
			if err != nil {
				return fmt.Errorf("failed to create discord session: %w", err)
			}

			grpcConn, err := createGrpcClientConn(apiTarget, systemToken, tlsCertPath)
			if err != nil {
				return fmt.Errorf("failed to dial GRPC connection to API: %w", err)
			}
			bot := discord.NewBot(
				logger,
				session,
				guildID,
				webUrl,
				archive.NewStore(archiveDir),
				createTranscriptClient(grpcConn),
				createSearchClient(grpcConn),
			)
			logger.Info("Starting bot...")
			if err = bot.Start(); err != nil {
				return fmt.Errorf("failed to start bot: %w", err)
			}

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)
			<-stop

			log.Println("Gracefully shutting down")
			if err = bot.Close(); err != nil {
				return fmt.Errorf("failed to gracefully shutdown bot: %w", err)
			}
			return nil
		},
	}

	flag.StringVarEnv(cmd.Flags(), &botToken, "", "discord-token", "", "discord token")
	flag.StringVarEnv(cmd.Flags(), &guildID, "", "discord-guild-id", "", "discord guild/server ID")
	flag.StringVarEnv(cmd.Flags(), &apiTarget, "", "api-target", "127.0.0.1:9090", "gRPC API target")
	flag.StringVarEnv(cmd.Flags(), &webUrl, "", "web-url", "http://127.0.0.1:4200", "Base web address used for links")
	flag.StringVarEnv(cmd.Flags(), &archiveDir, "", "archive-dir", "./var/archive", "Location to archive files via archive command")
	flag.StringVarEnv(cmd.Flags(), &tlsCertPath, "", "grpc-tls-cert", "./x509/server_cert.pem", "TLS certificate needed to access GRPC server")
	jwtConfig.RegisterFlags(cmd.Flags(), "")

	flag.Parse()

	return cmd
}

func createGrpcClientConn(apiTarget string, systemToken string, tlsCertificatePath string) (*grpc.ClientConn, error) {

	creds, err := credentials.NewClientTLSFromFile(tlsCertificatePath, "")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	return grpc.DialContext(
		context.Background(),
		apiTarget,
		//todo:
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(
			oauth.TokenSource{
				TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: systemToken}),
			},
		),
	)
}

func createTranscriptClient(conn *grpc.ClientConn) api.TranscriptServiceClient {
	return api.NewTranscriptServiceClient(conn)
}

func createSearchClient(conn *grpc.ClientConn) api.SearchServiceClient {
	return api.NewSearchServiceClient(conn)
}

func createLogger() (*zap.Logger, error) {
	var loggerConf zap.Config
	if os.Getenv("DEBUG") == "false" {
		loggerConf = zap.NewDevelopmentConfig()
	} else {
		loggerConf = zap.NewProductionConfig()
	}
	logger, loggerErr := loggerConf.Build()
	if loggerErr != nil {
		return nil, fmt.Errorf("failed to create logger: %w", loggerErr)
	}
	return logger, nil
}
