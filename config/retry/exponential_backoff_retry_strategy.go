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

type ExponentialBackoffRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	log                 logger.MomentoLogger
	initialDelayMillis float32
	maxBackoffMillis int
	growthFactor int
}

func NewExponentialBackoffRetryStrategy(logFactory logger.MomentoLoggerFactory) Strategy {
	return ExponentialBackoffRetryStrategy{
		eligibilityStrategy: DefaultEligibilityStrategy{},
		log:                 logFactory.GetLogger("exponential-backoff-retry-strategy"),
		initialDelayMillis:  DefaultInitialDelayMs,
		maxBackoffMillis:    DefaultMaxBackoffMs,
		growthFactor:        DefaultGrowthFactor,
	}
}

func (r ExponentialBackoffRetryStrategy) WithInitialDelayMillis(delay float32) Strategy {
	r.initialDelayMillis = delay
	return r
}

func (r ExponentialBackoffRetryStrategy) WithMaxBackoffMillis(backoff int) Strategy {
	r.maxBackoffMillis = backoff
	return r
}

func (r ExponentialBackoffRetryStrategy) WithEligibilityStrategy(s EligibilityStrategy) Strategy {
	r.eligibilityStrategy = s
	return r
}

func (r ExponentialBackoffRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	// attempt is 0-based, so we subtract 1 to get the correct attempt number
	attempt := props.AttemptNumber - 1
	r.log.Debug(
		"Determining whether request is eligible for retry; status code: %s, " +
		"request type: %s, attemptNumber: %d",
		props.GrpcStatusCode, props.GrpcMethod, attempt,
	)

	if !r.eligibilityStrategy.IsEligibleForRetry(props) {
		r.log.Debug("Request is not eligible for retry.")
		return nil
	}

	baseDelay := r.computeBaseDelay(attempt)
	previousBaseDelay := r.computePreviousBaseDelay(baseDelay)
	maxDelay := previousBaseDelay * 3
	jitteredDelay := randInRange(baseDelay, maxDelay)
	r.log.Debug("attempt #%d, base delay: %d, previous base delay: %d, max delay: %d, jittered delay: %d",
		attempt, baseDelay, previousBaseDelay, maxDelay, jitteredDelay)
	return &jitteredDelay
}

func (r ExponentialBackoffRetryStrategy) computeBaseDelay(attemptNumber int) int {
	baseDelay := int(math.Pow(DefaultGrowthFactor, float64(attemptNumber)) * float64(r.initialDelayMillis))
	if baseDelay > r.maxBackoffMillis {
		return r.maxBackoffMillis
	}
	return baseDelay
}

func (r ExponentialBackoffRetryStrategy) computePreviousBaseDelay(baseDelay int) int {
	return baseDelay / r.growthFactor
}

func randInRange(min, max int) int {
	if min > max || min == max {
		return min
	}
	return rand.Intn(max-min) + min
}
