package config

import (
	"fmt"
	"time"
)

type TopicsGrpcConfigurationProps struct {
	// The number of milliseconds the client is willing to wait for an RPC to complete before it is terminated
	// with a DeadlineExceeded error.
	ClientTimeout time.Duration

	// The maximum message length the client can send to the server.  If the client attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	MaxSendMessageLength int

	// The maximum message length the client can receive from the server.  If the server attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	MaxReceiveMessageLength int

	// NumUnaryGrpcChannels represents the number of GRPC channels the topic client
	// should open and work with for unary operations (i.e. topic publishes).
	NumUnaryGrpcChannels uint32
}

type StaticTopicsGrpcConfigurationProps struct {
	TopicsGrpcConfigurationProps

	// NumStreamGrpcChannels represents the number of GRPC channels the topic client
	// should open and work with for stream operations (i.e. topic subscriptions).
	NumStreamGrpcChannels uint32
}

type DynamicTopicsGrpcConfigurationProps struct {
	TopicsGrpcConfigurationProps

	// MaxSubscriptions represents the maximum number of subscriptions the topic client
	// should open and work with for stream operations (i.e. topic subscriptions).
	MaxSubscriptions uint32
}

// GrpcConfiguration Encapsulates gRPC configuration tunables.
type TopicsGrpcConfiguration interface {
	// GetClientTimeout Returns number of milliseconds the client is willing to wait for an RPC to complete before
	// it is terminated with a DeadlineExceeded error.
	GetClientTimeout() time.Duration

	// GetKeepAlivePermitWithoutCalls returns bool indicating if it is permissible to send keepalive pings from the client without any outstanding calls.
	GetKeepAlivePermitWithoutCalls() bool

	// GetKeepAliveTimeout returns number of milliseconds the client will wait for a response from a keepalive or ping.
	GetKeepAliveTimeout() time.Duration

	// GetKeepAliveTime returns the interval at which to send the keepalive or ping.
	GetKeepAliveTime() time.Duration

	// GetMaxSendMessageLength is the maximum message length the client can send to the server.  If the client attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	GetMaxSendMessageLength() int

	// GetMaxReceiveMessageLength is the maximum message length the client can receive from the server.  If the server attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	GetMaxReceiveMessageLength() int

	// GetNumUnaryGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for unary operations (i.e. topic publishes).
	GetNumUnaryGrpcChannels() uint32
}

// static version

type TopicsStaticGrpcConfigurationType interface {
	TopicsGrpcConfiguration

	// GetNumStreamGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for stream operations (i.e. topic subscriptions).
	GetNumStreamGrpcChannels() uint32

	// WithNumStreamGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for stream operations. Each GRPC connection can multiplex 100 concurrent subscriptions.
	// Defaults to 4.
	WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsStaticGrpcConfigurationType

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsStaticGrpcConfigurationType

	// WithKeepAliveTime Copy constructor for overriding the keepalive time.
	// After a duration of this time the client/server pings its peer to see if the transport is still alive.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAliveTime(keepAliveTime time.Duration) TopicsStaticGrpcConfigurationType

	// WithKeepAliveDisabled disables grpc keepalives```
	// Returns a new GrpcConfiguration with keepalive settings disabled (they're enabled by default)
	WithKeepAliveDisabled() TopicsStaticGrpcConfigurationType

	// WithKeepAliveTimeout Copy constructor for overriding the keepalive timeout. After waiting for a duration of this time,
	// if the keepalive ping sender does not receive the ping ack, it will close the transport.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsStaticGrpcConfigurationType

	// WithKeepAlivePermitWithoutCalls Copy constructor for overriding the keepalive permit without calls.
	// Indicates if it permissible to send keepalive pings from the client without any outstanding streams.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsStaticGrpcConfigurationType

	// WithClientTimeout Copy constructor for overriding the client-side deadline. Returns a new GrpcConfiguration
	// with the specified client-side deadline
	WithClientTimeout(deadline time.Duration) TopicsStaticGrpcConfigurationType
}

type TopicsStaticGrpcConfiguration struct {
	clientTimeout               time.Duration
	keepAlivePermitWithoutCalls bool
	keepAliveTimeout            time.Duration
	keepAliveTime               time.Duration
	maxSendMessageLength        int
	maxReceiveMessageLength     int
	numStreamGrpcChannels       uint32
	numUnaryGrpcChannels        uint32
}

