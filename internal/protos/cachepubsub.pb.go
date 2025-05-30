// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.18.1
// source: cachepubsub.proto

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

// A value to publish through a topic.
type XPublishRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Cache namespace for the topic to which you want to send the value.
	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
	// The literal topic name to which you want to send the value.
	Topic string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	// The value you want to send to the topic. All current subscribers will receive
	// this, should the whims of the Internet prove merciful.
	Value *XTopicValue `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *XPublishRequest) Reset() {
	*x = XPublishRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XPublishRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XPublishRequest) ProtoMessage() {}

func (x *XPublishRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XPublishRequest.ProtoReflect.Descriptor instead.
func (*XPublishRequest) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{0}
}

func (x *XPublishRequest) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

func (x *XPublishRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *XPublishRequest) GetValue() *XTopicValue {
	if x != nil {
		return x.Value
	}
	return nil
}

// A description of how you want to subscribe to a topic.
type XSubscriptionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Cache namespace for the topic to which you want to subscribe.
	CacheName string `protobuf:"bytes,1,opt,name=cache_name,json=cacheName,proto3" json:"cache_name,omitempty"`
	// The literal topic name to which you want to subscribe.
	Topic string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	// If provided, attempt to reconnect to the topic and replay messages starting from
	// the provided sequence number. You may get a discontinuity if some (or all) of the
	// messages are not available.
	// If provided at 1, you may receive some messages leading up to whatever the
	// newest message is. The exact amount is unspecified and subject to change.
	// If not provided (or 0), the subscription will begin with the latest messages.
	ResumeAtTopicSequenceNumber uint64 `protobuf:"varint,3,opt,name=resume_at_topic_sequence_number,json=resumeAtTopicSequenceNumber,proto3" json:"resume_at_topic_sequence_number,omitempty"`
	// Determined by the service when a topic is created. This clarifies the intent of
	// a subscription, and ensures the right messages are sent for a given
	// `resume_at_topic_sequence_number`.
	// Include this in your Subscribe() calls when you are reconnecting. The right value
	// is the last sequence_page you received.
	SequencePage uint64 `protobuf:"varint,4,opt,name=sequence_page,json=sequencePage,proto3" json:"sequence_page,omitempty"`
}

func (x *XSubscriptionRequest) Reset() {
	*x = XSubscriptionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XSubscriptionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XSubscriptionRequest) ProtoMessage() {}

func (x *XSubscriptionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XSubscriptionRequest.ProtoReflect.Descriptor instead.
func (*XSubscriptionRequest) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{1}
}

func (x *XSubscriptionRequest) GetCacheName() string {
	if x != nil {
		return x.CacheName
	}
	return ""
}

func (x *XSubscriptionRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *XSubscriptionRequest) GetResumeAtTopicSequenceNumber() uint64 {
	if x != nil {
		return x.ResumeAtTopicSequenceNumber
	}
	return 0
}

func (x *XSubscriptionRequest) GetSequencePage() uint64 {
	if x != nil {
		return x.SequencePage
	}
	return 0
}

// Possible message kinds from a topic. They can be items when they're from you, or
// other kinds when we have something we think you might need to know about the
// subscription's status.
type XSubscriptionItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Kind:
	//
	//	*XSubscriptionItem_Item
	//	*XSubscriptionItem_Discontinuity
	//	*XSubscriptionItem_Heartbeat
	Kind isXSubscriptionItem_Kind `protobuf_oneof:"kind"`
}

func (x *XSubscriptionItem) Reset() {
	*x = XSubscriptionItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XSubscriptionItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XSubscriptionItem) ProtoMessage() {}

func (x *XSubscriptionItem) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XSubscriptionItem.ProtoReflect.Descriptor instead.
func (*XSubscriptionItem) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{2}
}

func (m *XSubscriptionItem) GetKind() isXSubscriptionItem_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *XSubscriptionItem) GetItem() *XTopicItem {
	if x, ok := x.GetKind().(*XSubscriptionItem_Item); ok {
		return x.Item
	}
	return nil
}

func (x *XSubscriptionItem) GetDiscontinuity() *XDiscontinuity {
	if x, ok := x.GetKind().(*XSubscriptionItem_Discontinuity); ok {
		return x.Discontinuity
	}
	return nil
}

func (x *XSubscriptionItem) GetHeartbeat() *XHeartbeat {
	if x, ok := x.GetKind().(*XSubscriptionItem_Heartbeat); ok {
		return x.Heartbeat
	}
	return nil
}

type isXSubscriptionItem_Kind interface {
	isXSubscriptionItem_Kind()
}

type XSubscriptionItem_Item struct {
	// The subscription has yielded an item you previously published.
	Item *XTopicItem `protobuf:"bytes,1,opt,name=item,proto3,oneof"`
}

type XSubscriptionItem_Discontinuity struct {
	// Momento wants to let you know we detected some possible inconsistency at this
	// point in the subscription stream.
	//
	// A lack of a discontinuity does not mean the subscription is guaranteed to be
	// strictly perfect, but the presence of a discontinuity is very likely to
	Discontinuity *XDiscontinuity `protobuf:"bytes,2,opt,name=discontinuity,proto3,oneof"`
}

type XSubscriptionItem_Heartbeat struct {
	// The stream is still working, there's nothing to see here.
	Heartbeat *XHeartbeat `protobuf:"bytes,3,opt,name=heartbeat,proto3,oneof"`
}

func (*XSubscriptionItem_Item) isXSubscriptionItem_Kind() {}

func (*XSubscriptionItem_Discontinuity) isXSubscriptionItem_Kind() {}

func (*XSubscriptionItem_Heartbeat) isXSubscriptionItem_Kind() {}

// Your subscription has yielded an item you previously published. Here it is!
type XTopicItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Topic sequence numbers give an order of messages per-topic.
	// All subscribers to a topic will receive messages in the same order, with the same sequence numbers.
	TopicSequenceNumber uint64 `protobuf:"varint,1,opt,name=topic_sequence_number,json=topicSequenceNumber,proto3" json:"topic_sequence_number,omitempty"`
	// The value you previously published to this topic.
	Value *XTopicValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// Authenticated id from Publisher's disposable token
	PublisherId string `protobuf:"bytes,3,opt,name=publisher_id,json=publisherId,proto3" json:"publisher_id,omitempty"`
	// Sequence pages exist to determine which sequence number range a message belongs to. On a topic reset,
	// the sequence numbers reset and a new sequence_page is given.
	// For a given sequence_page, the next message in a topic is topic_sequence_number + 1.
	//
	// Later sequence pages are numbered greater than earlier pages, but they don't go one-by-one.
	SequencePage uint64 `protobuf:"varint,4,opt,name=sequence_page,json=sequencePage,proto3" json:"sequence_page,omitempty"`
}

func (x *XTopicItem) Reset() {
	*x = XTopicItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XTopicItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XTopicItem) ProtoMessage() {}

func (x *XTopicItem) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XTopicItem.ProtoReflect.Descriptor instead.
func (*XTopicItem) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{3}
}

func (x *XTopicItem) GetTopicSequenceNumber() uint64 {
	if x != nil {
		return x.TopicSequenceNumber
	}
	return 0
}

func (x *XTopicItem) GetValue() *XTopicValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *XTopicItem) GetPublisherId() string {
	if x != nil {
		return x.PublisherId
	}
	return ""
}

func (x *XTopicItem) GetSequencePage() uint64 {
	if x != nil {
		return x.SequencePage
	}
	return 0
}

// A value in a topic - published, duplicated and received in a subscription.
type XTopicValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types of messages a topic may relay. You can mix types or you can make conventionally
	// typed topics. Sticking with one kind will generally make your software easier to work
	// with though, so we recommend picking the kind you like and using it for a topic!
	//
	// Types that are assignable to Kind:
	//
	//	*XTopicValue_Text
	//	*XTopicValue_Binary
	Kind isXTopicValue_Kind `protobuf_oneof:"kind"`
}

func (x *XTopicValue) Reset() {
	*x = XTopicValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XTopicValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XTopicValue) ProtoMessage() {}

func (x *XTopicValue) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XTopicValue.ProtoReflect.Descriptor instead.
func (*XTopicValue) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{4}
}

func (m *XTopicValue) GetKind() isXTopicValue_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *XTopicValue) GetText() string {
	if x, ok := x.GetKind().(*XTopicValue_Text); ok {
		return x.Text
	}
	return ""
}

func (x *XTopicValue) GetBinary() []byte {
	if x, ok := x.GetKind().(*XTopicValue_Binary); ok {
		return x.Binary
	}
	return nil
}

type isXTopicValue_Kind interface {
	isXTopicValue_Kind()
}

type XTopicValue_Text struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3,oneof"`
}

