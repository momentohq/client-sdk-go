// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.20.3
// source: webhook.proto

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

type XWebhookId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
	// This is limited to 128 chars.
	WebhookName string `protobuf:"bytes,2,opt,name=webhook_name,json=webhookName,proto3" json:"webhook_name,omitempty"`
}

func (x *XWebhookId) Reset() {
	*x = XWebhookId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XWebhookId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XWebhookId) ProtoMessage() {}

func (x *XWebhookId) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XWebhookId.ProtoReflect.Descriptor instead.
func (*XWebhookId) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{0}
}

func (x *XWebhookId) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

func (x *XWebhookId) GetWebhookName() string {
	if x != nil {
		return x.WebhookName
	}
	return ""
}

type XWebhook struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WebhookId   *XWebhookId          `protobuf:"bytes,1,opt,name=webhook_id,json=webhookId,proto3" json:"webhook_id,omitempty"`
	TopicName   string               `protobuf:"bytes,2,opt,name=topic_name,json=topicName,proto3" json:"topic_name,omitempty"`
	Destination *XWebhookDestination `protobuf:"bytes,3,opt,name=destination,proto3" json:"destination,omitempty"`
}

func (x *XWebhook) Reset() {
	*x = XWebhook{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XWebhook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XWebhook) ProtoMessage() {}

func (x *XWebhook) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XWebhook.ProtoReflect.Descriptor instead.
func (*XWebhook) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{1}
}

func (x *XWebhook) GetWebhookId() *XWebhookId {
	if x != nil {
		return x.WebhookId
	}
	return nil
}

func (x *XWebhook) GetTopicName() string {
	if x != nil {
		return x.TopicName
	}
	return ""
}

func (x *XWebhook) GetDestination() *XWebhookDestination {
	if x != nil {
		return x.Destination
	}
	return nil
}

type XPutWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Webhook *XWebhook `protobuf:"bytes,1,opt,name=webhook,proto3" json:"webhook,omitempty"`
}

func (x *XPutWebhookRequest) Reset() {
	*x = XPutWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XPutWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XPutWebhookRequest) ProtoMessage() {}

func (x *XPutWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XPutWebhookRequest.ProtoReflect.Descriptor instead.
func (*XPutWebhookRequest) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{2}
}

func (x *XPutWebhookRequest) GetWebhook() *XWebhook {
	if x != nil {
		return x.Webhook
	}
	return nil
}

type XPutWebhookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SecretString string `protobuf:"bytes,1,opt,name=secret_string,json=secretString,proto3" json:"secret_string,omitempty"`
}

func (x *XPutWebhookResponse) Reset() {
	*x = XPutWebhookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XPutWebhookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XPutWebhookResponse) ProtoMessage() {}

func (x *XPutWebhookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XPutWebhookResponse.ProtoReflect.Descriptor instead.
func (*XPutWebhookResponse) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{3}
}

func (x *XPutWebhookResponse) GetSecretString() string {
	if x != nil {
		return x.SecretString
	}
	return ""
}

type XDeleteWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WebhookId *XWebhookId `protobuf:"bytes,1,opt,name=webhook_id,json=webhookId,proto3" json:"webhook_id,omitempty"`
}

func (x *XDeleteWebhookRequest) Reset() {
	*x = XDeleteWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XDeleteWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XDeleteWebhookRequest) ProtoMessage() {}

func (x *XDeleteWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XDeleteWebhookRequest.ProtoReflect.Descriptor instead.
func (*XDeleteWebhookRequest) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{4}
}

func (x *XDeleteWebhookRequest) GetWebhookId() *XWebhookId {
	if x != nil {
		return x.WebhookId
	}
	return nil
}

type XDeleteWebhookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *XDeleteWebhookResponse) Reset() {
	*x = XDeleteWebhookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XDeleteWebhookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XDeleteWebhookResponse) ProtoMessage() {}

func (x *XDeleteWebhookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XDeleteWebhookResponse.ProtoReflect.Descriptor instead.
func (*XDeleteWebhookResponse) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{5}
}

type XListWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
}

func (x *XListWebhookRequest) Reset() {
	*x = XListWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XListWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XListWebhookRequest) ProtoMessage() {}

