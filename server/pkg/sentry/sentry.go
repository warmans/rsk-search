package sentry

import (
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

func InitSentry(cfg *Config) error {
	if cfg.DSN == "" {
		return nil
	}
	err := sentryGo.Init(sentryGo.ClientOptions{
		Dsn: cfg.DSN,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	sentryInitialized = true
	return err
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
