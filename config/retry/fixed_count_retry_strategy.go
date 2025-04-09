package retry

import (
	"fmt"
	"strconv"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type fixedCountRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	maxAttempts         int
	log                 logger.MomentoLogger
}

type FixedCountRetryStrategy interface {
	Strategy
	WithMaxAttempts(attempts int) Strategy
	WithEligibilityStrategy(s EligibilityStrategy) Strategy
}

type FixedCountRetryStrategyProps struct {
	LoggerFactory       logger.MomentoLoggerFactory
	MaxAttempts         int
	EligibilityStrategy EligibilityStrategy
}

func NewFixedCountRetryStrategy(props FixedCountRetryStrategyProps) Strategy {
	eligibilityStrategy := EligibilityStrategy(DefaultEligibilityStrategy{})
	if props.EligibilityStrategy != nil {
		eligibilityStrategy = props.EligibilityStrategy
	}
	maxAttempts := 3
	if props.MaxAttempts != 0 {
		maxAttempts = props.MaxAttempts
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("fixed-count-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("fixed-count-retry-strategy")
	}

	return &fixedCountRetryStrategy{
		eligibilityStrategy: eligibilityStrategy,
		maxAttempts:         maxAttempts,
		log:                 log,
	}
}

func (r *fixedCountRetryStrategy) WithMaxAttempts(attempts int) Strategy {
	return &fixedCountRetryStrategy{
		log:                 r.log,
		maxAttempts:         attempts,
		eligibilityStrategy: r.eligibilityStrategy,
	}
}

func (r *fixedCountRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	return &fixedCountRetryStrategy{
		log:                 r.log,
		maxAttempts:         r.maxAttempts,
		eligibilityStrategy: s,
	}
}

func (r *fixedCountRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	if !r.eligibilityStrategy.IsEligibleForRetry(props) {
		r.log.Debug(
			"Request is not retryable: [method: %s, status: %s]", props.GrpcMethod, props.GrpcStatusCode.String(),
		)
		return nil
	}

	if props.AttemptNumber > r.maxAttempts {
		r.log.Debug(
			"Exceeded max retry attempts; not retrying: [method: %s, status: %s, attempt_count: %s, max_attempts: %s]",
			props.GrpcMethod,
			props.GrpcStatusCode.String(),
			strconv.Itoa(props.AttemptNumber),
			strconv.Itoa(r.maxAttempts),
		)
		return nil
	}

	r.log.Debug(
		"Determined request is retryable; retrying now: [method: %s, status: %s, attempt_count: %s, max_attempts: %s]",
		props.GrpcMethod,
		props.GrpcStatusCode.String(),
		strconv.Itoa(props.AttemptNumber),
		strconv.Itoa(r.maxAttempts),
	)
	timeTilNextRetry := 0
	return &timeTilNextRetry
}

func (r *fixedCountRetryStrategy) String() string {
	return fmt.Sprintf(
		"fixedCountRetryStrategy{eligibilityStrategy=%v, maxAttempts=%v, log=%v}",
		r.eligibilityStrategy,
		r.maxAttempts,
		r.log,
	)
}

func (r *fixedCountRetryStrategy) CalculateRetryDeadline(overallDeadline time.Time) *time.Time {
	return nil
}
