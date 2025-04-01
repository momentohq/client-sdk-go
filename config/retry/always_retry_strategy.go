package retry

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type alwaysRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	log                 logger.MomentoLogger
	reconnectMs         *int
}

type AlwaysRetryStrategy interface {
	Strategy
	WithEligibilityStrategy(s EligibilityStrategy) Strategy
	WithReconnectMs(ms int) Strategy
}

type AlwaysRetryStrategyProps struct {
	LoggerFactory       logger.MomentoLoggerFactory
	EligibilityStrategy EligibilityStrategy
	ReconnectMs         *int
}

func (r *alwaysRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	return &alwaysRetryStrategy{
		log:                 r.log,
		reconnectMs:         r.reconnectMs,
		eligibilityStrategy: s,
	}
}

func (r *alwaysRetryStrategy) WithReconnectMs(ms int) Strategy {
	return &alwaysRetryStrategy{
		log:                 r.log,
		reconnectMs:         &ms,
		eligibilityStrategy: r.eligibilityStrategy,
	}
}

func (r *alwaysRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	r.log.Debug(
		"Always retry strategy returning %d ms for [method: %s, status: %s]",
		*r.reconnectMs,
		props.GrpcMethod,
		props.GrpcStatusCode.String(),
	)
	return r.reconnectMs
}

func (r *alwaysRetryStrategy) String() string {
	return fmt.Sprintf(
		"alwaysRetryStrategy{eligibilityStrategy: %s, reconnectMs: %d, log: %s}\n",
		r.eligibilityStrategy,
		*r.reconnectMs,
		r.log,
	)
}

// NewAlwaysRetryStrategy is a retry strategy that always retries any request after a fixed delay.
// This maintains compatibility with the legacy behavior of the client, but is not the optimal
// implementation for most use cases. It is recommended to use a more sophisticated retry strategy.
func NewAlwaysRetryStrategy(props AlwaysRetryStrategyProps) Strategy {
	// the eligibility strategy is actually not used in this retry strategy as it duplicates
	// the legacy behavior of the client.
	eligibilityStrategy := EligibilityStrategy(DefaultEligibilityStrategy{})
	if props.EligibilityStrategy != nil {
		eligibilityStrategy = props.EligibilityStrategy
	}
	reconnectMsInt := 500
	reconnectMs := &reconnectMsInt
	if props.ReconnectMs != nil {
		reconnectMs = props.ReconnectMs
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("always-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("always-retry-strategy")
	}
	return &alwaysRetryStrategy{
		eligibilityStrategy: eligibilityStrategy,
		log:                 log,
		reconnectMs:         reconnectMs,
	}
}
