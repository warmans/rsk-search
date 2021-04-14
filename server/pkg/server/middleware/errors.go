package middleware

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
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
		s, err := status.New(codes.Internal, http.StatusText(http.StatusInternalServerError)).WithDetails(
			&errdetails.DebugInfo{
				Detail: "Unknown internal server error",
			},
		)
		if err != nil {
			return status.New(codes.Internal, "failed to create error").Err()
		}
		return s.Err()
	}
	return err
}