type XTopicValue_Binary struct {
	Binary []byte `protobuf:"bytes,2,opt,name=binary,proto3,oneof"`
}

func (*XTopicValue_Text) isXTopicValue_Kind() {}

func (*XTopicValue_Binary) isXTopicValue_Kind() {}

// A message from Momento when we know a subscription to have skipped some messages.
// We don't terminate your subscription in that case, but just in case you care, we
// do our best to let you know about it.
type XDiscontinuity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The last topic value sequence number known to have been attempted (if known, 0 otherwise).
	LastTopicSequence uint64 `protobuf:"varint,1,opt,name=last_topic_sequence,json=lastTopicSequence,proto3" json:"last_topic_sequence,omitempty"`
	// The new topic sequence number after which TopicItems will ostensibly resume.
	NewTopicSequence uint64 `protobuf:"varint,2,opt,name=new_topic_sequence,json=newTopicSequence,proto3" json:"new_topic_sequence,omitempty"`
	// The new topic sequence_page. If you had one before and this one is different, then your topic reset.
	// If you didn't have one, then this is just telling you what the sequence page is expected to be.
	// If you had one before, and this one is the same, then it's just telling you that you missed some messages
	// in the topic. Probably your client is consuming messages a little too slowly in this case!
	NewSequencePage uint64 `protobuf:"varint,3,opt,name=new_sequence_page,json=newSequencePage,proto3" json:"new_sequence_page,omitempty"`
}

