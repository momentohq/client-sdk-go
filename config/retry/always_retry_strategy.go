package retry

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type alwaysRetryStrategy struct {
	log     logger.MomentoLogger
	retryMs *int
}

type AlwaysRetryStrategy interface {
	Strategy
	WithRetryMs(ms int) Strategy
}

type AlwaysRetryStrategyProps struct {
	LoggerFactory logger.MomentoLoggerFactory
	RetryMs       *int
}

func (r *alwaysRetryStrategy) WithRetryMs(ms int) Strategy {
	return &alwaysRetryStrategy{
		log:     r.log,
		retryMs: &ms,
	}
}

func (r *alwaysRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	r.log.Debug(
		"Always retry strategy returning %d ms for [method: %s, status: %s]",
		*r.retryMs,
		props.GrpcMethod,
		props.GrpcStatusCode.String(),
	)
	return r.retryMs
}

func (r *alwaysRetryStrategy) String() string {
	return fmt.Sprintf(
		"alwaysRetryStrategy{retryMs: %d, log: %s}\n",
		*r.retryMs,
		r.log,
	)
}

// NewAlwaysRetryStrategy is a retry strategy that always retries any request after a fixed delay.
// It is intended to maintain compatibility with the legacy behavior of Momento Topic subscriptions,
// which are now able to specify a retry strategy for determining when to reconnect to a subscription
// that has been interrupted. Switching to any of the other available retry strategies is recommended
// but not required and may require additional error handling after `Item()` or `Event()` calls as
// errors that previously resulted in a retry will now be returned to the caller.
// Deprecated: This strategy is deprecated and will be removed in a future release.
func NewAlwaysRetryStrategy(props AlwaysRetryStrategyProps) Strategy {
	retryMsInt := 500
	retryMs := &retryMsInt
	if props.RetryMs != nil {
		retryMs = props.RetryMs
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("always-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("always-retry-strategy")
	}
	return &alwaysRetryStrategy{
		log:     log,
		retryMs: retryMs,
	}
}
