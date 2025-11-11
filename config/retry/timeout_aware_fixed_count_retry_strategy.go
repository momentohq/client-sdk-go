package retry

import (
	"fmt"
	"strconv"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

const (
	// DefaultRetryTimeoutDuration is the duration the client is willing to wait for
	// a response on a retry attempt. Defaults to 500 milliseconds.
	DefaultRetryTimeoutDuration = 500 * time.Millisecond
)

type timeoutAwareFixedCountRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	maxAttempts         int
	log                 logger.MomentoLogger
	timeoutDuration     time.Duration
}

// TimeoutAwareFixedCountRetryStrategy is a retry strategy that treats DeadlineExceeded as retryable.
// It is used in conjunction with the TimeoutAwareEligibilityStrategy.
//
// Typically, a DeadlineExceeded would indicate that the overall client timeout has been reached and no further
// retries should be attempted. However, there are some cases where timeouts may occur due to the client being
// overloaded or experiencing transient network issues. In these cases, it may be beneficial to retry the request
// even if the overall timeout has been reached.
//
// Note that this is different than using the FixedTimeoutRetryStrategy, which does not attempt a retry with a
// shorter timeout until after the first request fails due to a non-timeout error, and will retry only until the
// client's overall timeout is reached.
type TimeoutAwareFixedCountRetryStrategy interface {
	Strategy
	WithMaxAttempts(attempts int) Strategy
	WithEligibilityStrategy(s EligibilityStrategy) Strategy
	WithTimeoutDuration(duration time.Duration) Strategy
}

type TimeoutAwareFixedCountRetryStrategyProps struct {
	LoggerFactory       logger.MomentoLoggerFactory
	MaxAttempts         int
	EligibilityStrategy EligibilityStrategy
	TimeoutDuration     time.Duration
}

func NewTimeoutAwareFixedCountRetryStrategy(props TimeoutAwareFixedCountRetryStrategyProps) Strategy {
	eligibilityStrategy := EligibilityStrategy(TimeoutAwareEligibilityStrategy{})
	if props.EligibilityStrategy != nil {
		eligibilityStrategy = props.EligibilityStrategy
	}
	maxAttempts := 3
	if props.MaxAttempts != 0 {
		maxAttempts = props.MaxAttempts
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("timeout-aware-fixed-count-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("timeout-aware-fixed-count-retry-strategy")
	}

	var timeoutDuration time.Duration
	if props.TimeoutDuration != 0 {
		timeoutDuration = props.TimeoutDuration
	} else {
		timeoutDuration = DefaultRetryTimeoutDuration
	}
	fmt.Printf("\nUsing timeout duration: %v for timeout-aware fixed count retry strategy\n", timeoutDuration)

	return &timeoutAwareFixedCountRetryStrategy{
		eligibilityStrategy: eligibilityStrategy,
		maxAttempts:         maxAttempts,
		log:                 log,
		timeoutDuration:     timeoutDuration,
	}
}

func (r *timeoutAwareFixedCountRetryStrategy) WithMaxAttempts(attempts int) Strategy {
	return &timeoutAwareFixedCountRetryStrategy{
		log:                 r.log,
		maxAttempts:         attempts,
		eligibilityStrategy: r.eligibilityStrategy,
		timeoutDuration:     r.timeoutDuration,
	}
}

func (r *timeoutAwareFixedCountRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	return &timeoutAwareFixedCountRetryStrategy{
		log:                 r.log,
		maxAttempts:         r.maxAttempts,
		eligibilityStrategy: s,
		timeoutDuration:     r.timeoutDuration,
	}
}

func (r *timeoutAwareFixedCountRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	fmt.Printf("\noverall deadline in DetermineWhenToRetry: %v\n", props.OverallDeadline.Local())
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
		"\nDetermined request is retryable; retrying now: [method: %s, status: %s, attempt_count: %s, max_attempts: %s]",
		props.GrpcMethod,
		props.GrpcStatusCode.String(),
		strconv.Itoa(props.AttemptNumber),
		strconv.Itoa(r.maxAttempts),
	)
	timeTilNextRetry := 0
	return &timeTilNextRetry
}

func (r *timeoutAwareFixedCountRetryStrategy) String() string {
	return fmt.Sprintf(
		"timeoutAwareFixedCountRetryStrategy{eligibilityStrategy=%v, maxAttempts=%v, log=%v, timeoutDuration=%v}",
		r.eligibilityStrategy,
		r.maxAttempts,
		r.log,
		r.timeoutDuration,
	)
}

// CalculateRetryDeadline calculates the new deadline for a retry attempt if using a TimeoutAwareEligibilityStrategy.
func (r *timeoutAwareFixedCountRetryStrategy) CalculateNewOverallDeadline() time.Time {
	deadline := time.Now().Add(r.timeoutDuration)
	r.log.Debug("Added %v to current time, new deadline is: %v", r.timeoutDuration, deadline.Local())
	return deadline
}
