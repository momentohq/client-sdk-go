package retry

import (
	"fmt"
	"github.com/momentohq/client-sdk-go/config/logger"
	"strconv"
)

type FixedCountRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	maxAttempts         int
	log                 logger.MomentoLogger
}

func NewFixedCountRetryStrategy(logFactory logger.MomentoLoggerFactory) Strategy {
	return FixedCountRetryStrategy{
		eligibilityStrategy: DefaultEligibilityStrategy{},
		maxAttempts:         3,
		log:                 logFactory.GetLogger("fixed-count-retry-strategy"),
	}
}

func (r FixedCountRetryStrategy) WithMaxAttempts(attempts int) Strategy {
	r.maxAttempts = attempts
	return r
}

func (r FixedCountRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	r.eligibilityStrategy = s
	return r
}

func (r FixedCountRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
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

func (r FixedCountRetryStrategy) String() string {
	return fmt.Sprintf(
		"FixedCountRetryStrategy{eligibilityStrategy=%v, maxAttempts=%v, log=%v}",
		r.eligibilityStrategy,
		r.maxAttempts,
		r.log)
}

