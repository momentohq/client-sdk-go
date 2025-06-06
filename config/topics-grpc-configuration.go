package config

import "time"

// The maximum number of concurrent streams that can be created on a single gRPC channel.
const MAX_CONCURRENT_STREAMS_PER_CHANNEL int = 100

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
