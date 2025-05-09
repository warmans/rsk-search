// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: community.proto

package api

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommunityProject struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Summary       string                 `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Url           string                 `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommunityProject) Reset() {
	*x = CommunityProject{}
	mi := &file_community_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommunityProject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommunityProject) ProtoMessage() {}

func (x *CommunityProject) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommunityProject.ProtoReflect.Descriptor instead.
func (*CommunityProject) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{0}
}

func (x *CommunityProject) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CommunityProject) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CommunityProject) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *CommunityProject) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *CommunityProject) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CommunityProject) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type ListCommunityProjectsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filter        string                 `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	SortField     string                 `protobuf:"bytes,2,opt,name=sort_field,json=sortField,proto3" json:"sort_field,omitempty"`
	SortDirection string                 `protobuf:"bytes,3,opt,name=sort_direction,json=sortDirection,proto3" json:"sort_direction,omitempty"`
	Page          int32                  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,5,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCommunityProjectsRequest) Reset() {
	*x = ListCommunityProjectsRequest{}
	mi := &file_community_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCommunityProjectsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommunityProjectsRequest) ProtoMessage() {}

func (x *ListCommunityProjectsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommunityProjectsRequest.ProtoReflect.Descriptor instead.
func (*ListCommunityProjectsRequest) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{1}
}

func (x *ListCommunityProjectsRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *ListCommunityProjectsRequest) GetSortField() string {
	if x != nil {
		return x.SortField
	}
	return ""
}

func (x *ListCommunityProjectsRequest) GetSortDirection() string {
	if x != nil {
		return x.SortDirection
	}
	return ""
}

func (x *ListCommunityProjectsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCommunityProjectsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type CommunityProjectList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Projects      []*CommunityProject    `protobuf:"bytes,1,rep,name=projects,proto3" json:"projects,omitempty"`
	ResultCount   int32                  `protobuf:"varint,2,opt,name=result_count,json=resultCount,proto3" json:"result_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommunityProjectList) Reset() {
	*x = CommunityProjectList{}
	mi := &file_community_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommunityProjectList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommunityProjectList) ProtoMessage() {}

func (x *CommunityProjectList) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommunityProjectList.ProtoReflect.Descriptor instead.
func (*CommunityProjectList) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{2}
}

func (x *CommunityProjectList) GetProjects() []*CommunityProject {
	if x != nil {
		return x.Projects
	}
	return nil
}

func (x *CommunityProjectList) GetResultCount() int32 {
	if x != nil {
		return x.ResultCount
	}
	return 0
}

type ListArchiveRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	EpisodeIds    []string               `protobuf:"bytes,1,rep,name=episode_ids,json=episodeIds,proto3" json:"episode_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListArchiveRequest) Reset() {
	*x = ListArchiveRequest{}
	mi := &file_community_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListArchiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListArchiveRequest) ProtoMessage() {}

func (x *ListArchiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListArchiveRequest.ProtoReflect.Descriptor instead.
func (*ListArchiveRequest) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{3}
}

func (x *ListArchiveRequest) GetEpisodeIds() []string {
	if x != nil {
		return x.EpisodeIds
	}
	return nil
}

type ArchiveList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Archive             `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ArchiveList) Reset() {
	*x = ArchiveList{}
	mi := &file_community_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ArchiveList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArchiveList) ProtoMessage() {}

func (x *ArchiveList) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArchiveList.ProtoReflect.Descriptor instead.
func (*ArchiveList) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{4}
}

func (x *ArchiveList) GetItems() []*Archive {
	if x != nil {
		return x.Items
	}
	return nil
}

type Archive struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description    string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	RelatedEpisode string                 `protobuf:"bytes,3,opt,name=related_episode,json=relatedEpisode,proto3" json:"related_episode,omitempty"`
	// Deprecated: Marked as deprecated in community.proto.
	Files         []string `protobuf:"bytes,4,rep,name=files,proto3" json:"files,omitempty"`
	Media         []*File  `protobuf:"bytes,5,rep,name=media,proto3" json:"media,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Archive) Reset() {
	*x = Archive{}
	mi := &file_community_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Archive) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Archive) ProtoMessage() {}

func (x *Archive) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Archive.ProtoReflect.Descriptor instead.
func (*Archive) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{5}
}

func (x *Archive) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Archive) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Archive) GetRelatedEpisode() string {
	if x != nil {
		return x.RelatedEpisode
	}
	return ""
}

// Deprecated: Marked as deprecated in community.proto.
func (x *Archive) GetFiles() []string {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *Archive) GetMedia() []*File {
	if x != nil {
		return x.Media
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ThumbnailName string                 `protobuf:"bytes,2,opt,name=thumbnail_name,json=thumbnailName,proto3" json:"thumbnail_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *File) Reset() {
	*x = File{}
	mi := &file_community_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_community_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_community_proto_rawDescGZIP(), []int{6}
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetThumbnailName() string {
	if x != nil {
		return x.ThumbnailName
	}
	return ""
}

var File_community_proto protoreflect.FileDescriptor

