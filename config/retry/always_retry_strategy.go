package retry

import (
	"fmt"
	"github.com/momentohq/client-sdk-go/config/logger"
)

type alwaysRetryStrategy struct{
	eligibilityStrategy EligibilityStrategy
	log logger.MomentoLogger
	reconnectMs *int
}

type AlwaysRetryStrategy interface {
	Strategy
	// TODO: implement these
	WithEligibilityStrategy(s EligibilityStrategy) Strategy
	WithReconnectMs(ms int) Strategy
}

type AlwaysRetryStrategyProps struct {
	LoggerFactory       logger.MomentoLoggerFactory
	EligibilityStrategy EligibilityStrategy
	ReconnectMs         *int
}

func (r *alwaysRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	fmt.Printf("always retry determining when for props: %v", props)
	if !r.eligibilityStrategy.IsEligibleForRetry(props) {
		r.log.Debug(
		"Request is not retryable: [method: %s, status: %s]", props.GrpcMethod, props.GrpcStatusCode.String(),
		)
		//return nil
	}
	r.log.Debug("Request is retryable: [method: %s, status: %s]", props.GrpcMethod, props.GrpcStatusCode.String())
	r.log.Debug("returning reconnectMs: %d", *r.reconnectMs)
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
