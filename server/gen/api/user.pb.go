// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: user.proto

package api

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Notification_NotificationKind int32

const (
	Notification_UNDEFINED_KIND Notification_NotificationKind = 0
	Notification_CONFIRMATION   Notification_NotificationKind = 1
	Notification_INFO           Notification_NotificationKind = 2
	Notification_WARNING        Notification_NotificationKind = 3
	Notification_SPAM           Notification_NotificationKind = 4
)

// Enum value maps for Notification_NotificationKind.
var (
	Notification_NotificationKind_name = map[int32]string{
		0: "UNDEFINED_KIND",
		1: "CONFIRMATION",
		2: "INFO",
		3: "WARNING",
		4: "SPAM",
	}
	Notification_NotificationKind_value = map[string]int32{
		"UNDEFINED_KIND": 0,
		"CONFIRMATION":   1,
		"INFO":           2,
		"WARNING":        3,
		"SPAM":           4,
	}
)

func (x Notification_NotificationKind) Enum() *Notification_NotificationKind {
	p := new(Notification_NotificationKind)
	*p = x
	return p
}

func (x Notification_NotificationKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Notification_NotificationKind) Descriptor() protoreflect.EnumDescriptor {
	return file_user_proto_enumTypes[0].Descriptor()
}

func (Notification_NotificationKind) Type() protoreflect.EnumType {
	return &file_user_proto_enumTypes[0]
}

func (x Notification_NotificationKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Notification_NotificationKind.Descriptor instead.
func (Notification_NotificationKind) EnumDescriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2, 0}
}

type ListNotificationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filter        string `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	SortField     string `protobuf:"bytes,2,opt,name=sort_field,json=sortField,proto3" json:"sort_field,omitempty"`
	SortDirection string `protobuf:"bytes,3,opt,name=sort_direction,json=sortDirection,proto3" json:"sort_direction,omitempty"`
	Page          int32  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32  `protobuf:"varint,5,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *ListNotificationsRequest) Reset() {
	*x = ListNotificationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListNotificationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNotificationsRequest) ProtoMessage() {}

func (x *ListNotificationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNotificationsRequest.ProtoReflect.Descriptor instead.
func (*ListNotificationsRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *ListNotificationsRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *ListNotificationsRequest) GetSortField() string {
	if x != nil {
		return x.SortField
	}
	return ""
}

func (x *ListNotificationsRequest) GetSortDirection() string {
	if x != nil {
		return x.SortDirection
	}
	return ""
}

func (x *ListNotificationsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListNotificationsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type NotificationsList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Notifications []*Notification `protobuf:"bytes,1,rep,name=notifications,proto3" json:"notifications,omitempty"`
}

func (x *NotificationsList) Reset() {
	*x = NotificationsList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationsList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationsList) ProtoMessage() {}

func (x *NotificationsList) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationsList.ProtoReflect.Descriptor instead.
func (*NotificationsList) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *NotificationsList) GetNotifications() []*Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

type Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Kind           Notification_NotificationKind `protobuf:"varint,2,opt,name=kind,proto3,enum=rsk.Notification_NotificationKind" json:"kind,omitempty"`
	Message        string                        `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	ClickThoughUrl string                        `protobuf:"bytes,4,opt,name=click_though_url,json=clickThoughUrl,proto3" json:"click_though_url,omitempty"`
	CreatedAt      string                        `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ReadAt         string                        `protobuf:"bytes,6,opt,name=read_at,json=readAt,proto3" json:"read_at,omitempty"`
}

func (x *Notification) Reset() {
	*x = Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *Notification) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Notification) GetKind() Notification_NotificationKind {
	if x != nil {
		return x.Kind
	}
	return Notification_UNDEFINED_KIND
}

func (x *Notification) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Notification) GetClickThoughUrl() string {
	if x != nil {
		return x.ClickThoughUrl
	}
	return ""
}

func (x *Notification) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Notification) GetReadAt() string {
	if x != nil {
		return x.ReadAt
	}
	return ""
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x72, 0x73,
	0x6b, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xa9, 0x01, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6f, 0x72, 0x74, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x6f, 0x72,
	0x74, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x4c, 0x0a, 0x11, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x37, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x73, 0x6b, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xad, 0x02, 0x0a, 0x0c, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x36, 0x0a, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x72, 0x73, 0x6b, 0x2e, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x10,
	0x63, 0x6c, 0x69, 0x63, 0x6b, 0x5f, 0x74, 0x68, 0x6f, 0x75, 0x67, 0x68, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x54, 0x68, 0x6f,
	0x75, 0x67, 0x68, 0x55, 0x72, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x64, 0x41, 0x74, 0x22, 0x59,
	0x0a, 0x10, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x69,
	0x6e, 0x64, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x5f,
	0x4b, 0x49, 0x4e, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x52,
	0x4d, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x46, 0x4f,
	0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12,
	0x08, 0x0a, 0x04, 0x53, 0x50, 0x41, 0x4d, 0x10, 0x04, 0x32, 0xba, 0x01, 0x0a, 0x0b, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xaa, 0x01, 0x0a, 0x11, 0x4c, 0x69,
	0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x1d, 0x2e, 0x72, 0x73, 0x6b, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x72, 0x73, 0x6b, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x5e, 0x92, 0x41, 0x3c, 0x0a, 0x06, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x12, 0x1f, 0x4c, 0x69, 0x73, 0x74, 0x20, 0x6d, 0x6f, 0x73, 0x74, 0x20, 0x72,
	0x65, 0x63, 0x65, 0x6e, 0x74, 0x20, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x2a, 0x11, 0x6c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x12, 0x17, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x8e, 0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x72, 0x6d, 0x61, 0x6e, 0x73, 0x2f, 0x72, 0x73,
	0x6b, 0x2d, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x92, 0x41, 0x5d, 0x12, 0x05, 0x32, 0x03, 0x31, 0x2e,
	0x30, 0x2a, 0x01, 0x01, 0x72, 0x51, 0x0a, 0x38, 0x55, 0x73, 0x65, 0x72, 0x20, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x20, 0x68, 0x61, 0x73, 0x20, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x20, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x61, 0x20,
	0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x75, 0x6c, 0x61, 0x72, 0x20, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x12, 0x15, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x73, 0x63, 0x72, 0x69, 0x6d, 0x70,
	0x74, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_user_proto_goTypes = []interface{}{
	(Notification_NotificationKind)(0), // 0: rsk.Notification.NotificationKind
	(*ListNotificationsRequest)(nil),   // 1: rsk.ListNotificationsRequest
	(*NotificationsList)(nil),          // 2: rsk.NotificationsList
	(*Notification)(nil),               // 3: rsk.Notification
}
var file_user_proto_depIdxs = []int32{
	3, // 0: rsk.NotificationsList.notifications:type_name -> rsk.Notification
	0, // 1: rsk.Notification.kind:type_name -> rsk.Notification.NotificationKind
	1, // 2: rsk.UserService.ListNotifications:input_type -> rsk.ListNotificationsRequest
	2, // 3: rsk.UserService.ListNotifications:output_type -> rsk.NotificationsList
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListNotificationsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationsList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notification); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		EnumInfos:         file_user_proto_enumTypes,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
