package retry

import (
	"github.com/momentohq/client-sdk-go/config/logger"
	"math"
	"math/rand"
)

const (
	DefaultInitialDelayMs = 0.5
	DefaultGrowthFactor = 2
	DefaultMaxBackoffMs   = 8
)

type exponentialBackoffRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	log logger.MomentoLogger
	initialDelayMillis float32
	maxBackoffMillis int
	growthFactor int
}

func NewExponentialBackoffRetryStrategy(logFactory logger.MomentoLoggerFactory) Strategy {
	return exponentialBackoffRetryStrategy{
		eligibilityStrategy: DefaultEligibilityStrategy{},
		log:                 logFactory.GetLogger("exponential-backoff-retry-strategy"),
		initialDelayMillis:  DefaultInitialDelayMs,
		maxBackoffMillis:    DefaultMaxBackoffMs,
		growthFactor:        DefaultGrowthFactor,
	}
}

func (r exponentialBackoffRetryStrategy) WithInitialDelayMillis(delay float32) Strategy {
	r.initialDelayMillis = delay
	return r
}

func (r exponentialBackoffRetryStrategy) WithMaxBackoffMillis(backoff int) Strategy {
	r.maxBackoffMillis = backoff
	return r
}

func (r exponentialBackoffRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	r.eligibilityStrategy = s
	return r
}

func (r exponentialBackoffRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	r.log.Debug(
		"Determining whether request is eligible for retry; status code: %s, " +
		"request type: %s, attemptNumber: %d",
		props.GrpcStatusCode, props.GrpcMethod, props.AttemptNumber,
	)

	if !r.eligibilityStrategy.IsEligibleForRetry(props) {
		r.log.Debug("Request is not eligible for retry.")
		return nil
	}

	baseDelay := r.computeBaseDelay(props.AttemptNumber)
	previousBaseDelay := r.computePreviousBaseDelay(baseDelay)
	maxDelay := previousBaseDelay * 3
	jitteredDelay := randInRange(baseDelay, maxDelay)
	r.log.Debug("attempt #%d, base delay: %d, previous base delay: %d, max delay: %d, jittered delay: %d",
		props.AttemptNumber, baseDelay, previousBaseDelay, maxDelay, jitteredDelay)
	return &jitteredDelay
}

func (r exponentialBackoffRetryStrategy) computeBaseDelay(attemptNumber int) int {
	baseDelay := int(math.Pow(DefaultGrowthFactor, float64(attemptNumber)) * float64(r.initialDelayMillis))
	if baseDelay > r.maxBackoffMillis {
		return r.maxBackoffMillis
	}
	return baseDelay
}

func (r exponentialBackoffRetryStrategy) computePreviousBaseDelay(baseDelay int) int {
	return baseDelay / r.growthFactor
}

func randInRange(min, max int) int {
	if min > max {
		return min
	}
	return rand.Intn(max-min) + min
}