// NewStaticGrpcConfiguration constructs a new TopicsGrpcConfiguration to tune lower-level grpc settings.
// Note: keepalive settings are enabled by default, use WithKeepAliveDisabled() to disable all of them,
// or use the appropriate copy constructor to override individual settings.
func NewTopicsStaticGrpcConfiguration(grpcConfiguration *StaticTopicsGrpcConfigurationProps) *TopicsStaticGrpcConfiguration {
	// We set keepalive values to defaults because we can't tell if users set them to zero values
	// or if the settings were just omitted from the props struct, thus defaulting them to zero values.

	maxSendLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.MaxSendMessageLength > 0 {
		maxSendLength = grpcConfiguration.MaxSendMessageLength
	}

	maxReceiveLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.MaxReceiveMessageLength > 0 {
		maxReceiveLength = grpcConfiguration.MaxReceiveMessageLength
	}

	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               grpcConfiguration.ClientTimeout,
		keepAlivePermitWithoutCalls: DEFAULT_KEEPALIVE_WITHOUT_STREAM,
		keepAliveTimeout:            DEFAULT_KEEPALIVE_TIMEOUT,
		keepAliveTime:               DEFAULT_KEEPALIVE_TIME,
		maxSendMessageLength:        maxSendLength,
		maxReceiveMessageLength:     maxReceiveLength,
		numStreamGrpcChannels:       grpcConfiguration.NumStreamGrpcChannels,
		numUnaryGrpcChannels:        grpcConfiguration.NumUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) GetClientTimeout() time.Duration {
	return s.clientTimeout
}

func (s *TopicsStaticGrpcConfiguration) GetKeepAlivePermitWithoutCalls() bool {
	return s.keepAlivePermitWithoutCalls
}

func (s *TopicsStaticGrpcConfiguration) GetKeepAliveTimeout() time.Duration {
	return s.keepAliveTimeout
}

func (s *TopicsStaticGrpcConfiguration) GetKeepAliveTime() time.Duration {
	return s.keepAliveTime
}

func (s *TopicsStaticGrpcConfiguration) GetMaxSendMessageLength() int {
	return s.maxSendMessageLength
}

func (s *TopicsStaticGrpcConfiguration) GetMaxReceiveMessageLength() int {
	return s.maxReceiveMessageLength
}

func (s *TopicsStaticGrpcConfiguration) GetNumStreamGrpcChannels() uint32 {
	return s.numStreamGrpcChannels
}

func (s *TopicsStaticGrpcConfiguration) GetNumUnaryGrpcChannels() uint32 {
	return s.numUnaryGrpcChannels
}

func (s *TopicsStaticGrpcConfiguration) WithClientTimeout(clientTimeout time.Duration) TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveTime(keepAliveTime time.Duration) TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveDisabled() TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: false,
		keepAliveTimeout:            0,
		keepAliveTime:               0,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsStaticGrpcConfigurationType {
	return &TopicsStaticGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) String() string {
	return fmt.Sprintf(
		"TopicsGrpcConfiguration{clientTimeout=%v, keepAlivePermitWithoutCalls=%v, keepAliveTimeout=%v, keepAliveTime=%v, maxSendMessageLength=%v, maxReceiveMessageLength=%v, numStreamGrpcChannels=%v, numUnaryGrpcChannels=%v}",
		s.clientTimeout, s.keepAlivePermitWithoutCalls, s.keepAliveTimeout, s.keepAliveTime, s.maxSendMessageLength, s.maxReceiveMessageLength, s.numStreamGrpcChannels, s.numUnaryGrpcChannels,
	)
}

// dynamic version

type TopicsDynamicGrpcConfigurationType interface {
	TopicsGrpcConfiguration

	// GetMaxSubscriptions Returns the configuration option for the maximum number of subscriptions the topic client
	// should open and work with for stream operations (i.e. topic subscriptions).
	GetMaxSubscriptions() uint32

	// WithMaxSubscriptions is currently implemented to create the specified number of GRPC connections
	// for stream operations. Each GRPC connection can multiplex 100 concurrent subscriptions.
	// Defaults to 4.
	WithMaxSubscriptions(maxSubscriptions uint32) TopicsDynamicGrpcConfigurationType

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsDynamicGrpcConfigurationType

	// WithKeepAliveTime Copy constructor for overriding the keepalive time.
	// After a duration of this time the client/server pings its peer to see if the transport is still alive.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAliveTime(keepAliveTime time.Duration) TopicsDynamicGrpcConfigurationType

	// WithKeepAliveDisabled disables grpc keepalives```
	// Returns a new GrpcConfiguration with keepalive settings disabled (they're enabled by default)
	WithKeepAliveDisabled() TopicsDynamicGrpcConfigurationType

	// WithKeepAliveTimeout Copy constructor for overriding the keepalive timeout. After waiting for a duration of this time,
	// if the keepalive ping sender does not receive the ping ack, it will close the transport.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsDynamicGrpcConfigurationType

	// WithKeepAlivePermitWithoutCalls Copy constructor for overriding the keepalive permit without calls.
	// Indicates if it permissible to send keepalive pings from the client without any outstanding streams.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsDynamicGrpcConfigurationType

	// WithClientTimeout Copy constructor for overriding the client-side deadline. Returns a new GrpcConfiguration
	// with the specified client-side deadline
	WithClientTimeout(deadline time.Duration) TopicsDynamicGrpcConfigurationType
}

