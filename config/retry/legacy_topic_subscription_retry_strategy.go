package retry

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type legacyTopicSubscriptionRetryStrategy struct {
	log     logger.MomentoLogger
	retryMs *int
}

type LegacyTopicSubscriptionRetryStrategy interface {
	Strategy
	WithRetryMs(ms int) Strategy
}

type LegacyTopicSubscriptionRetryStrategyProps struct {
	LoggerFactory logger.MomentoLoggerFactory
	RetryMs       *int
}

func (r *legacyTopicSubscriptionRetryStrategy) WithRetryMs(ms int) Strategy {
	return &legacyTopicSubscriptionRetryStrategy{
		log:     r.log,
		retryMs: &ms,
	}
}

func (r *legacyTopicSubscriptionRetryStrategy) GetResponseDataReceivedTimeoutMillis() *int {
	return nil
}

func (r *legacyTopicSubscriptionRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	r.log.Debug(
		"Always retry strategy returning %d ms for [method: %s, status: %s]",
		*r.retryMs,
		props.GrpcMethod,
		props.GrpcStatusCode.String(),
	)
	return r.retryMs
}

func (r *legacyTopicSubscriptionRetryStrategy) String() string {
	return fmt.Sprintf(
		"legacyTopicSubscriptionRetryStrategy{retryMs: %d, log: %s}\n",
		*r.retryMs,
		r.log,
	)
}

// NewLegacyTopicSubscriptionRetryStrategy returns a strategy that always retries any request after a fixed delay.
// It is intended to maintain compatibility with the legacy behavior of Momento Topic subscriptions,
// which are now able to specify a retry strategy for determining when to reconnect to a subscription
// that has been interrupted. Switching to any of the other available retry strategies is recommended
// but not required and may require additional error handling after `Item()` or `Event()` calls as
// errors that previously resulted in a retry will now be returned to the caller.
// Deprecated: This strategy is deprecated and will be removed in a future release.
func NewLegacyTopicSubscriptionRetryStrategy(props LegacyTopicSubscriptionRetryStrategyProps) Strategy {
	retryMsInt := 500
	retryMs := &retryMsInt
	if props.RetryMs != nil {
		retryMs = props.RetryMs
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("legacy-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("legacy-retry-strategy")
	}
	return &legacyTopicSubscriptionRetryStrategy{
		log:     log,
		retryMs: retryMs,
	}
}
