// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: radio.proto

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

const (
	RadioService_GetRadioState_FullMethodName = "/rsk.RadioService/GetRadioState"
	RadioService_GetRadioNext_FullMethodName  = "/rsk.RadioService/GetRadioNext"
	RadioService_PutRadioState_FullMethodName = "/rsk.RadioService/PutRadioState"
)

// RadioServiceClient is the client API for RadioService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RadioServiceClient interface {
	GetRadioState(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RadioState, error)
	GetRadioNext(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NextRadioEpisode, error)
	PutRadioState(ctx context.Context, in *PutRadioStateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type radioServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRadioServiceClient(cc grpc.ClientConnInterface) RadioServiceClient {
	return &radioServiceClient{cc}
}

func (c *radioServiceClient) GetRadioState(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RadioState, error) {
	out := new(RadioState)
	err := c.cc.Invoke(ctx, RadioService_GetRadioState_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *radioServiceClient) GetRadioNext(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NextRadioEpisode, error) {
	out := new(NextRadioEpisode)
	err := c.cc.Invoke(ctx, RadioService_GetRadioNext_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *radioServiceClient) PutRadioState(ctx context.Context, in *PutRadioStateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RadioService_PutRadioState_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RadioServiceServer is the server API for RadioService service.
// All implementations should embed UnimplementedRadioServiceServer
// for forward compatibility
type RadioServiceServer interface {
	GetRadioState(context.Context, *emptypb.Empty) (*RadioState, error)
	GetRadioNext(context.Context, *emptypb.Empty) (*NextRadioEpisode, error)
	PutRadioState(context.Context, *PutRadioStateRequest) (*emptypb.Empty, error)
}

// UnimplementedRadioServiceServer should be embedded to have forward compatible implementations.
type UnimplementedRadioServiceServer struct {
}

func (UnimplementedRadioServiceServer) GetRadioState(context.Context, *emptypb.Empty) (*RadioState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRadioState not implemented")
}
func (UnimplementedRadioServiceServer) GetRadioNext(context.Context, *emptypb.Empty) (*NextRadioEpisode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRadioNext not implemented")
}
func (UnimplementedRadioServiceServer) PutRadioState(context.Context, *PutRadioStateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutRadioState not implemented")
}

// UnsafeRadioServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RadioServiceServer will
// result in compilation errors.
type UnsafeRadioServiceServer interface {
	mustEmbedUnimplementedRadioServiceServer()
}

func RegisterRadioServiceServer(s grpc.ServiceRegistrar, srv RadioServiceServer) {
	s.RegisterService(&RadioService_ServiceDesc, srv)
}

func _RadioService_GetRadioState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RadioServiceServer).GetRadioState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RadioService_GetRadioState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RadioServiceServer).GetRadioState(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RadioService_GetRadioNext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RadioServiceServer).GetRadioNext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RadioService_GetRadioNext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RadioServiceServer).GetRadioNext(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RadioService_PutRadioState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRadioStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RadioServiceServer).PutRadioState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RadioService_PutRadioState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RadioServiceServer).PutRadioState(ctx, req.(*PutRadioStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RadioService_ServiceDesc is the grpc.ServiceDesc for RadioService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RadioService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rsk.RadioService",
	HandlerType: (*RadioServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRadioState",
			Handler:    _RadioService_GetRadioState_Handler,
		},
		{
			MethodName: "GetRadioNext",
			Handler:    _RadioService_GetRadioNext_Handler,
		},
		{
			MethodName: "PutRadioState",
			Handler:    _RadioService_PutRadioState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "radio.proto",
}