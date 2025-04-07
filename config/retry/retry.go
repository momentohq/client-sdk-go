package retry

import (
	"time"

	"google.golang.org/grpc/codes"
)

type StrategyProps struct {
	GrpcStatusCode  codes.Code
	GrpcMethod      string
	AttemptNumber   int
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

	// GetResponseDataReceivedTimeoutMillis returns the timeout for a retry attemptin milliseconds.
	GetResponseDataReceivedTimeoutMillis() *int
}
