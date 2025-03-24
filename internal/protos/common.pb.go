// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.18.1
// source: common.proto

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

type Present struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Present) Reset() {
	*x = Present{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Present) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Present) ProtoMessage() {}

func (x *Present) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Present.ProtoReflect.Descriptor instead.
func (*Present) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

type PresentAndNotEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValueToCheck []byte `protobuf:"bytes,1,opt,name=value_to_check,json=valueToCheck,proto3" json:"value_to_check,omitempty"`
}

func (x *PresentAndNotEqual) Reset() {
	*x = PresentAndNotEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PresentAndNotEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresentAndNotEqual) ProtoMessage() {}

func (x *PresentAndNotEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresentAndNotEqual.ProtoReflect.Descriptor instead.
func (*PresentAndNotEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

func (x *PresentAndNotEqual) GetValueToCheck() []byte {
	if x != nil {
		return x.ValueToCheck
	}
	return nil
}

type Absent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Absent) Reset() {
	*x = Absent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Absent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Absent) ProtoMessage() {}

func (x *Absent) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Absent.ProtoReflect.Descriptor instead.
func (*Absent) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

type Equal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValueToCheck []byte `protobuf:"bytes,1,opt,name=value_to_check,json=valueToCheck,proto3" json:"value_to_check,omitempty"`
}

func (x *Equal) Reset() {
	*x = Equal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Equal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Equal) ProtoMessage() {}

func (x *Equal) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Equal.ProtoReflect.Descriptor instead.
func (*Equal) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{3}
}

func (x *Equal) GetValueToCheck() []byte {
	if x != nil {
		return x.ValueToCheck
	}
	return nil
}

type AbsentOrEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValueToCheck []byte `protobuf:"bytes,1,opt,name=value_to_check,json=valueToCheck,proto3" json:"value_to_check,omitempty"`
}

func (x *AbsentOrEqual) Reset() {
	*x = AbsentOrEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbsentOrEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbsentOrEqual) ProtoMessage() {}

func (x *AbsentOrEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbsentOrEqual.ProtoReflect.Descriptor instead.
func (*AbsentOrEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{4}
}

func (x *AbsentOrEqual) GetValueToCheck() []byte {
	if x != nil {
		return x.ValueToCheck
	}
	return nil
}

type NotEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValueToCheck []byte `protobuf:"bytes,1,opt,name=value_to_check,json=valueToCheck,proto3" json:"value_to_check,omitempty"`
}

func (x *NotEqual) Reset() {
	*x = NotEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotEqual) ProtoMessage() {}

func (x *NotEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotEqual.ProtoReflect.Descriptor instead.
func (*NotEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{5}
}

func (x *NotEqual) GetValueToCheck() []byte {
	if x != nil {
		return x.ValueToCheck
	}
	return nil
}

type XUnbounded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *XUnbounded) Reset() {
	*x = XUnbounded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XUnbounded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XUnbounded) ProtoMessage() {}

func (x *XUnbounded) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XUnbounded.ProtoReflect.Descriptor instead.
func (*XUnbounded) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{6}
}

type XEmpty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *XEmpty) Reset() {
	*x = XEmpty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XEmpty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XEmpty) ProtoMessage() {}

func (x *XEmpty) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XEmpty.ProtoReflect.Descriptor instead.
func (*XEmpty) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{7}
}

type PresentAndNotHashEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashToCheck []byte `protobuf:"bytes,1,opt,name=hash_to_check,json=hashToCheck,proto3" json:"hash_to_check,omitempty"`
}

func (x *PresentAndNotHashEqual) Reset() {
	*x = PresentAndNotHashEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PresentAndNotHashEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresentAndNotHashEqual) ProtoMessage() {}

func (x *PresentAndNotHashEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresentAndNotHashEqual.ProtoReflect.Descriptor instead.
func (*PresentAndNotHashEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{8}
}

func (x *PresentAndNotHashEqual) GetHashToCheck() []byte {
	if x != nil {
		return x.HashToCheck
	}
	return nil
}

type PresentAndHashEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashToCheck []byte `protobuf:"bytes,1,opt,name=hash_to_check,json=hashToCheck,proto3" json:"hash_to_check,omitempty"`
}

func (x *PresentAndHashEqual) Reset() {
	*x = PresentAndHashEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PresentAndHashEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresentAndHashEqual) ProtoMessage() {}

func (x *PresentAndHashEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresentAndHashEqual.ProtoReflect.Descriptor instead.
func (*PresentAndHashEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{9}
}

func (x *PresentAndHashEqual) GetHashToCheck() []byte {
	if x != nil {
		return x.HashToCheck
	}
	return nil
}

type AbsentOrHashEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashToCheck []byte `protobuf:"bytes,1,opt,name=hash_to_check,json=hashToCheck,proto3" json:"hash_to_check,omitempty"`
}

func (x *AbsentOrHashEqual) Reset() {
	*x = AbsentOrHashEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbsentOrHashEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbsentOrHashEqual) ProtoMessage() {}

func (x *AbsentOrHashEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbsentOrHashEqual.ProtoReflect.Descriptor instead.
func (*AbsentOrHashEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{10}
}

func (x *AbsentOrHashEqual) GetHashToCheck() []byte {
	if x != nil {
		return x.HashToCheck
	}
	return nil
}

type AbsentOrNotHashEqual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashToCheck []byte `protobuf:"bytes,1,opt,name=hash_to_check,json=hashToCheck,proto3" json:"hash_to_check,omitempty"`
}