const file_community_proto_rawDesc = "" +
	"\n" +
	"\x0fcommunity.proto\x12\x03rsk\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\x1a\x1bgoogle/protobuf/empty.proto\"\x9b\x01\n" +
	"\x10CommunityProject\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x18\n" +
	"\asummary\x18\x03 \x01(\tR\asummary\x12\x18\n" +
	"\acontent\x18\x04 \x01(\tR\acontent\x12\x10\n" +
	"\x03url\x18\x05 \x01(\tR\x03url\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\tR\tcreatedAt\"\xad\x01\n" +
	"\x1cListCommunityProjectsRequest\x12\x16\n" +
	"\x06filter\x18\x01 \x01(\tR\x06filter\x12\x1d\n" +
	"\n" +
	"sort_field\x18\x02 \x01(\tR\tsortField\x12%\n" +
	"\x0esort_direction\x18\x03 \x01(\tR\rsortDirection\x12\x12\n" +
	"\x04page\x18\x04 \x01(\x05R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x05 \x01(\x05R\bpageSize\"l\n" +
	"\x14CommunityProjectList\x121\n" +
	"\bprojects\x18\x01 \x03(\v2\x15.rsk.CommunityProjectR\bprojects\x12!\n" +
	"\fresult_count\x18\x02 \x01(\x05R\vresultCount\"5\n" +
	"\x12ListArchiveRequest\x12\x1f\n" +
	"\vepisode_ids\x18\x01 \x03(\tR\n" +
	"episodeIds\"1\n" +
	"\vArchiveList\x12\"\n" +
	"\x05items\x18\x01 \x03(\v2\f.rsk.ArchiveR\x05items\"\x9f\x01\n" +
	"\aArchive\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12'\n" +
	"\x0frelated_episode\x18\x03 \x01(\tR\x0erelatedEpisode\x12\x18\n" +
	"\x05files\x18\x04 \x03(\tB\x02\x18\x01R\x05files\x12\x1f\n" +
	"\x05media\x18\x05 \x03(\v2\t.rsk.FileR\x05media\"A\n" +
	"\x04File\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12%\n" +
	"\x0ethumbnail_name\x18\x02 \x01(\tR\rthumbnailName2\xcd\x02\n" +
	"\x10CommunityService\x12\xac\x01\n" +
	"\fListProjects\x12!.rsk.ListCommunityProjectsRequest\x1a\x19.rsk.CommunityProjectList\"^\x92A=\n" +
	"\tcommunity\x12\x19Lists community projects.*\x15listCommunityProjects\x82\xd3\xe4\x93\x02\x18\x12\x16/api/community/project\x12\x89\x01\n" +
	"\vListArchive\x12\x17.rsk.ListArchiveRequest\x1a\x10.rsk.ArchiveList\"O\x92A.\n" +
	"\tcommunity\x12\x14Lists archive items.*\vlistArchive\x82\xd3\xe4\x93\x02\x18\x12\x16/api/community/archiveBj\x92A9\x12\x052\x031.0*\x01\x01r-\n" +
	"\x14Community functions.\x12\x15https://scrimpton.comZ,github.com/warmans/rsk-search/server/gen/apib\x06proto3"

var (
	file_community_proto_rawDescOnce sync.Once
	file_community_proto_rawDescData []byte
)

func file_community_proto_rawDescGZIP() []byte {
	file_community_proto_rawDescOnce.Do(func() {
		file_community_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_community_proto_rawDesc), len(file_community_proto_rawDesc)))
	})
	return file_community_proto_rawDescData
}

var file_community_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_community_proto_goTypes = []any{
	(*CommunityProject)(nil),             // 0: rsk.CommunityProject
	(*ListCommunityProjectsRequest)(nil), // 1: rsk.ListCommunityProjectsRequest
	(*CommunityProjectList)(nil),         // 2: rsk.CommunityProjectList
	(*ListArchiveRequest)(nil),           // 3: rsk.ListArchiveRequest
	(*ArchiveList)(nil),                  // 4: rsk.ArchiveList
	(*Archive)(nil),                      // 5: rsk.Archive
	(*File)(nil),                         // 6: rsk.File
}
var file_community_proto_depIdxs = []int32{
	0, // 0: rsk.CommunityProjectList.projects:type_name -> rsk.CommunityProject
	5, // 1: rsk.ArchiveList.items:type_name -> rsk.Archive
	6, // 2: rsk.Archive.media:type_name -> rsk.File
	1, // 3: rsk.CommunityService.ListProjects:input_type -> rsk.ListCommunityProjectsRequest
	3, // 4: rsk.CommunityService.ListArchive:input_type -> rsk.ListArchiveRequest
	2, // 5: rsk.CommunityService.ListProjects:output_type -> rsk.CommunityProjectList
	4, // 6: rsk.CommunityService.ListArchive:output_type -> rsk.ArchiveList
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_community_proto_init() }
func file_community_proto_init() {
	if File_community_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_community_proto_rawDesc), len(file_community_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_community_proto_goTypes,
		DependencyIndexes: file_community_proto_depIdxs,
		MessageInfos:      file_community_proto_msgTypes,
	}.Build()
	File_community_proto = out.File
	file_community_proto_goTypes = nil
	file_community_proto_depIdxs = nil
}
