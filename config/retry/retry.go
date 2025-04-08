package retry

import (
	"google.golang.org/grpc/codes"
	"time"
)

type StrategyProps struct {
	GrpcStatusCode  codes.Code
	GrpcMethod      string
	AttemptNumber   int
	// OverallDeadline is the overall deadline for the request. It is only used by FixedTimeoutStrategy types
	// and is currently ignored by all other strategy types.
	OverallDeadline time.Time
}

type Strategy interface {
	// DetermineWhenToRetry Determines whether a grpc call can be retried and how long to wait before that retry.
	//
	// StrategyProps - Information about the grpc call, its last invocation, and how many times the call
	// has been made.
	//
	// Returns The time in milliseconds before the next retry should occur or nil if no retry should be attempted.
	DetermineWhenToRetry(props StrategyProps) *int
}

type FixedTimeoutStrategy interface {
	Strategy
	// GetResponseDataReceivedTimeoutMillis returns the timeout for a retry attemptin milliseconds.
	GetResponseDataReceivedTimeoutMillis() *int
}
