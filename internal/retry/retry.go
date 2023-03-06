package retry

import (
	"strconv"

	"github.com/momentohq/client-sdk-go/config/logger"
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
	// RetryableProps - Information about the grpc call, its last invocation, and how many times the call
	//            has been made.
	//
	// Returns The time in seconds before the next retry should occur or nil if no retry should be attempted.
	DetermineWhenToRetry(props StrategyProps) *int
}

type fixedCountRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	maxAttempts         int
	log                 logger.MomentoLogger
}

func NewFixedCountRetryStrategy(logFactory logger.MomentoLoggerFactory) Strategy {
	return fixedCountRetryStrategy{
		eligibilityStrategy: DefaultEligibilityStrategy{},
		maxAttempts:         3,
		log:                 logFactory.GetLogger("fixed-count-retry-strategy", logger.DEBUG),
	}
}

func (r fixedCountRetryStrategy) WithMaxAttempts(attempts int) Strategy {
	r.maxAttempts = attempts
	return r
}

func (r fixedCountRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	r.eligibilityStrategy = s
	return r
}

func (r fixedCountRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	if !r.eligibilityStrategy.IsEligibleForRetry(props) {
		r.log.Debug(
			"Request is not retryable",
			"method", props.GrpcMethod,
			"status", props.GrpcStatusCode.String(),
		)
		return nil
	}

	if props.AttemptNumber > r.maxAttempts {
		r.log.Debug(
			"Exceeded max retry attempts not retrying",
			"method", props.GrpcMethod,
			"status", props.GrpcStatusCode.String(),
			"attempt_count", strconv.Itoa(props.AttemptNumber),
			"max_attempts", strconv.Itoa(r.maxAttempts),
		)
		return nil
	}

	r.log.Debug(
		"Determined request is retryable retrying now",
		"method", props.GrpcMethod,
		"status", props.GrpcStatusCode.String(),
		"attempt_count", strconv.Itoa(props.AttemptNumber),
		"max_attempts", strconv.Itoa(r.maxAttempts),
	)
	timeTilNextRetry := 0
	return &timeTilNextRetry
}