func (x *XDiscontinuity) Reset() {
	*x = XDiscontinuity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XDiscontinuity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XDiscontinuity) ProtoMessage() {}

func (x *XDiscontinuity) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XDiscontinuity.ProtoReflect.Descriptor instead.
func (*XDiscontinuity) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{5}
}

func (x *XDiscontinuity) GetLastTopicSequence() uint64 {
	if x != nil {
		return x.LastTopicSequence
	}
	return 0
}

func (x *XDiscontinuity) GetNewTopicSequence() uint64 {
	if x != nil {
		return x.NewTopicSequence
	}
	return 0
}

func (x *XDiscontinuity) GetNewSequencePage() uint64 {
	if x != nil {
		return x.NewSequencePage
	}
	return 0
}

// A message from Momento for when we want to reassure clients or frameworks that a
// Subscription is still healthy.
// These are synthetic meta-events and do not increase the topic sequence count.
// Different subscribers may receive a different cadence of heartbeat, and no guarantee
// is made about the cadence or even presence or absence of heartbeats in a stream.
// They are a tool for helping ensure that socket timeouts and the like don't impact
// subscriptions you may care about, but that aren't receiving a substantial publish rate.
type XHeartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *XHeartbeat) Reset() {
	*x = XHeartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cachepubsub_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XHeartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XHeartbeat) ProtoMessage() {}

func (x *XHeartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_cachepubsub_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XHeartbeat.ProtoReflect.Descriptor instead.
func (*XHeartbeat) Descriptor() ([]byte, []int) {
	return file_cachepubsub_proto_rawDescGZIP(), []int{6}
}

var File_cachepubsub_proto protoreflect.FileDescriptor

var file_cachepubsub_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x61, 0x63, 0x68, 0x65, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x13, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x0f, 0x5f, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69,
	0x63, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x5f, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x04, 0x80, 0xb5, 0x18, 0x00, 0x22,
	0xbc, 0x01, 0x0a, 0x14, 0x5f, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61,
	0x63, 0x68, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x44, 0x0a,
	0x1f, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x61, 0x74, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x1b, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x41, 0x74,
	0x54, 0x6f, 0x70, 0x69, 0x63, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x73, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x50, 0x61, 0x67, 0x65, 0x3a, 0x04, 0x80, 0xb5, 0x18, 0x01, 0x22, 0xe0,
	0x01, 0x0a, 0x11, 0x5f, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x35, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x5f, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x49,
	0x74, 0x65, 0x6d, 0x48, 0x00, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x12, 0x4b, 0x0a, 0x0d, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x5f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e,
	0x74, 0x69, 0x6e, 0x75, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x6e, 0x74, 0x69, 0x6e, 0x75, 0x69, 0x74, 0x79, 0x12, 0x3f, 0x0a, 0x09, 0x68, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x61,
	0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75,
	0x62, 0x2e, 0x5f, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x48, 0x00, 0x52, 0x09,
	0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x22, 0xc0, 0x01, 0x0a, 0x0a, 0x5f, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x32, 0x0a, 0x15, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x13, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x5f, 0x54, 0x6f, 0x70, 0x69, 0x63,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
	0x50, 0x61, 0x67, 0x65, 0x22, 0x45, 0x0a, 0x0b, 0x5f, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x06, 0x62, 0x69, 0x6e,
	0x61, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x06, 0x62, 0x69, 0x6e,
	0x61, 0x72, 0x79, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x9a, 0x01, 0x0a, 0x0e,
	0x5f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x69, 0x74, 0x79, 0x12, 0x2e,
	0x0a, 0x13, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x11, 0x6c, 0x61, 0x73,
	0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x2c,
	0x0a, 0x12, 0x6e, 0x65, 0x77, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x73, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x6e, 0x65, 0x77, 0x54,
	0x6f, 0x70, 0x69, 0x63, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x11,
	0x6e, 0x65, 0x77, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x6e, 0x65, 0x77, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x50, 0x61, 0x67, 0x65, 0x22, 0x0c, 0x0a, 0x0a, 0x5f, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x32, 0xab, 0x01, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x73, 0x75,
	0x62, 0x12, 0x3f, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12, 0x24, 0x2e, 0x63,
	0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x75, 0x62, 0x73,
	0x75, 0x62, 0x2e, 0x5f, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x5f, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x60, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12,
	0x29, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x5f, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62,
	0x2e, 0x5f, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x74,
	0x65, 0x6d, 0x30, 0x01, 0x42, 0x72, 0x0a, 0x18, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62,
	0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d,
	0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x68, 0x71, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d,
	0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x3b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x64,
	0x6b, 0x5f, 0x67, 0x6f, 0xaa, 0x02, 0x21, 0x4d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x6f, 0x2e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x61, 0x63, 0x68, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x50, 0x75, 0x62, 0x73, 0x75, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cachepubsub_proto_rawDescOnce sync.Once
	file_cachepubsub_proto_rawDescData = file_cachepubsub_proto_rawDesc
)

