package config

import (
	"fmt"
	"time"
)

type TopicsGrpcConfigurationProps struct {
	// The number of milliseconds the client is willing to wait for an RPC to complete before it is terminated
	// with a DeadlineExceeded error.
	client_timeout time.Duration

	// The maximum message length the client can send to the server.  If the client attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	maxSendMessageLength int

	// The maximum message length the client can receive from the server.  If the server attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	maxReceiveMessageLength int

	// NumStreamGrpcChannels represents the number of GRPC channels the topic client
	// should open and work with for stream operations (i.e. topic subscriptions).
	numStreamGrpcChannels uint32

	// NumUnaryGrpcChannels represents the number of GRPC channels the topic client
	// should open and work with for unary operations (i.e. topic publishes).
	numUnaryGrpcChannels uint32
}

// GrpcConfiguration Encapsulates gRPC configuration tunables.
type TopicsGrpcConfiguration interface {
	// GetClientTimeout Returns number of milliseconds the client is willing to wait for an RPC to complete before
	// it is terminated with a DeadlineExceeded error.
	GetClientTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding the client-side deadline. Returns a new GrpcConfiguration
	// with the specified client-side deadline
	WithClientTimeout(deadline time.Duration) TopicsGrpcConfiguration

	// GetKeepAlivePermitWithoutCalls returns bool indicating if it is permissible to send keepalive pings from the client without any outstanding calls.
	GetKeepAlivePermitWithoutCalls() bool

	// WithKeepAlivePermitWithoutCalls Copy constructor for overriding the keepalive permit without calls.
	// Indicates if it permissible to send keepalive pings from the client without any outstanding streams.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsGrpcConfiguration

	// GetKeepAliveTimeout returns number of milliseconds the client will wait for a response from a keepalive or ping.
	GetKeepAliveTimeout() time.Duration

	// WithKeepAliveTimeout Copy constructor for overriding the keepalive timeout. After waiting for a duration of this time,
	// if the keepalive ping sender does not receive the ping ack, it will close the transport.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsGrpcConfiguration

	// GetKeepAliveTime returns the interval at which to send the keepalive or ping.
	GetKeepAliveTime() time.Duration

	// WithKeepAliveTime Copy constructor for overriding the keepalive time.
	// After a duration of this time the client/server pings its peer to see if the transport is still alive.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	WithKeepAliveTime(keepAliveTime time.Duration) TopicsGrpcConfiguration

	// WithKeepAliveDisabled disables grpc keepalives```
	// Returns a new GrpcConfiguration with keepalive settings disabled (they're enabled by default)
	WithKeepAliveDisabled() TopicsGrpcConfiguration

	// GetMaxSendMessageLength is the maximum message length the client can send to the server.  If the client attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	GetMaxSendMessageLength() int

	// GetMaxReceiveMessageLength is the maximum message length the client can receive from the server.  If the server attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	GetMaxReceiveMessageLength() int

	// GetNumStreamGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for stream operations (i.e. topic subscriptions).
	GetNumStreamGrpcChannels() uint32

	// WithNumStreamGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for stream operations. Each GRPC connection can multiplex 100 concurrent subscriptions.
	// Defaults to 4.
	WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsGrpcConfiguration

	// GetNumUnaryGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for unary operations (i.e. topic publishes).
	GetNumUnaryGrpcChannels() uint32

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsGrpcConfiguration
}

type TopicsStaticGrpcConfiguration struct {
	client_timeout              time.Duration
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
func NewTopicsStaticGrpcConfiguration(grpcConfiguration *TopicsGrpcConfigurationProps) *TopicsStaticGrpcConfiguration {
	// We set keepalive values to defaults because we can't tell if users set them to zero values
	// or if the settings were just omitted from the props struct, thus defaulting them to zero values.

	maxSendLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.maxSendMessageLength > 0 {
		maxSendLength = grpcConfiguration.maxSendMessageLength
	}

	maxReceiveLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.maxReceiveMessageLength > 0 {
		maxReceiveLength = grpcConfiguration.maxReceiveMessageLength
	}

	return &TopicsStaticGrpcConfiguration{
		client_timeout:              grpcConfiguration.client_timeout,
		keepAlivePermitWithoutCalls: DEFAULT_KEEPALIVE_WITHOUT_STREAM,
		keepAliveTimeout:            DEFAULT_KEEPALIVE_TIMEOUT,
		keepAliveTime:               DEFAULT_KEEPALIVE_TIME,
		maxSendMessageLength:        maxSendLength,
		maxReceiveMessageLength:     maxReceiveLength,
		numStreamGrpcChannels:       grpcConfiguration.numStreamGrpcChannels,
		numUnaryGrpcChannels:        grpcConfiguration.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) GetClientTimeout() time.Duration {
	return s.client_timeout
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

func (s *TopicsStaticGrpcConfiguration) WithClientTimeout(client_timeout time.Duration) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveTime(keepAliveTime time.Duration) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveDisabled() TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: false,
		keepAliveTimeout:            0,
		keepAliveTime:               0,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
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
		"TopicsGrpcConfiguration{client_timeout=%v, keepAlivePermitWithoutCalls=%v, keepAliveTimeout=%v, keepAliveTime=%v, maxSendMessageLength=%v, maxReceiveMessageLength=%v, numStreamGrpcChannels=%v, numUnaryGrpcChannels=%v}",
		s.client_timeout, s.keepAlivePermitWithoutCalls, s.keepAliveTimeout, s.keepAliveTime, s.maxSendMessageLength, s.maxReceiveMessageLength, s.numStreamGrpcChannels, s.numUnaryGrpcChannels,
	)
}
