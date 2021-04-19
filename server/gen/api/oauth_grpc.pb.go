// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OauthServiceClient is the client API for OauthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OauthServiceClient interface {
	GetRedditAuthURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RedditAuthURL, error)
}

type oauthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOauthServiceClient(cc grpc.ClientConnInterface) OauthServiceClient {
	return &oauthServiceClient{cc}
}

func (c *oauthServiceClient) GetRedditAuthURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RedditAuthURL, error) {
	out := new(RedditAuthURL)
	err := c.cc.Invoke(ctx, "/rsk.OauthService/GetRedditAuthURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OauthServiceServer is the server API for OauthService service.
// All implementations should embed UnimplementedOauthServiceServer
// for forward compatibility
type OauthServiceServer interface {
	GetRedditAuthURL(context.Context, *emptypb.Empty) (*RedditAuthURL, error)
}

// UnimplementedOauthServiceServer should be embedded to have forward compatible implementations.
type UnimplementedOauthServiceServer struct {
}

func (UnimplementedOauthServiceServer) GetRedditAuthURL(context.Context, *emptypb.Empty) (*RedditAuthURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRedditAuthURL not implemented")
}

// UnsafeOauthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OauthServiceServer will
// result in compilation errors.
type UnsafeOauthServiceServer interface {
	mustEmbedUnimplementedOauthServiceServer()
}

func RegisterOauthServiceServer(s grpc.ServiceRegistrar, srv OauthServiceServer) {
	s.RegisterService(&OauthService_ServiceDesc, srv)
}

func _OauthService_GetRedditAuthURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OauthServiceServer).GetRedditAuthURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.OauthService/GetRedditAuthURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OauthServiceServer).GetRedditAuthURL(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// OauthService_ServiceDesc is the grpc.ServiceDesc for OauthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OauthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rsk.OauthService",
	HandlerType: (*OauthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRedditAuthURL",
			Handler:    _OauthService_GetRedditAuthURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "oauth.proto",
}