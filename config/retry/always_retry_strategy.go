package retry

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type alwaysRetryStrategy struct {
	log         logger.MomentoLogger
	reconnectMs *int
}

type AlwaysRetryStrategy interface {
	Strategy
	WithReconnectMs(ms int) Strategy
}

type AlwaysRetryStrategyProps struct {
	LoggerFactory logger.MomentoLoggerFactory
	ReconnectMs   *int
}

func (r *alwaysRetryStrategy) WithReconnectMs(ms int) Strategy {
	return &alwaysRetryStrategy{
		log:         r.log,
		reconnectMs: &ms,
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
		"alwaysRetryStrategy{reconnectMs: %d, log: %s}\n",
		*r.reconnectMs,
		r.log,
	)
}

// NewAlwaysRetryStrategy is a retry strategy that always retries any request after a fixed delay.
// This maintains compatibility with the legacy behavior of the client, but is not the optimal
// implementation for most use cases. It is recommended to use a more sophisticated retry strategy.
func NewAlwaysRetryStrategy(props AlwaysRetryStrategyProps) Strategy {
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
		log:         log,
		reconnectMs: reconnectMs,
	}
}
