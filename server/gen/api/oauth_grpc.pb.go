// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: oauth.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	OauthService_GetAuthURL_FullMethodName = "/rsk.OauthService/GetAuthURL"
)

// OauthServiceClient is the client API for OauthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OauthServiceClient interface {
	GetAuthURL(ctx context.Context, in *GetAuthURLRequest, opts ...grpc.CallOption) (*AuthURL, error)
}

type oauthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOauthServiceClient(cc grpc.ClientConnInterface) OauthServiceClient {
	return &oauthServiceClient{cc}
}

func (c *oauthServiceClient) GetAuthURL(ctx context.Context, in *GetAuthURLRequest, opts ...grpc.CallOption) (*AuthURL, error) {
	out := new(AuthURL)
	err := c.cc.Invoke(ctx, OauthService_GetAuthURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OauthServiceServer is the server API for OauthService service.
// All implementations should embed UnimplementedOauthServiceServer
// for forward compatibility
type OauthServiceServer interface {
	GetAuthURL(context.Context, *GetAuthURLRequest) (*AuthURL, error)
}

// UnimplementedOauthServiceServer should be embedded to have forward compatible implementations.
type UnimplementedOauthServiceServer struct {
}

func (UnimplementedOauthServiceServer) GetAuthURL(context.Context, *GetAuthURLRequest) (*AuthURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthURL not implemented")
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

func _OauthService_GetAuthURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OauthServiceServer).GetAuthURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OauthService_GetAuthURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OauthServiceServer).GetAuthURL(ctx, req.(*GetAuthURLRequest))
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
			MethodName: "GetAuthURL",
			Handler:    _OauthService_GetAuthURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "oauth.proto",
}
