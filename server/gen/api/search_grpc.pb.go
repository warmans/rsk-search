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

// SearchServiceClient is the client API for SearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchServiceClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResultList, error)
	GetSearchMetadata(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SearchMetadata, error)
	ListFieldValues(ctx context.Context, in *ListFieldValuesRequest, opts ...grpc.CallOption) (*FieldValueList, error)
	GetEpisode(ctx context.Context, in *GetEpisodeRequest, opts ...grpc.CallOption) (*Episode, error)
	ListEpisodes(ctx context.Context, in *ListEpisodesRequest, opts ...grpc.CallOption) (*EpisodeList, error)
	// tscript is an incomplete transcription
	// chunks are ~2 min sections of the transcription
	GetTscriptChunkStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ChunkStats, error)
	GetTscriptChunk(ctx context.Context, in *GetTscriptChunkRequest, opts ...grpc.CallOption) (*TscriptChunk, error)
	ListChunkContributions(ctx context.Context, in *ListChunkContributionsRequest, opts ...grpc.CallOption) (*ChunkContributionList, error)
	ListAuthorContributions(ctx context.Context, in *ListAuthorContributionsRequest, opts ...grpc.CallOption) (*ChunkContributionList, error)
	GetChunkContribution(ctx context.Context, in *GetChunkContributionRequest, opts ...grpc.CallOption) (*ChunkContribution, error)
	CreateChunkContribution(ctx context.Context, in *CreateChunkContributionRequest, opts ...grpc.CallOption) (*ChunkContribution, error)
	UpdateChunkContribution(ctx context.Context, in *UpdateChunkContributionRequest, opts ...grpc.CallOption) (*ChunkContribution, error)
	RequestChunkContributionState(ctx context.Context, in *RequestChunkContributionStateRequest, opts ...grpc.CallOption) (*ChunkContribution, error)
	SubmitDialogCorrection(ctx context.Context, in *SubmitDialogCorrectionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetRedditAuthURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RedditAuthURL, error)
}

type searchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchServiceClient(cc grpc.ClientConnInterface) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResultList, error) {
	out := new(SearchResultList)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetSearchMetadata(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SearchMetadata, error) {
	out := new(SearchMetadata)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/GetSearchMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListFieldValues(ctx context.Context, in *ListFieldValuesRequest, opts ...grpc.CallOption) (*FieldValueList, error) {
	out := new(FieldValueList)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/ListFieldValues", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetEpisode(ctx context.Context, in *GetEpisodeRequest, opts ...grpc.CallOption) (*Episode, error) {
	out := new(Episode)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/GetEpisode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListEpisodes(ctx context.Context, in *ListEpisodesRequest, opts ...grpc.CallOption) (*EpisodeList, error) {
	out := new(EpisodeList)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/ListEpisodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetTscriptChunkStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ChunkStats, error) {
	out := new(ChunkStats)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/GetTscriptChunkStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetTscriptChunk(ctx context.Context, in *GetTscriptChunkRequest, opts ...grpc.CallOption) (*TscriptChunk, error) {
	out := new(TscriptChunk)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/GetTscriptChunk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListChunkContributions(ctx context.Context, in *ListChunkContributionsRequest, opts ...grpc.CallOption) (*ChunkContributionList, error) {
	out := new(ChunkContributionList)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/ListChunkContributions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) ListAuthorContributions(ctx context.Context, in *ListAuthorContributionsRequest, opts ...grpc.CallOption) (*ChunkContributionList, error) {
	out := new(ChunkContributionList)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/ListAuthorContributions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetChunkContribution(ctx context.Context, in *GetChunkContributionRequest, opts ...grpc.CallOption) (*ChunkContribution, error) {
	out := new(ChunkContribution)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/GetChunkContribution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) CreateChunkContribution(ctx context.Context, in *CreateChunkContributionRequest, opts ...grpc.CallOption) (*ChunkContribution, error) {
	out := new(ChunkContribution)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/CreateChunkContribution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) UpdateChunkContribution(ctx context.Context, in *UpdateChunkContributionRequest, opts ...grpc.CallOption) (*ChunkContribution, error) {
	out := new(ChunkContribution)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/UpdateChunkContribution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) RequestChunkContributionState(ctx context.Context, in *RequestChunkContributionStateRequest, opts ...grpc.CallOption) (*ChunkContribution, error) {
	out := new(ChunkContribution)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/RequestChunkContributionState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) SubmitDialogCorrection(ctx context.Context, in *SubmitDialogCorrectionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/SubmitDialogCorrection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetRedditAuthURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RedditAuthURL, error) {
	out := new(RedditAuthURL)
	err := c.cc.Invoke(ctx, "/rsksearch.SearchService/GetRedditAuthURL", in, out, opts...)
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
	GetSearchMetadata(context.Context, *emptypb.Empty) (*SearchMetadata, error)
	ListFieldValues(context.Context, *ListFieldValuesRequest) (*FieldValueList, error)
	GetEpisode(context.Context, *GetEpisodeRequest) (*Episode, error)
	ListEpisodes(context.Context, *ListEpisodesRequest) (*EpisodeList, error)
	// tscript is an incomplete transcription
	// chunks are ~2 min sections of the transcription
	GetTscriptChunkStats(context.Context, *emptypb.Empty) (*ChunkStats, error)
	GetTscriptChunk(context.Context, *GetTscriptChunkRequest) (*TscriptChunk, error)
	ListChunkContributions(context.Context, *ListChunkContributionsRequest) (*ChunkContributionList, error)
	ListAuthorContributions(context.Context, *ListAuthorContributionsRequest) (*ChunkContributionList, error)
	GetChunkContribution(context.Context, *GetChunkContributionRequest) (*ChunkContribution, error)
	CreateChunkContribution(context.Context, *CreateChunkContributionRequest) (*ChunkContribution, error)
	UpdateChunkContribution(context.Context, *UpdateChunkContributionRequest) (*ChunkContribution, error)
	RequestChunkContributionState(context.Context, *RequestChunkContributionStateRequest) (*ChunkContribution, error)
	SubmitDialogCorrection(context.Context, *SubmitDialogCorrectionRequest) (*emptypb.Empty, error)
	GetRedditAuthURL(context.Context, *emptypb.Empty) (*RedditAuthURL, error)
}

// UnimplementedSearchServiceServer should be embedded to have forward compatible implementations.
type UnimplementedSearchServiceServer struct {
}

func (UnimplementedSearchServiceServer) Search(context.Context, *SearchRequest) (*SearchResultList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedSearchServiceServer) GetSearchMetadata(context.Context, *emptypb.Empty) (*SearchMetadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSearchMetadata not implemented")
}
func (UnimplementedSearchServiceServer) ListFieldValues(context.Context, *ListFieldValuesRequest) (*FieldValueList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFieldValues not implemented")
}
func (UnimplementedSearchServiceServer) GetEpisode(context.Context, *GetEpisodeRequest) (*Episode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEpisode not implemented")
}
func (UnimplementedSearchServiceServer) ListEpisodes(context.Context, *ListEpisodesRequest) (*EpisodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEpisodes not implemented")
}
func (UnimplementedSearchServiceServer) GetTscriptChunkStats(context.Context, *emptypb.Empty) (*ChunkStats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTscriptChunkStats not implemented")
}
func (UnimplementedSearchServiceServer) GetTscriptChunk(context.Context, *GetTscriptChunkRequest) (*TscriptChunk, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTscriptChunk not implemented")
}
func (UnimplementedSearchServiceServer) ListChunkContributions(context.Context, *ListChunkContributionsRequest) (*ChunkContributionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListChunkContributions not implemented")
}
func (UnimplementedSearchServiceServer) ListAuthorContributions(context.Context, *ListAuthorContributionsRequest) (*ChunkContributionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAuthorContributions not implemented")
}
func (UnimplementedSearchServiceServer) GetChunkContribution(context.Context, *GetChunkContributionRequest) (*ChunkContribution, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChunkContribution not implemented")
}
func (UnimplementedSearchServiceServer) CreateChunkContribution(context.Context, *CreateChunkContributionRequest) (*ChunkContribution, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChunkContribution not implemented")
}
func (UnimplementedSearchServiceServer) UpdateChunkContribution(context.Context, *UpdateChunkContributionRequest) (*ChunkContribution, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChunkContribution not implemented")
}
func (UnimplementedSearchServiceServer) RequestChunkContributionState(context.Context, *RequestChunkContributionStateRequest) (*ChunkContribution, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestChunkContributionState not implemented")
}
func (UnimplementedSearchServiceServer) SubmitDialogCorrection(context.Context, *SubmitDialogCorrectionRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitDialogCorrection not implemented")
}
func (UnimplementedSearchServiceServer) GetRedditAuthURL(context.Context, *emptypb.Empty) (*RedditAuthURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRedditAuthURL not implemented")
}

// UnsafeSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServiceServer will
// result in compilation errors.
type UnsafeSearchServiceServer interface {
	mustEmbedUnimplementedSearchServiceServer()
}

func RegisterSearchServiceServer(s grpc.ServiceRegistrar, srv SearchServiceServer) {
	s.RegisterService(&SearchService_ServiceDesc, srv)
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
		FullMethod: "/rsksearch.SearchService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetSearchMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetSearchMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/GetSearchMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetSearchMetadata(ctx, req.(*emptypb.Empty))
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
		FullMethod: "/rsksearch.SearchService/ListFieldValues",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListFieldValues(ctx, req.(*ListFieldValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetEpisode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEpisodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetEpisode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/GetEpisode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetEpisode(ctx, req.(*GetEpisodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ListEpisodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEpisodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ListEpisodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/ListEpisodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListEpisodes(ctx, req.(*ListEpisodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetTscriptChunkStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetTscriptChunkStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/GetTscriptChunkStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetTscriptChunkStats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetTscriptChunk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTscriptChunkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetTscriptChunk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/GetTscriptChunk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetTscriptChunk(ctx, req.(*GetTscriptChunkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ListChunkContributions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChunkContributionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ListChunkContributions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/ListChunkContributions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListChunkContributions(ctx, req.(*ListChunkContributionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_ListAuthorContributions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAuthorContributionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).ListAuthorContributions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/ListAuthorContributions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).ListAuthorContributions(ctx, req.(*ListAuthorContributionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetChunkContribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChunkContributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetChunkContribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/GetChunkContribution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetChunkContribution(ctx, req.(*GetChunkContributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_CreateChunkContribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChunkContributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).CreateChunkContribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/CreateChunkContribution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).CreateChunkContribution(ctx, req.(*CreateChunkContributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_UpdateChunkContribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateChunkContributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).UpdateChunkContribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/UpdateChunkContribution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).UpdateChunkContribution(ctx, req.(*UpdateChunkContributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_RequestChunkContributionState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestChunkContributionStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).RequestChunkContributionState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/RequestChunkContributionState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).RequestChunkContributionState(ctx, req.(*RequestChunkContributionStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_SubmitDialogCorrection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitDialogCorrectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).SubmitDialogCorrection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/SubmitDialogCorrection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).SubmitDialogCorrection(ctx, req.(*SubmitDialogCorrectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetRedditAuthURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetRedditAuthURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsksearch.SearchService/GetRedditAuthURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetRedditAuthURL(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchService_ServiceDesc is the grpc.ServiceDesc for SearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rsksearch.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _SearchService_Search_Handler,
		},
		{
			MethodName: "GetSearchMetadata",
			Handler:    _SearchService_GetSearchMetadata_Handler,
		},
		{
			MethodName: "ListFieldValues",
			Handler:    _SearchService_ListFieldValues_Handler,
		},
		{
			MethodName: "GetEpisode",
			Handler:    _SearchService_GetEpisode_Handler,
		},
		{
			MethodName: "ListEpisodes",
			Handler:    _SearchService_ListEpisodes_Handler,
		},
		{
			MethodName: "GetTscriptChunkStats",
			Handler:    _SearchService_GetTscriptChunkStats_Handler,
		},
		{
			MethodName: "GetTscriptChunk",
			Handler:    _SearchService_GetTscriptChunk_Handler,
		},
		{
			MethodName: "ListChunkContributions",
			Handler:    _SearchService_ListChunkContributions_Handler,
		},
		{
			MethodName: "ListAuthorContributions",
			Handler:    _SearchService_ListAuthorContributions_Handler,
		},
		{
			MethodName: "GetChunkContribution",
			Handler:    _SearchService_GetChunkContribution_Handler,
		},
		{
			MethodName: "CreateChunkContribution",
			Handler:    _SearchService_CreateChunkContribution_Handler,
		},
		{
			MethodName: "UpdateChunkContribution",
			Handler:    _SearchService_UpdateChunkContribution_Handler,
		},
		{
			MethodName: "RequestChunkContributionState",
			Handler:    _SearchService_RequestChunkContributionState_Handler,
		},
		{
			MethodName: "SubmitDialogCorrection",
			Handler:    _SearchService_SubmitDialogCorrection_Handler,
		},
		{
			MethodName: "GetRedditAuthURL",
			Handler:    _SearchService_GetRedditAuthURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search.proto",
}