func (x *XListWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XListWebhookRequest.ProtoReflect.Descriptor instead.
func (*XListWebhookRequest) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{6}
}

func (x *XListWebhookRequest) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

type XListWebhooksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Webhook []*XWebhook `protobuf:"bytes,1,rep,name=webhook,proto3" json:"webhook,omitempty"`
}

func (x *XListWebhooksResponse) Reset() {
	*x = XListWebhooksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XListWebhooksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XListWebhooksResponse) ProtoMessage() {}

func (x *XListWebhooksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XListWebhooksResponse.ProtoReflect.Descriptor instead.
func (*XListWebhooksResponse) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{7}
}

func (x *XListWebhooksResponse) GetWebhook() []*XWebhook {
	if x != nil {
		return x.Webhook
	}
	return nil
}

type XGetWebhookSecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName   string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
	WebhookName string `protobuf:"bytes,2,opt,name=webhook_name,json=webhookName,proto3" json:"webhook_name,omitempty"`
}

func (x *XGetWebhookSecretRequest) Reset() {
	*x = XGetWebhookSecretRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XGetWebhookSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XGetWebhookSecretRequest) ProtoMessage() {}

func (x *XGetWebhookSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XGetWebhookSecretRequest.ProtoReflect.Descriptor instead.
func (*XGetWebhookSecretRequest) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{8}
}

func (x *XGetWebhookSecretRequest) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

func (x *XGetWebhookSecretRequest) GetWebhookName() string {
	if x != nil {
		return x.WebhookName
	}
	return ""
}

type XGetWebhookSecretResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CacheName    string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
	WebhookName  string `protobuf:"bytes,2,opt,name=webhook_name,json=webhookName,proto3" json:"webhook_name,omitempty"`
	SecretString string `protobuf:"bytes,3,opt,name=secret_string,json=secretString,proto3" json:"secret_string,omitempty"`
}

func (x *XGetWebhookSecretResponse) Reset() {
	*x = XGetWebhookSecretResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XGetWebhookSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XGetWebhookSecretResponse) ProtoMessage() {}

func (x *XGetWebhookSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XGetWebhookSecretResponse.ProtoReflect.Descriptor instead.
func (*XGetWebhookSecretResponse) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{9}
}

func (x *XGetWebhookSecretResponse) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

func (x *XGetWebhookSecretResponse) GetWebhookName() string {
	if x != nil {
		return x.WebhookName
	}
	return ""
}

func (x *XGetWebhookSecretResponse) GetSecretString() string {
	if x != nil {
		return x.SecretString
	}
	return ""
}

type XWebhookDestination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Kind:
	//
	//	*XWebhookDestination_PostUrl
	Kind isXWebhookDestination_Kind `protobuf_oneof:"kind"`
}

func (x *XWebhookDestination) Reset() {
	*x = XWebhookDestination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XWebhookDestination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XWebhookDestination) ProtoMessage() {}

func (x *XWebhookDestination) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XWebhookDestination.ProtoReflect.Descriptor instead.
func (*XWebhookDestination) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{10}
}

func (m *XWebhookDestination) GetKind() isXWebhookDestination_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *XWebhookDestination) GetPostUrl() string {
	if x, ok := x.GetKind().(*XWebhookDestination_PostUrl); ok {
		return x.PostUrl
	}
	return ""
}

type isXWebhookDestination_Kind interface {
	isXWebhookDestination_Kind()
}

type XWebhookDestination_PostUrl struct {
	PostUrl string `protobuf:"bytes,1,opt,name=post_url,json=postUrl,proto3,oneof"`
}

func (*XWebhookDestination_PostUrl) isXWebhookDestination_Kind() {}

type XRotateWebhookSecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WebhookId *XWebhookId `protobuf:"bytes,1,opt,name=webhook_id,json=webhookId,proto3" json:"webhook_id,omitempty"`
}

func (x *XRotateWebhookSecretRequest) Reset() {
	*x = XRotateWebhookSecretRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XRotateWebhookSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XRotateWebhookSecretRequest) ProtoMessage() {}

func (x *XRotateWebhookSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XRotateWebhookSecretRequest.ProtoReflect.Descriptor instead.
func (*XRotateWebhookSecretRequest) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{11}
}

