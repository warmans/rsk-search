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
const _ = grpc.SupportPackageIsVersion7

// SearchServiceClient is the client API for SearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchServiceClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResultList, error)
	GetMetadata(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Metadata, error)
	ListFieldValues(ctx context.Context, in *ListFieldValuesRequest, opts ...grpc.CallOption) (*FieldValueList, error)
	PredictSearchTerm(ctx context.Context, in *PredictSearchTermRequest, opts ...grpc.CallOption) (*SearchTermPredictions, error)
	GetTranscript(ctx context.Context, in *GetTranscriptRequest, opts ...grpc.CallOption) (*Transcript, error)
	ListTranscripts(ctx context.Context, in *ListTranscriptsRequest, opts ...grpc.CallOption) (*TranscriptList, error)
	// changelogs
	ListChangelogs(ctx context.Context, in *ListChangelogsRequest, opts ...grpc.CallOption) (*ChangelogList, error)
}

type searchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchServiceClient(cc grpc.ClientConnInterface) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResultList, error) {
	out := new(SearchResultList)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetMetadata(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Metadata, error) {
	out := new(Metadata)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/GetMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListFieldValues(ctx context.Context, in *ListFieldValuesRequest, opts ...grpc.CallOption) (*FieldValueList, error) {
	out := new(FieldValueList)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/ListFieldValues", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) PredictSearchTerm(ctx context.Context, in *PredictSearchTermRequest, opts ...grpc.CallOption) (*SearchTermPredictions, error) {
	out := new(SearchTermPredictions)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/PredictSearchTerm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetTranscript(ctx context.Context, in *GetTranscriptRequest, opts ...grpc.CallOption) (*Transcript, error) {
	out := new(Transcript)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/GetTranscript", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListTranscripts(ctx context.Context, in *ListTranscriptsRequest, opts ...grpc.CallOption) (*TranscriptList, error) {
	out := new(TranscriptList)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/ListTranscripts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListChangelogs(ctx context.Context, in *ListChangelogsRequest, opts ...grpc.CallOption) (*ChangelogList, error) {
	out := new(ChangelogList)
	err := c.cc.Invoke(ctx, "/rsk.SearchService/ListChangelogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServiceServer is the server API for SearchService service.
// All implementations should embed UnimplementedSearchServiceServer
// for forward compatibility
type SearchServiceServer interface {
	Search(context.Context, *SearchRequest) (*SearchResultList, error)
	GetMetadata(context.Context, *emptypb.Empty) (*Metadata, error)
	ListFieldValues(context.Context, *ListFieldValuesRequest) (*FieldValueList, error)
	PredictSearchTerm(context.Context, *PredictSearchTermRequest) (*SearchTermPredictions, error)
	GetTranscript(context.Context, *GetTranscriptRequest) (*Transcript, error)
	ListTranscripts(context.Context, *ListTranscriptsRequest) (*TranscriptList, error)
	// changelogs
	ListChangelogs(context.Context, *ListChangelogsRequest) (*ChangelogList, error)
}

// UnimplementedSearchServiceServer should be embedded to have forward compatible implementations.
type UnimplementedSearchServiceServer struct {
}

func (UnimplementedSearchServiceServer) Search(context.Context, *SearchRequest) (*SearchResultList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedSearchServiceServer) GetMetadata(context.Context, *emptypb.Empty) (*Metadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetadata not implemented")
}
func (UnimplementedSearchServiceServer) ListFieldValues(context.Context, *ListFieldValuesRequest) (*FieldValueList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFieldValues not implemented")
}
func (UnimplementedSearchServiceServer) PredictSearchTerm(context.Context, *PredictSearchTermRequest) (*SearchTermPredictions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PredictSearchTerm not implemented")
}
func (UnimplementedSearchServiceServer) GetTranscript(context.Context, *GetTranscriptRequest) (*Transcript, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTranscript not implemented")
}
func (UnimplementedSearchServiceServer) ListTranscripts(context.Context, *ListTranscriptsRequest) (*TranscriptList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTranscripts not implemented")
}
func (UnimplementedSearchServiceServer) ListChangelogs(context.Context, *ListChangelogsRequest) (*ChangelogList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListChangelogs not implemented")
}

// UnsafeSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServiceServer will
// result in compilation errors.
type UnsafeSearchServiceServer interface {
	mustEmbedUnimplementedSearchServiceServer()
}

func RegisterSearchServiceServer(s *grpc.Server, srv SearchServiceServer) {
	s.RegisterService(&_SearchService_serviceDesc, srv)
}

func _SearchService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/GetMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetMetadata(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ListFieldValues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFieldValuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ListFieldValues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/ListFieldValues",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListFieldValues(ctx, req.(*ListFieldValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_PredictSearchTerm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictSearchTermRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).PredictSearchTerm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/PredictSearchTerm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).PredictSearchTerm(ctx, req.(*PredictSearchTermRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetTranscript_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTranscriptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetTranscript(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/GetTranscript",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetTranscript(ctx, req.(*GetTranscriptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ListTranscripts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTranscriptsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ListTranscripts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/ListTranscripts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListTranscripts(ctx, req.(*ListTranscriptsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ListChangelogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChangelogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ListChangelogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsk.SearchService/ListChangelogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListChangelogs(ctx, req.(*ListChangelogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SearchService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rsk.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _SearchService_Search_Handler,
		},
		{
			MethodName: "GetMetadata",
			Handler:    _SearchService_GetMetadata_Handler,
		},
		{
			MethodName: "ListFieldValues",
			Handler:    _SearchService_ListFieldValues_Handler,
		},
		{
			MethodName: "PredictSearchTerm",
			Handler:    _SearchService_PredictSearchTerm_Handler,
		},
		{
			MethodName: "GetTranscript",
			Handler:    _SearchService_GetTranscript_Handler,
		},
		{
			MethodName: "ListTranscripts",
			Handler:    _SearchService_ListTranscripts_Handler,
		},
		{
			MethodName: "ListChangelogs",
			Handler:    _SearchService_ListChangelogs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search.proto",
}
