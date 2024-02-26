package config

import "time"

type GrpcConfigurationProps struct {
	// number of milliseconds the client is willing to wait for an RPC to complete before it is terminated
	// with a DeadlineExceeded error.
	deadline time.Duration

	// Indicates if it permissible to send keepalive pings from the client without any outstanding streams.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	keepAlivePermitWithoutCalls bool

	// After waiting for a duration of this time, if the keepalive ping sender does not receive the ping ack,
	// it will close the transport.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	keepAliveTimeout time.Duration

	// After a duration of this time the client/server pings its peer to see if the transport is still alive.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	keepAliveTime time.Duration

	// The maximum message length the client can send to the server.  If the client attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	maxSendMessageLength int

	// The maximum message length the client can receive from the server.  If the server attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	maxReceiveMessageLength int
}

// GrpcConfiguration Encapsulates gRPC configuration tunables.
type GrpcConfiguration interface {
	// GetDeadline Returns number of milliseconds the client is willing to wait for an RPC to complete before
	// it is terminated with a DeadlineExceeded error.
	GetDeadline() time.Duration

	// WithDeadline Copy constructor for overriding the client-side deadline. Returns a new GrpcConfiguration
	// with the specified client-side deadline
	WithDeadline(deadline time.Duration) GrpcConfiguration

	// Returns bool indicating if it is permissible to send keepalive pings from the client without any outstanding calls.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	GetKeepAlivePermitWithoutCalls() *bool

	// Returns number of milliseconds the client will wait for a response from a keepalive or ping.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	GetKeepAliveTimeout() *time.Duration

	// Returns the interval at which to send the keepalive or ping.
	//
	// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
	// when the connection is idle. However, they are very problematic for lambda environments where the lambda
	// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
	// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
	// Therefore, keep-alives should be disabled in lambda and similar environments.
	GetKeepAliveTime() *time.Duration

	// The maximum message length the client can send to the server.  If the client attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	GetMaxSendMessageLength() *int

	// The maximum message length the client can receive from the server.  If the server attempts to send a message
	// larger than this size, it will result in a RESOURCE_EXHAUSTED error.
	GetMaxReceiveMessageLength() *int
}