func (x *XRotateWebhookSecretRequest) GetWebhookId() *XWebhookId {
	if x != nil {
		return x.WebhookId
	}
	return nil
}

type XRotateWebhookSecretResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SecretString string `protobuf:"bytes,1,opt,name=secret_string,json=secretString,proto3" json:"secret_string,omitempty"`
}

func (x *XRotateWebhookSecretResponse) Reset() {
	*x = XRotateWebhookSecretResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webhook_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XRotateWebhookSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XRotateWebhookSecretResponse) ProtoMessage() {}

func (x *XRotateWebhookSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webhook_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XRotateWebhookSecretResponse.ProtoReflect.Descriptor instead.
func (*XRotateWebhookSecretResponse) Descriptor() ([]byte, []int) {
	return file_webhook_proto_rawDescGZIP(), []int{12}
}

func (x *XRotateWebhookSecretResponse) GetSecretString() string {
	if x != nil {
		return x.SecretString
	}
	return ""
}

var File_webhook_proto protoreflect.FileDescriptor

var file_webhook_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x22, 0x4e, 0x0a, 0x0a, 0x5f, 0x57, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x9d, 0x01, 0x0a, 0x08, 0x5f, 0x57, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x32, 0x0a, 0x0a, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x52, 0x09,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x41, 0x0a, 0x12, 0x5f, 0x50, 0x75, 0x74,
	0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b,
	0x0a, 0x07, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x57, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x52, 0x07, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x22, 0x3a, 0x0a, 0x13, 0x5f,
	0x50, 0x75, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x5f, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x4b, 0x0a, 0x15, 0x5f, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x32, 0x0a, 0x0a, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f,
	0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x52, 0x09, 0x77, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x49, 0x64, 0x22, 0x18, 0x0a, 0x16, 0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x34,
	0x0a, 0x13, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x44, 0x0a, 0x15, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a,
	0x07, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x52, 0x07, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x22, 0x5c, 0x0a, 0x18, 0x5f, 0x47,
	0x65, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x82, 0x01, 0x0a, 0x19, 0x5f, 0x47, 0x65,
	0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x3a, 0x0a,
	0x13, 0x5f, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x72,
	0x6c, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x51, 0x0a, 0x1b, 0x5f, 0x52, 0x6f,
	0x74, 0x61, 0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x0a, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x77,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49,
	0x64, 0x52, 0x09, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x22, 0x43, 0x0a, 0x1c,
	0x5f, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x32, 0xbb, 0x03, 0x0a, 0x07, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x49, 0x0a,
	0x0a, 0x50, 0x75, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x1b, 0x2e, 0x77, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x50, 0x75, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x2e, 0x5f, 0x50, 0x75, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x1e, 0x2e, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x0c,
	0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x1c, 0x2e, 0x77,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5b, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x21, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x47, 0x65, 0x74, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x47,
	0x65, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x64, 0x0a, 0x13, 0x52, 0x6f, 0x74,
	0x61, 0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x24, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x5f, 0x52, 0x6f, 0x74, 0x61,
	0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x2e, 0x5f, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x5b, 0x0a, 0x0c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x50,
	0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f,
	0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x68, 0x71, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x73,
	0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x3b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x64, 0x6b,
	0x5f, 0x67, 0x6f, 0xaa, 0x02, 0x16, 0x4d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x2e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_webhook_proto_rawDescOnce sync.Once
	file_webhook_proto_rawDescData = file_webhook_proto_rawDesc
)

func file_webhook_proto_rawDescGZIP() []byte {
	file_webhook_proto_rawDescOnce.Do(func() {
		file_webhook_proto_rawDescData = protoimpl.X.CompressGZIP(file_webhook_proto_rawDescData)
	})
	return file_webhook_proto_rawDescData
}