func file_cachepubsub_proto_rawDescGZIP() []byte {
	file_cachepubsub_proto_rawDescOnce.Do(func() {
		file_cachepubsub_proto_rawDescData = protoimpl.X.CompressGZIP(file_cachepubsub_proto_rawDescData)
	})
	return file_cachepubsub_proto_rawDescData
}

var file_cachepubsub_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_cachepubsub_proto_goTypes = []any{
	(*XPublishRequest)(nil),      // 0: cache_client.pubsub._PublishRequest
	(*XSubscriptionRequest)(nil), // 1: cache_client.pubsub._SubscriptionRequest
	(*XSubscriptionItem)(nil),    // 2: cache_client.pubsub._SubscriptionItem
	(*XTopicItem)(nil),           // 3: cache_client.pubsub._TopicItem
	(*XTopicValue)(nil),          // 4: cache_client.pubsub._TopicValue
	(*XDiscontinuity)(nil),       // 5: cache_client.pubsub._Discontinuity
	(*XHeartbeat)(nil),           // 6: cache_client.pubsub._Heartbeat
	(*XEmpty)(nil),               // 7: common._Empty
}
var file_cachepubsub_proto_depIdxs = []int32{
	4, // 0: cache_client.pubsub._PublishRequest.value:type_name -> cache_client.pubsub._TopicValue
	3, // 1: cache_client.pubsub._SubscriptionItem.item:type_name -> cache_client.pubsub._TopicItem
	5, // 2: cache_client.pubsub._SubscriptionItem.discontinuity:type_name -> cache_client.pubsub._Discontinuity
	6, // 3: cache_client.pubsub._SubscriptionItem.heartbeat:type_name -> cache_client.pubsub._Heartbeat
	4, // 4: cache_client.pubsub._TopicItem.value:type_name -> cache_client.pubsub._TopicValue
	0, // 5: cache_client.pubsub.Pubsub.Publish:input_type -> cache_client.pubsub._PublishRequest
	1, // 6: cache_client.pubsub.Pubsub.Subscribe:input_type -> cache_client.pubsub._SubscriptionRequest
	7, // 7: cache_client.pubsub.Pubsub.Publish:output_type -> common._Empty
	2, // 8: cache_client.pubsub.Pubsub.Subscribe:output_type -> cache_client.pubsub._SubscriptionItem
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_cachepubsub_proto_init() }
func file_cachepubsub_proto_init() {
	if File_cachepubsub_proto != nil {
		return
	}
	file_common_proto_init()
	file_extensions_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_cachepubsub_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*XPublishRequest); i {
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
		file_cachepubsub_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*XSubscriptionRequest); i {
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
		file_cachepubsub_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*XSubscriptionItem); i {
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
		file_cachepubsub_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*XTopicItem); i {
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
		file_cachepubsub_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*XTopicValue); i {
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
		file_cachepubsub_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*XDiscontinuity); i {
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
		file_cachepubsub_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*XHeartbeat); i {
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
	file_cachepubsub_proto_msgTypes[2].OneofWrappers = []any{
		(*XSubscriptionItem_Item)(nil),
		(*XSubscriptionItem_Discontinuity)(nil),
		(*XSubscriptionItem_Heartbeat)(nil),
	}
	file_cachepubsub_proto_msgTypes[4].OneofWrappers = []any{
		(*XTopicValue_Text)(nil),
		(*XTopicValue_Binary)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cachepubsub_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cachepubsub_proto_goTypes,
		DependencyIndexes: file_cachepubsub_proto_depIdxs,
		MessageInfos:      file_cachepubsub_proto_msgTypes,
	}.Build()
	File_cachepubsub_proto = out.File
	file_cachepubsub_proto_rawDesc = nil
	file_cachepubsub_proto_goTypes = nil
	file_cachepubsub_proto_depIdxs = nil
}
