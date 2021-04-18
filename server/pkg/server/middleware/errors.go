package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		resp, err := handler(ctx, req)
		return resp, handleError(err)
	}
}

func StreamErrorServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		return handleError(handler(srv, stream))
	}
}

func handleError(err error) error {
	if err == nil {
		return nil
	}
	staErr, ok := status.FromError(err)
	if !ok || ok &&  staErr.Code() == codes.Internal {
		s := status.New(codes.Internal, "Demicky response.")
		return s.Err()
	}
	return err
}
