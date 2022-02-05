// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: protos/controlclient.proto

package client_sdk_go

import (
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

type XDeleteCacheRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
}

func (x *XDeleteCacheRequest) Reset() {
	*x = XDeleteCacheRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XDeleteCacheRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XDeleteCacheRequest) ProtoMessage() {}

func (x *XDeleteCacheRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XDeleteCacheRequest.ProtoReflect.Descriptor instead.
func (*XDeleteCacheRequest) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{0}
}

func (x *XDeleteCacheRequest) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

type XDeleteCacheResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *XDeleteCacheResponse) Reset() {
	*x = XDeleteCacheResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XDeleteCacheResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XDeleteCacheResponse) ProtoMessage() {}

func (x *XDeleteCacheResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XDeleteCacheResponse.ProtoReflect.Descriptor instead.
func (*XDeleteCacheResponse) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{1}
}

type XCreateCacheRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
}

func (x *XCreateCacheRequest) Reset() {
	*x = XCreateCacheRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XCreateCacheRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XCreateCacheRequest) ProtoMessage() {}

func (x *XCreateCacheRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XCreateCacheRequest.ProtoReflect.Descriptor instead.
func (*XCreateCacheRequest) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{2}
}

func (x *XCreateCacheRequest) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

type XCreateCacheResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *XCreateCacheResponse) Reset() {
	*x = XCreateCacheResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XCreateCacheResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XCreateCacheResponse) ProtoMessage() {}

func (x *XCreateCacheResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XCreateCacheResponse.ProtoReflect.Descriptor instead.
func (*XCreateCacheResponse) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{3}
}

type XListCachesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NextToken string `protobuf:"bytes,1,opt,name=next_token,json=nextToken,proto3" json:"next_token,omitempty"`
}

func (x *XListCachesRequest) Reset() {
	*x = XListCachesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XListCachesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XListCachesRequest) ProtoMessage() {}

func (x *XListCachesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XListCachesRequest.ProtoReflect.Descriptor instead.
func (*XListCachesRequest) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{4}
}

func (x *XListCachesRequest) GetNextToken() string {
	if x != nil {
		return x.NextToken
	}
	return ""
}

type XCache struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
}

func (x *XCache) Reset() {
	*x = XCache{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XCache) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XCache) ProtoMessage() {}

func (x *XCache) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XCache.ProtoReflect.Descriptor instead.
func (*XCache) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{5}
}

func (x *XCache) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

type XListCachesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cache     []*XCache `protobuf:"bytes,1,rep,name=cache,proto3" json:"cache,omitempty"`
	NextToken string    `protobuf:"bytes,2,opt,name=next_token,json=nextToken,proto3" json:"next_token,omitempty"`
}

func (x *XListCachesResponse) Reset() {
	*x = XListCachesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_controlclient_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XListCachesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XListCachesResponse) ProtoMessage() {}

func (x *XListCachesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_controlclient_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XListCachesResponse.ProtoReflect.Descriptor instead.
func (*XListCachesResponse) Descriptor() ([]byte, []int) {
	return file_protos_controlclient_proto_rawDescGZIP(), []int{6}
}

func (x *XListCachesResponse) GetCache() []*XCache {
	if x != nil {
		return x.Cache
	}
	return nil
}

func (x *XListCachesResponse) GetNextToken() string {
	if x != nil {
		return x.NextToken
	}
	return ""
}

var File_protos_controlclient_proto protoreflect.FileDescriptor

var file_protos_controlclient_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x34, 0x0a, 0x13,
	0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x63,
	0x68, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x34, 0x0a, 0x13, 0x5f, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x16, 0x0a, 0x14, 0x5f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x33, 0x0a, 0x12, 0x5f, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x61, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x27, 0x0a,
	0x06, 0x5f, 0x43, 0x61, 0x63, 0x68, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x62, 0x0a, 0x13, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a,
	0x05, 0x63, 0x61, 0x63, 0x68, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x5f, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x52, 0x05, 0x63, 0x61, 0x63, 0x68, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e,
	0x65, 0x78, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0x9d, 0x02, 0x0a, 0x0a, 0x53,
	0x63, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x5a, 0x0a, 0x0b, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x12, 0x23, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x5f, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x5f,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x12, 0x23, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x63,
	0x68, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x5f, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x57, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x63, 0x68, 0x65, 0x73, 0x12,
	0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2e, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x63, 0x68, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x49, 0x0a, 0x13, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x68, 0x71, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x3b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73,
	0x64, 0x6b, 0x5f, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_controlclient_proto_rawDescOnce sync.Once
	file_protos_controlclient_proto_rawDescData = file_protos_controlclient_proto_rawDesc
)

func file_protos_controlclient_proto_rawDescGZIP() []byte {
	file_protos_controlclient_proto_rawDescOnce.Do(func() {
		file_protos_controlclient_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_controlclient_proto_rawDescData)
	})
	return file_protos_controlclient_proto_rawDescData
}

var file_protos_controlclient_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_protos_controlclient_proto_goTypes = []interface{}{
	(*XDeleteCacheRequest)(nil),  // 0: control_client._DeleteCacheRequest
	(*XDeleteCacheResponse)(nil), // 1: control_client._DeleteCacheResponse
	(*XCreateCacheRequest)(nil),  // 2: control_client._CreateCacheRequest
	(*XCreateCacheResponse)(nil), // 3: control_client._CreateCacheResponse
	(*XListCachesRequest)(nil),   // 4: control_client._ListCachesRequest
	(*XCache)(nil),               // 5: control_client._Cache
	(*XListCachesResponse)(nil),  // 6: control_client._ListCachesResponse
}
var file_protos_controlclient_proto_depIdxs = []int32{
	5, // 0: control_client._ListCachesResponse.cache:type_name -> control_client._Cache
	2, // 1: control_client.ScsControl.CreateCache:input_type -> control_client._CreateCacheRequest
	0, // 2: control_client.ScsControl.DeleteCache:input_type -> control_client._DeleteCacheRequest
	4, // 3: control_client.ScsControl.ListCaches:input_type -> control_client._ListCachesRequest
	3, // 4: control_client.ScsControl.CreateCache:output_type -> control_client._CreateCacheResponse
	1, // 5: control_client.ScsControl.DeleteCache:output_type -> control_client._DeleteCacheResponse
	6, // 6: control_client.ScsControl.ListCaches:output_type -> control_client._ListCachesResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_controlclient_proto_init() }
func file_protos_controlclient_proto_init() {
	if File_protos_controlclient_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_controlclient_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XDeleteCacheRequest); i {
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
		file_protos_controlclient_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XDeleteCacheResponse); i {
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
		file_protos_controlclient_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XCreateCacheRequest); i {
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
		file_protos_controlclient_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XCreateCacheResponse); i {
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
		file_protos_controlclient_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XListCachesRequest); i {
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
		file_protos_controlclient_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XCache); i {
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
		file_protos_controlclient_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XListCachesResponse); i {
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
			RawDescriptor: file_protos_controlclient_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_controlclient_proto_goTypes,
		DependencyIndexes: file_protos_controlclient_proto_depIdxs,
		MessageInfos:      file_protos_controlclient_proto_msgTypes,
	}.Build()
	File_protos_controlclient_proto = out.File
	file_protos_controlclient_proto_rawDesc = nil
	file_protos_controlclient_proto_goTypes = nil
	file_protos_controlclient_proto_depIdxs = nil
}
