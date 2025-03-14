package retry

import (
	"google.golang.org/grpc/codes"
)

type StrategyProps struct {
	GrpcStatusCode codes.Code
	GrpcMethod     string
	AttemptNumber  int
}
type Strategy interface {
	// DetermineWhenToRetry Determines whether a grpc call can be retried and how long to wait before that retry.
	//
	// StrategyProps - Information about the grpc call, its last invocation, and how many times the call
	// has been made.
	//
	// Returns The time in seconds before the next retry should occur or nil if no retry should be attempted.
	DetermineWhenToRetry(props StrategyProps) *int
}