type TopicsDynamicGrpcConfiguration struct {
	clientTimeout               time.Duration
	keepAlivePermitWithoutCalls bool
	keepAliveTimeout            time.Duration
	keepAliveTime               time.Duration
	maxSendMessageLength        int
	maxReceiveMessageLength     int
	maxSubscriptions            uint32
	numUnaryGrpcChannels        uint32
}

// NewDynamicGrpcConfiguration constructs a new TopicsGrpcConfiguration to tune lower-level grpc settings
// for a static number of unary channels and a dynamic number of stream channels.
// Note: keepalive settings are enabled by default, use WithKeepAliveDisabled() to disable all of them,
// or use the appropriate copy constructor to override individual settings.
func NewTopicsDynamicGrpcConfiguration(grpcConfiguration *DynamicTopicsGrpcConfigurationProps) *TopicsDynamicGrpcConfiguration {
	// We set keepalive values to defaults because we can't tell if users set them to zero values
	// or if the settings were just omitted from the props struct, thus defaulting them to zero values.

	maxSendLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.MaxSendMessageLength > 0 {
		maxSendLength = grpcConfiguration.MaxSendMessageLength
	}

	maxReceiveLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.MaxReceiveMessageLength > 0 {
		maxReceiveLength = grpcConfiguration.MaxReceiveMessageLength
	}

	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               grpcConfiguration.ClientTimeout,
		keepAlivePermitWithoutCalls: DEFAULT_KEEPALIVE_WITHOUT_STREAM,
		keepAliveTimeout:            DEFAULT_KEEPALIVE_TIMEOUT,
		keepAliveTime:               DEFAULT_KEEPALIVE_TIME,
		maxSendMessageLength:        maxSendLength,
		maxReceiveMessageLength:     maxReceiveLength,
		maxSubscriptions:            grpcConfiguration.MaxSubscriptions,
		numUnaryGrpcChannels:        grpcConfiguration.NumUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) GetClientTimeout() time.Duration {
	return s.clientTimeout
}

func (s *TopicsDynamicGrpcConfiguration) GetKeepAlivePermitWithoutCalls() bool {
	return s.keepAlivePermitWithoutCalls
}

func (s *TopicsDynamicGrpcConfiguration) GetKeepAliveTimeout() time.Duration {
	return s.keepAliveTimeout
}

func (s *TopicsDynamicGrpcConfiguration) GetKeepAliveTime() time.Duration {
	return s.keepAliveTime
}

func (s *TopicsDynamicGrpcConfiguration) GetMaxSendMessageLength() int {
	return s.maxSendMessageLength
}

func (s *TopicsDynamicGrpcConfiguration) GetMaxReceiveMessageLength() int {
	return s.maxReceiveMessageLength
}

func (s *TopicsDynamicGrpcConfiguration) GetMaxSubscriptions() uint32 {
	return s.maxSubscriptions
}

func (s *TopicsDynamicGrpcConfiguration) GetNumUnaryGrpcChannels() uint32 {
	return s.numUnaryGrpcChannels
}

func (s *TopicsDynamicGrpcConfiguration) WithClientTimeout(clientTimeout time.Duration) TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            s.maxSubscriptions,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            s.maxSubscriptions,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            s.maxSubscriptions,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) WithKeepAliveTime(keepAliveTime time.Duration) TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            s.maxSubscriptions,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) WithKeepAliveDisabled() TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: false,
		keepAliveTimeout:            0,
		keepAliveTime:               0,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            s.maxSubscriptions,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) WithMaxSubscriptions(maxSubscriptions uint32) TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            maxSubscriptions,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsDynamicGrpcConfigurationType {
	return &TopicsDynamicGrpcConfiguration{
		clientTimeout:               s.clientTimeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		maxSubscriptions:            s.maxSubscriptions,
		numUnaryGrpcChannels:        numUnaryGrpcChannels,
	}
}

func (s *TopicsDynamicGrpcConfiguration) String() string {
	return fmt.Sprintf(
		"TopicsGrpcConfiguration{clientTimeout=%v, keepAlivePermitWithoutCalls=%v, keepAliveTimeout=%v, keepAliveTime=%v, maxSendMessageLength=%v, maxReceiveMessageLength=%v, maxSubscriptions=%v, numUnaryGrpcChannels=%v}",
		s.clientTimeout, s.keepAlivePermitWithoutCalls, s.keepAliveTimeout, s.keepAliveTime, s.maxSendMessageLength, s.maxReceiveMessageLength, s.maxSubscriptions, s.numUnaryGrpcChannels,
	)
}
