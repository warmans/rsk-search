package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/warmans/rsk-search/service/grpc"
	"go.uber.org/zap"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor() googleGrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *googleGrpc.UnaryServerInfo, handler googleGrpc.UnaryHandler) (_ interface{}, err error) {
		resp, err := handler(ctx, req)
		return resp, handleError(ctx, err)
	}
}

func StreamErrorServerInterceptor() googleGrpc.StreamServerInterceptor {
	return func(srv interface{}, stream googleGrpc.ServerStream, info *googleGrpc.StreamServerInfo, handler googleGrpc.StreamHandler) (err error) {
		return handleError(stream.Context(), handler(srv, stream))
	}
}

func handleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	staErr, ok := err.(*grpc.Status)
	if !ok || ok && staErr.Sta.Code() == codes.Internal {
		ctxzap.Extract(ctx).Error("Internal error was returned by the API", zap.String("reason", err.Error()))

		// obfuscate internal errors
		s := status.New(codes.Internal, "Bauhaus is not working! Internal server error returned.")
		return s.Err()
	}
	return err
}
