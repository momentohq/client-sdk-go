// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.20.3
// source: token.proto

package client_sdk_go

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type XGenerateDisposableTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expires     *XGenerateDisposableTokenRequest_Expires `protobuf:"bytes,1,opt,name=expires,proto3" json:"expires,omitempty"`
	AuthToken   string                                   `protobuf:"bytes,2,opt,name=auth_token,json=authToken,proto3" json:"auth_token,omitempty"`
	Permissions *Permissions                             `protobuf:"bytes,3,opt,name=permissions,proto3" json:"permissions,omitempty"`
}

func (x *XGenerateDisposableTokenRequest) Reset() {
	*x = XGenerateDisposableTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XGenerateDisposableTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XGenerateDisposableTokenRequest) ProtoMessage() {}

func (x *XGenerateDisposableTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XGenerateDisposableTokenRequest.ProtoReflect.Descriptor instead.
func (*XGenerateDisposableTokenRequest) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{0}
}

func (x *XGenerateDisposableTokenRequest) GetExpires() *XGenerateDisposableTokenRequest_Expires {
	if x != nil {
		return x.Expires
	}
	return nil
}

func (x *XGenerateDisposableTokenRequest) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *XGenerateDisposableTokenRequest) GetPermissions() *Permissions {
	if x != nil {
		return x.Permissions
	}
	return nil
}

type XGenerateDisposableTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the new api key used for authentication against Momento backend
	ApiKey string `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	// the Momento endpoint that this token is allowed to make requests against
	Endpoint string `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// epoch seconds when the api token expires
	ValidUntil uint64 `protobuf:"varint,3,opt,name=valid_until,json=validUntil,proto3" json:"valid_until,omitempty"`
}

func (x *XGenerateDisposableTokenResponse) Reset() {
	*x = XGenerateDisposableTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XGenerateDisposableTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XGenerateDisposableTokenResponse) ProtoMessage() {}

func (x *XGenerateDisposableTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XGenerateDisposableTokenResponse.ProtoReflect.Descriptor instead.
func (*XGenerateDisposableTokenResponse) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{1}
}

func (x *XGenerateDisposableTokenResponse) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *XGenerateDisposableTokenResponse) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *XGenerateDisposableTokenResponse) GetValidUntil() uint64 {
	if x != nil {
		return x.ValidUntil
	}
	return 0
}

// generate a token that has an expiry
type XGenerateDisposableTokenRequest_Expires struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// how many seconds do you want the api token to be valid for?
	ValidForSeconds uint32 `protobuf:"varint,1,opt,name=valid_for_seconds,json=validForSeconds,proto3" json:"valid_for_seconds,omitempty"`
}

func (x *XGenerateDisposableTokenRequest_Expires) Reset() {
	*x = XGenerateDisposableTokenRequest_Expires{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XGenerateDisposableTokenRequest_Expires) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XGenerateDisposableTokenRequest_Expires) ProtoMessage() {}

func (x *XGenerateDisposableTokenRequest_Expires) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XGenerateDisposableTokenRequest_Expires.ProtoReflect.Descriptor instead.
func (*XGenerateDisposableTokenRequest_Expires) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{0, 0}
}

func (x *XGenerateDisposableTokenRequest_Expires) GetValidForSeconds() uint32 {
	if x != nil {
		return x.ValidForSeconds
	}
	return 0
}

var File_token_proto protoreflect.FileDescriptor

var file_token_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x1a, 0x18, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85,
	0x02, 0x0a, 0x1f, 0x5f, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x70,
	0x6f, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x48, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x5f, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x52, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a,
	0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x42, 0x0a, 0x0b, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x1a,
	0x35, 0x0a, 0x07, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x46, 0x6f, 0x72, 0x53,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x78, 0x0a, 0x20, 0x5f, 0x47, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x70,
	0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69,
	0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x75, 0x6e, 0x74, 0x69, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x55, 0x6e, 0x74, 0x69, 0x6c,
	0x32, 0x75, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x6c, 0x0a, 0x17, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x26, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x5f, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x70, 0x6f, 0x73, 0x61, 0x62, 0x6c, 0x65,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x5f, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x44, 0x69,
	0x73, 0x70, 0x6f, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x60, 0x0a, 0x0d, 0x6d, 0x6f, 0x6d, 0x65, 0x6e,
	0x74, 0x6f, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x68, 0x71,
	0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x3b, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x64, 0x6b, 0x5f, 0x67, 0x6f, 0xaa, 0x02, 0x1a, 0x4d,
	0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_token_proto_rawDescOnce sync.Once
	file_token_proto_rawDescData = file_token_proto_rawDesc
)

func file_token_proto_rawDescGZIP() []byte {
	file_token_proto_rawDescOnce.Do(func() {
		file_token_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_proto_rawDescData)
	})
	return file_token_proto_rawDescData
}

var file_token_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_token_proto_goTypes = []interface{}{
	(*XGenerateDisposableTokenRequest)(nil),         // 0: token._GenerateDisposableTokenRequest
	(*XGenerateDisposableTokenResponse)(nil),        // 1: token._GenerateDisposableTokenResponse
	(*XGenerateDisposableTokenRequest_Expires)(nil), // 2: token._GenerateDisposableTokenRequest.Expires
	(*Permissions)(nil),                             // 3: permission_messages.Permissions
}
var file_token_proto_depIdxs = []int32{
	2, // 0: token._GenerateDisposableTokenRequest.expires:type_name -> token._GenerateDisposableTokenRequest.Expires
	3, // 1: token._GenerateDisposableTokenRequest.permissions:type_name -> permission_messages.Permissions
	0, // 2: token.Token.GenerateDisposableToken:input_type -> token._GenerateDisposableTokenRequest
	1, // 3: token.Token.GenerateDisposableToken:output_type -> token._GenerateDisposableTokenResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_token_proto_init() }
func file_token_proto_init() {
	if File_token_proto != nil {
		return
	}
	file_permissionmessages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_token_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XGenerateDisposableTokenRequest); i {
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
		file_token_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XGenerateDisposableTokenResponse); i {
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
		file_token_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XGenerateDisposableTokenRequest_Expires); i {
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
			RawDescriptor: file_token_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_token_proto_goTypes,
		DependencyIndexes: file_token_proto_depIdxs,
		MessageInfos:      file_token_proto_msgTypes,
	}.Build()
	File_token_proto = out.File
	file_token_proto_rawDesc = nil
	file_token_proto_goTypes = nil
	file_token_proto_depIdxs = nil
}
