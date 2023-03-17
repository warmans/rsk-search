package sentry

import (
	"github.com/TheZeroSlave/zapsentry"
	sentryGo "github.com/getsentry/sentry-go"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sentryInitialized = false

type Config struct {
	DSN string
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.DSN, prefix, "sentry-dsn", "", "Sentry DSN")
}

func NewClient(cfg *Config) (*sentryGo.Client, error) {
	sentryClient, err := sentryGo.NewClient(sentryGo.ClientOptions{
		Dsn: cfg.DSN,
	})
	if err != nil {
		return nil, err
	}
	return sentryClient, nil
}

func ZapHook(env string) func(zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		if entry.Level < zap.ErrorLevel {
			return nil
		}
		if sentryInitialized {
			sentryGo.WithScope(func(scope *sentryGo.Scope) {
				scope.SetLevel(zapLevelToSentry(entry.Level))

				scope.SetTag("level", entry.Level.String())
				scope.SetTag("environment", env)

				scope.SetExtra("stacktrace", entry.Stack)
				scope.SetExtra("caller", entry.Caller.String())
				scope.SetExtra("logger_name", entry.LoggerName)

				sentryGo.CaptureMessage(entry.Message)
			})
		}
		return nil
	}
}

func zapLevelToSentry(level zapcore.Level) sentryGo.Level {
	switch level {
	case zapcore.DebugLevel:
		return sentryGo.LevelDebug
	case zapcore.InfoLevel:
		return sentryGo.LevelInfo
	case zapcore.WarnLevel:
		return sentryGo.LevelWarning
	case zapcore.ErrorLevel:
		return sentryGo.LevelError
	case zapcore.FatalLevel, zapcore.DPanicLevel, zapcore.PanicLevel:
		return sentryGo.LevelFatal
	default:
		return sentryGo.LevelInfo
	}
}

func LoggerWithSentry(log *zap.Logger, client *sentryGo.Client, env string) *zap.Logger {

	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"component":   "server",
			"environment": env,
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))

	if err != nil {
		//in case of err it will return noop core. so we can safely attach it
		log.Warn("failed to init zap", zap.Error(err))
	}

	log = zapsentry.AttachCoreToLogger(core, log)
	return log.With(zapsentry.NewScope())
}
