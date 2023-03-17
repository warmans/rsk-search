package middleware

import (
	"context"
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CodeToLevel() grpc_zap.CodeToLevel {
	return func(code codes.Code) zapcore.Level {
		if code == codes.OK {
			return zapcore.DebugLevel
		}
		// actual errors
		if code == codes.Internal || code == codes.Unavailable || code == codes.ResourceExhausted {
			return zapcore.ErrorLevel
		}
		// user errors
		if code == codes.FailedPrecondition || code == codes.InvalidArgument || code == codes.NotFound {
			return zapcore.InfoLevel
		}
		// anything else
		return zapcore.WarnLevel
	}
}

func LogMessageProducer() grpc_zap.MessageProducer {
	return func(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {

		// don't spam logs with OK responses
		if level < zapcore.WarnLevel {
			return
		}

		fields := []zap.Field{
			duration,
		}
		if err != nil {
			fields = append(fields, zap.String("error", err.Error()))
			fields = append(fields, zap.String("error.cause", errors.Cause(err).Error()))
		}

		if sta, ok := status.FromError(err); ok {
			fields = append(fields, zap.String("grpc.code", sta.Code().String()))
			for _, detail := range sta.Details() {
				switch t := detail.(type) {
				case *errdetails.DebugInfo:
					fields = append(fields, zap.String("err.debug.detail", t.Detail))
					if len(t.StackEntries) > 0 {
						fields = append(fields, zap.Strings("err.debug.stack", t.StackEntries))
					}
				case *errdetails.ErrorInfo:
					fields = append(fields, zap.String("err.error.reason", t.Reason))
					if t.Domain != "" {
						fields = append(fields, zap.String("err.error.domain", t.Domain))
					}
					if t.Metadata != nil {
						for k, v := range t.Metadata {
							fields = append(fields, zap.String(fmt.Sprintf("err.meta.%s", k), v))
						}
					}
				case *errdetails.BadRequest:
					for k, violation := range t.GetFieldViolations() {
						fields = append(fields, zap.String(fmt.Sprintf("err.request.field[%d]", k), violation.GetField()))
						fields = append(fields, zap.String(fmt.Sprintf("err.request.violation[%d]", k), violation.GetDescription()))
					}
				case *errdetails.BadRequest_FieldViolation:
					fields = append(fields, zap.String("err.request.field", t.GetField()))
					fields = append(fields, zap.String("err.request.violation", t.GetDescription()))
				}
			}
		}
		ctxzap.Extract(ctx).Check(level, msg).Write(
			fields...,
		)
	}
}
