package retry

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"google.golang.org/grpc/codes"
)

const (
	// DefaultRetryTimeoutMillis is the number of milliseconds the client is willing to wait for
	// a response on a retry attempt. Defaults to 1 second.
	DefaultRetryTimeoutMillis = 1000
	// DefaultRetryDelayIntervalMillis is the number of milliseconds to wait between retry attempts.
	// Defaults to 100ms +/- jitter.
	DefaultRetryDelayIntervalMillis = 100
)

type fixedTimeoutRetryStrategy struct {
	eligibilityStrategy      EligibilityStrategy
	log                      logger.MomentoLogger
	retryTimeoutMillis       int
	retryDelayIntervalMillis int
}

// FixedTimeoutRetryStrategy is a retry strategy that retries a request up until the client timeout
// is reached. After the initial request fails, the next retry will be scheduled for retryDelayIntervalMillis
// from the current time, and the retried request will timeout after retryTimeoutMillis if there is no response.
type FixedTimeoutRetryStrategy interface {
	Strategy
	WithRetryTimeoutMillis(timeout int) Strategy
	WithRetryDelayIntervalMillis(delay int) Strategy
	WithEligibilityStrategy(s EligibilityStrategy) Strategy
}

type FixedTimeoutRetryStrategyProps struct {
	EligibilityStrategy      EligibilityStrategy
	LoggerFactory            logger.MomentoLoggerFactory
	RetryTimeoutMillis       int
	RetryDelayIntervalMillis int
}

func NewFixedTimeoutRetryStrategy(props FixedTimeoutRetryStrategyProps) Strategy {
	eligibilityStrategy := EligibilityStrategy(DefaultEligibilityStrategy{})
	if props.EligibilityStrategy != nil {
		eligibilityStrategy = props.EligibilityStrategy
	}
	retryTimeoutMillis := DefaultRetryTimeoutMillis
	if props.RetryTimeoutMillis != 0 {
		retryTimeoutMillis = props.RetryTimeoutMillis
	}
	retryDelayIntervalMillis := DefaultRetryDelayIntervalMillis
	if props.RetryDelayIntervalMillis != 0 {
		retryDelayIntervalMillis = props.RetryDelayIntervalMillis
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("fixed-timeout-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("fixed-timeout-retry-strategy")
	}

	return &fixedTimeoutRetryStrategy{
		eligibilityStrategy:      eligibilityStrategy,
		log:                      log,
		retryTimeoutMillis:       retryTimeoutMillis,
		retryDelayIntervalMillis: retryDelayIntervalMillis,
	}
}

func (r *fixedTimeoutRetryStrategy) WithRetryTimeoutMillis(timeout int) Strategy {
	return &fixedTimeoutRetryStrategy{
		log:                      r.log,
		eligibilityStrategy:      r.eligibilityStrategy,
		retryTimeoutMillis:       timeout,
		retryDelayIntervalMillis: r.retryDelayIntervalMillis,
	}
}

func (r *fixedTimeoutRetryStrategy) WithRetryDelayIntervalMillis(delay int) Strategy {
	return &fixedTimeoutRetryStrategy{
		log:                      r.log,
		eligibilityStrategy:      r.eligibilityStrategy,
		retryTimeoutMillis:       r.retryTimeoutMillis,
		retryDelayIntervalMillis: delay,
	}
}

func (r *fixedTimeoutRetryStrategy) WithEligibilityStrategy(strategy EligibilityStrategy) Strategy {
	return &fixedTimeoutRetryStrategy{
		log:                      r.log,
		eligibilityStrategy:      strategy,
		retryTimeoutMillis:       r.retryTimeoutMillis,
		retryDelayIntervalMillis: r.retryDelayIntervalMillis,
	}
}

// CalculateRetryDeadline calculates the deadline for a retry attempt using the retry timeout,
// but clips it to the overall deadline if the overall deadline is sooner.
func (r *fixedTimeoutRetryStrategy) CalculateRetryDeadline(overallDeadline time.Time) *time.Time {
	deadlineOffset := time.Duration(r.retryTimeoutMillis) * time.Millisecond
	deadline := time.Now().Add(deadlineOffset)
	if deadline.After(overallDeadline) {
		deadline = overallDeadline
	}
	return &deadline
}

func (r *fixedTimeoutRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	r.log.Debug(
		"Determining whether request is eligible for retry; status code: %s, "+
			"request type: %s, attemptNumber: %d",
		props.GrpcStatusCode, props.GrpcMethod, props.AttemptNumber,
	)

	// If a retry attempt's timeout has passed but the client's overall timeout has not yet passed,
	// we should reset the deadline and retry.
	if props.AttemptNumber > 0 && props.GrpcStatusCode == codes.DeadlineExceeded && props.OverallDeadline.After(time.Now()) {
		timeoutWithJitter := addJitter(r.retryDelayIntervalMillis)
		r.log.Debug(
			"Determined request is retryable; retrying after %d ms: [method: %s, status: %s, attempt: %d]",
			timeoutWithJitter,
			props.GrpcMethod,
			props.GrpcStatusCode.String(),
			props.AttemptNumber,
		)
		return &timeoutWithJitter
	}

	if !r.eligibilityStrategy.IsEligibleForRetry(props) {
		r.log.Debug(
			"Request is not retryable: [method: %s, status: %s]", props.GrpcMethod, props.GrpcStatusCode.String(),
		)
		return nil
	}

	timeoutWithJitter := addJitter(r.retryDelayIntervalMillis)

	r.log.Debug(
		"Determined request is retryable; retrying after %d ms: [method: %s, status: %s, attempt: %d]",
		timeoutWithJitter,
		props.GrpcMethod,
		props.GrpcStatusCode.String(),
		props.AttemptNumber,
	)
	return &timeoutWithJitter
}

func addJitter(whenToRetry int) int {
	return int((0.2*rand.Float64() + 0.9) * float64(whenToRetry))
}

func (r *fixedTimeoutRetryStrategy) String() string {
	return fmt.Sprintf(
		"fixedTimeoutRetryStrategy{eligibilityStrategy=%v, retryTimeoutMillis=%v, retryDelayIntervalMillis=%v}",
		r.eligibilityStrategy,
		r.retryTimeoutMillis,
		r.retryDelayIntervalMillis,
	)
}
