package retry

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"google.golang.org/grpc/codes"
)

const (
	DefaultResponseDataReceivedTimeoutMillis = 1000 // 1 second default timeout for retry attempts
	DefaultRetryDelayIntervalMillis          = 100  // Schedule retry attempt for 100ms later +/- jitter
)

type fixedTimeoutRetryStrategy struct {
	eligibilityStrategy               EligibilityStrategy
	log                               logger.MomentoLogger
	responseDataReceivedTimeoutMillis int
	retryDelayIntervalMillis          int
}

type FixedTimeoutRetryStrategy interface {
	Strategy
	WithResponseDataReceivedTimeoutMillis(timeout int) Strategy
	WithRetryDelayIntervalMillis(delay int) Strategy
	WithEligibilityStrategy(s EligibilityStrategy) Strategy
}

type FixedTimeoutRetryStrategyProps struct {
	EligibilityStrategy               EligibilityStrategy
	LoggerFactory                     logger.MomentoLoggerFactory
	ResponseDataReceivedTimeoutMillis int
	RetryDelayIntervalMillis          int
}

func NewFixedTimeoutRetryStrategy(props FixedTimeoutRetryStrategyProps) Strategy {
	eligibilityStrategy := EligibilityStrategy(DefaultEligibilityStrategy{})
	if props.EligibilityStrategy != nil {
		eligibilityStrategy = props.EligibilityStrategy
	}
	responseDataReceivedTimeoutMillis := DefaultResponseDataReceivedTimeoutMillis
	if props.ResponseDataReceivedTimeoutMillis != 0 {
		responseDataReceivedTimeoutMillis = props.ResponseDataReceivedTimeoutMillis
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
		eligibilityStrategy:               eligibilityStrategy,
		log:                               log,
		responseDataReceivedTimeoutMillis: responseDataReceivedTimeoutMillis,
		retryDelayIntervalMillis:          retryDelayIntervalMillis,
	}
}

func (r *fixedTimeoutRetryStrategy) WithResponseDataReceivedTimeoutMillis(timeout int) Strategy {
	return &fixedTimeoutRetryStrategy{
		log:                               r.log,
		eligibilityStrategy:               r.eligibilityStrategy,
		responseDataReceivedTimeoutMillis: timeout,
		retryDelayIntervalMillis:          r.retryDelayIntervalMillis,
	}
}

func (r *fixedTimeoutRetryStrategy) WithRetryDelayIntervalMillis(delay int) Strategy {
	return &fixedTimeoutRetryStrategy{
		log:                               r.log,
		eligibilityStrategy:               r.eligibilityStrategy,
		responseDataReceivedTimeoutMillis: r.responseDataReceivedTimeoutMillis,
		retryDelayIntervalMillis:          delay,
	}
}

func (r *fixedTimeoutRetryStrategy) WithEligibilityStrategy(strategy EligibilityStrategy) Strategy {
	return &fixedTimeoutRetryStrategy{
		log:                               r.log,
		eligibilityStrategy:               strategy,
		responseDataReceivedTimeoutMillis: r.responseDataReceivedTimeoutMillis,
		retryDelayIntervalMillis:          r.retryDelayIntervalMillis,
	}
}

func (r *fixedTimeoutRetryStrategy) GetResponseDataReceivedTimeoutMillis() *int {
	return &r.responseDataReceivedTimeoutMillis
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
		"fixedTimeoutRetryStrategy{eligibilityStrategy=%v, responseDataReceivedTimeoutMillis=%v, retryDelayIntervalMillis=%v}",
		r.eligibilityStrategy,
		r.responseDataReceivedTimeoutMillis,
		r.retryDelayIntervalMillis,
	)
}
