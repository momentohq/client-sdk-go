package config

import "time"

type GrpcConfigurationProps struct {
	// number of milliseconds the client is willing to wait for an RPC to complete before it is terminated
	// with a DeadlineExceeded error.
	deadline time.Duration
}

// GrpcConfiguration Encapsulates gRPC configuration tunables.
type GrpcConfiguration interface {
	// GetDeadline Returns number of milliseconds the client is willing to wait for an RPC to complete before
	//it is terminated with a DeadlineExceeded error.
	GetDeadline() time.Duration

	// WithDeadline Copy constructor for overriding the client-side deadline. Returns a new GrpcConfiguration
	//with the specified client-side deadline
	WithDeadline(deadline time.Duration) GrpcConfiguration
}
