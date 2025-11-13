package retry

import (
	"time"

	"google.golang.org/grpc/codes"
)

type StrategyProps struct {
	GrpcStatusCode codes.Code
	GrpcMethod     string
	AttemptNumber  int

	// OverallDeadline is the overall deadline for the request. It is currently only used
	// by the FixedTimeoutRetryStrategy to determine when to stop retrying.
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

type DeadlineAwareRetryStrategy interface {
	Strategy

	// CalculateRetryDeadline calculates the deadline for a retry attempt.
	// Returns nil if there is no adjustment to the deadline.
	CalculateRetryDeadline(overallDeadline time.Time) time.Time
}

type OverrideDeadlineRetryStrategy interface {
	Strategy

	// CalculateNewOverallDeadline calculates the deadline for a retry attempt and uses it to override the overall deadline.
	// Returns nil if there is no adjustment to the deadline.
	CalculateNewOverallDeadline() time.Time
}