var file_webhook_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_webhook_proto_goTypes = []interface{}{
	(*XWebhookId)(nil),                   // 0: webhook._WebhookId
	(*XWebhook)(nil),                     // 1: webhook._Webhook
	(*XPutWebhookRequest)(nil),           // 2: webhook._PutWebhookRequest
	(*XPutWebhookResponse)(nil),          // 3: webhook._PutWebhookResponse
	(*XDeleteWebhookRequest)(nil),        // 4: webhook._DeleteWebhookRequest
	(*XDeleteWebhookResponse)(nil),       // 5: webhook._DeleteWebhookResponse
	(*XListWebhookRequest)(nil),          // 6: webhook._ListWebhookRequest
	(*XListWebhooksResponse)(nil),        // 7: webhook._ListWebhooksResponse
	(*XGetWebhookSecretRequest)(nil),     // 8: webhook._GetWebhookSecretRequest
	(*XGetWebhookSecretResponse)(nil),    // 9: webhook._GetWebhookSecretResponse
	(*XWebhookDestination)(nil),          // 10: webhook._WebhookDestination
	(*XRotateWebhookSecretRequest)(nil),  // 11: webhook._RotateWebhookSecretRequest
	(*XRotateWebhookSecretResponse)(nil), // 12: webhook._RotateWebhookSecretResponse
}
var file_webhook_proto_depIdxs = []int32{
	0,  // 0: webhook._Webhook.webhook_id:type_name -> webhook._WebhookId
	10, // 1: webhook._Webhook.destination:type_name -> webhook._WebhookDestination
	1,  // 2: webhook._PutWebhookRequest.webhook:type_name -> webhook._Webhook
	0,  // 3: webhook._DeleteWebhookRequest.webhook_id:type_name -> webhook._WebhookId
	1,  // 4: webhook._ListWebhooksResponse.webhook:type_name -> webhook._Webhook
	0,  // 5: webhook._RotateWebhookSecretRequest.webhook_id:type_name -> webhook._WebhookId
	2,  // 6: webhook.Webhook.PutWebhook:input_type -> webhook._PutWebhookRequest
	4,  // 7: webhook.Webhook.DeleteWebhook:input_type -> webhook._DeleteWebhookRequest
	6,  // 8: webhook.Webhook.ListWebhooks:input_type -> webhook._ListWebhookRequest
	8,  // 9: webhook.Webhook.GetWebhookSecret:input_type -> webhook._GetWebhookSecretRequest
	11, // 10: webhook.Webhook.RotateWebhookSecret:input_type -> webhook._RotateWebhookSecretRequest
	3,  // 11: webhook.Webhook.PutWebhook:output_type -> webhook._PutWebhookResponse
	5,  // 12: webhook.Webhook.DeleteWebhook:output_type -> webhook._DeleteWebhookResponse
	7,  // 13: webhook.Webhook.ListWebhooks:output_type -> webhook._ListWebhooksResponse
	9,  // 14: webhook.Webhook.GetWebhookSecret:output_type -> webhook._GetWebhookSecretResponse
	12, // 15: webhook.Webhook.RotateWebhookSecret:output_type -> webhook._RotateWebhookSecretResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_webhook_proto_init() }
func file_webhook_proto_init() {
	if File_webhook_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_webhook_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XWebhookId); i {
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
		file_webhook_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XWebhook); i {
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
		file_webhook_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XPutWebhookRequest); i {
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
		file_webhook_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XPutWebhookResponse); i {
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
		file_webhook_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XDeleteWebhookRequest); i {
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
		file_webhook_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XDeleteWebhookResponse); i {
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
		file_webhook_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XListWebhookRequest); i {
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
		file_webhook_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XListWebhooksResponse); i {
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
		file_webhook_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XGetWebhookSecretRequest); i {
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
		file_webhook_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XGetWebhookSecretResponse); i {
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
		file_webhook_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XWebhookDestination); i {
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
		file_webhook_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XRotateWebhookSecretRequest); i {
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
		file_webhook_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XRotateWebhookSecretResponse); i {
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
	file_webhook_proto_msgTypes[10].OneofWrappers = []interface{}{
		(*XWebhookDestination_PostUrl)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_webhook_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_webhook_proto_goTypes,
		DependencyIndexes: file_webhook_proto_depIdxs,
		MessageInfos:      file_webhook_proto_msgTypes,
	}.Build()
	File_webhook_proto = out.File
	file_webhook_proto_rawDesc = nil
	file_webhook_proto_goTypes = nil
	file_webhook_proto_depIdxs = nil
}
