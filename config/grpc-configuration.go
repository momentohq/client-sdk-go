package config

import "time"

type GrpcConfigurationProps struct {
	// number of milliseconds the client is willing to wait for an RPC to complete before it is terminated
	// with a DeadlineExceeded error.
	deadline time.Duration
	// the maximum amount of memory, in megabytes, that a session is allowed to consume.  Sessions that consume
	// more than this amount will return a ResourceExhausted error.
	maxSessionMemoryMb uint32
}

// GrpcConfiguration Encapsulates gRPC configuration tunables.
type GrpcConfiguration interface {
	// GetDeadline Returns number of milliseconds the client is willing to wait for an RPC to complete before
	//it is terminated with a DeadlineExceeded error.
	GetDeadline() time.Duration

	// WithDeadline Copy constructor for overriding the client-side deadline. Returns a new GrpcConfiguration
	//with the specified client-side deadline
	WithDeadline(deadline time.Duration) GrpcConfiguration

	// GetMaxSessionMemoryMb the maximum amount of memory, in megabytes, that a session is allowed to consume.
	//Sessions that consume more than this amount will return a ResourceExhausted error.
	GetMaxSessionMemoryMb() uint32

	// WithMaxSessionMb Copy constructor for overriding the max session memory. maxSessionMemoryMb is the desired maximum
	//amount of memory, in megabytes, to allow a client session to consume. Returns  a new GrpcConfiguration with the
	//specified maximum memory
	WithMaxSessionMb(maxSessionMemoryMb uint32) GrpcConfiguration
}