func (x *AbsentOrNotHashEqual) Reset() {
	*x = AbsentOrNotHashEqual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbsentOrNotHashEqual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbsentOrNotHashEqual) ProtoMessage() {}

func (x *AbsentOrNotHashEqual) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbsentOrNotHashEqual.ProtoReflect.Descriptor instead.
func (*AbsentOrNotHashEqual) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{11}
}

func (x *AbsentOrNotHashEqual) GetHashToCheck() []byte {
	if x != nil {
		return x.HashToCheck
	}
	return nil
}

type Unconditional struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Unconditional) Reset() {
	*x = Unconditional{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unconditional) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unconditional) ProtoMessage() {}

func (x *Unconditional) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unconditional.ProtoReflect.Descriptor instead.
func (*Unconditional) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{12}
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0x09, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e,
	0x74, 0x22, 0x3a, 0x0a, 0x12, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x41, 0x6e, 0x64, 0x4e,
	0x6f, 0x74, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x12, 0x24, 0x0a, 0x0e, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0c, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x6f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x08, 0x0a,
	0x06, 0x41, 0x62, 0x73, 0x65, 0x6e, 0x74, 0x22, 0x2d, 0x0a, 0x05, 0x45, 0x71, 0x75, 0x61, 0x6c,
	0x12, 0x24, 0x0a, 0x0e, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x54,
	0x6f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x35, 0x0a, 0x0d, 0x41, 0x62, 0x73, 0x65, 0x6e, 0x74,
	0x4f, 0x72, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x12, 0x24, 0x0a, 0x0e, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0c, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x6f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x30, 0x0a,
	0x08, 0x4e, 0x6f, 0x74, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x12, 0x24, 0x0a, 0x0e, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0c, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x6f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22,
	0x0c, 0x0a, 0x0a, 0x5f, 0x55, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x65, 0x64, 0x22, 0x08, 0x0a,
	0x06, 0x5f, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x3c, 0x0a, 0x16, 0x50, 0x72, 0x65, 0x73, 0x65,
	0x6e, 0x74, 0x41, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68, 0x45, 0x71, 0x75, 0x61,
	0x6c, 0x12, 0x22, 0x0a, 0x0d, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x68, 0x54, 0x6f,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x39, 0x0a, 0x13, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74,
	0x41, 0x6e, 0x64, 0x48, 0x61, 0x73, 0x68, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x12, 0x22, 0x0a, 0x0d,
	0x68, 0x61, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x68, 0x54, 0x6f, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x22, 0x37, 0x0a, 0x11, 0x41, 0x62, 0x73, 0x65, 0x6e, 0x74, 0x4f, 0x72, 0x48, 0x61, 0x73, 0x68,
	0x45, 0x71, 0x75, 0x61, 0x6c, 0x12, 0x22, 0x0a, 0x0d, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74, 0x6f,
	0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x68, 0x61,
	0x73, 0x68, 0x54, 0x6f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x3a, 0x0a, 0x14, 0x41, 0x62, 0x73,
	0x65, 0x6e, 0x74, 0x4f, 0x72, 0x4e, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68, 0x45, 0x71, 0x75, 0x61,
	0x6c, 0x12, 0x22, 0x0a, 0x0d, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x68, 0x54, 0x6f,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x0f, 0x0a, 0x0d, 0x55, 0x6e, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x42, 0x59, 0x0a, 0x0b, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x68, 0x71, 0x2f, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x3b, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x73, 0x64, 0x6b, 0x5f, 0x67, 0x6f, 0xaa, 0x02, 0x15, 0x4d, 0x6f, 0x6d, 0x65,
	0x6e, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_common_proto_goTypes = []any{
	(*Present)(nil),                // 0: common.Present
	(*PresentAndNotEqual)(nil),     // 1: common.PresentAndNotEqual
	(*Absent)(nil),                 // 2: common.Absent
	(*Equal)(nil),                  // 3: common.Equal
	(*AbsentOrEqual)(nil),          // 4: common.AbsentOrEqual
	(*NotEqual)(nil),               // 5: common.NotEqual
	(*XUnbounded)(nil),             // 6: common._Unbounded
	(*XEmpty)(nil),                 // 7: common._Empty
	(*PresentAndNotHashEqual)(nil), // 8: common.PresentAndNotHashEqual
	(*PresentAndHashEqual)(nil),    // 9: common.PresentAndHashEqual
	(*AbsentOrHashEqual)(nil),      // 10: common.AbsentOrHashEqual
	(*AbsentOrNotHashEqual)(nil),   // 11: common.AbsentOrNotHashEqual
	(*Unconditional)(nil),          // 12: common.Unconditional
}
var file_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Present); i {
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
		file_common_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*PresentAndNotEqual); i {
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
		file_common_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Absent); i {
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
		file_common_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Equal); i {
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
		file_common_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*AbsentOrEqual); i {
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
		file_common_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*NotEqual); i {
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
		file_common_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*XUnbounded); i {
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
		file_common_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*XEmpty); i {
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
		file_common_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*PresentAndNotHashEqual); i {
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
		file_common_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*PresentAndHashEqual); i {
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
		file_common_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*AbsentOrHashEqual); i {
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
		file_common_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*AbsentOrNotHashEqual); i {
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
		file_common_proto_msgTypes[12].Exporter = func(v any, i int) any {
			switch v := v.(*Unconditional); i {
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
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}
