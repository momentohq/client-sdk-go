package retry

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/momentohq/client-sdk-go/config/logger"
)

const (
	DefaultInitialDelayMs = 0.5
	DefaultGrowthFactor   = 2
	DefaultMaxBackoffMs   = 8
)

type exponentialBackoffRetryStrategy struct {
	eligibilityStrategy EligibilityStrategy
	log                 logger.MomentoLogger
	initialDelayMillis  float64
	maxBackoffMillis    int
	growthFactor        int
}

type ExponentialBackoffRetryStrategy interface {
	Strategy
	WithInitialDelayMillis(float32) Strategy
	WithMaxBackoffMillis(int) Strategy
	WithEligibilityStrategy(EligibilityStrategy) Strategy
	WithGrowthFactor(int) Strategy
}

type ExponentialBackoffRetryStrategyProps struct {
	LoggerFactory       logger.MomentoLoggerFactory
	EligibilityStrategy EligibilityStrategy
	InitialDelayMillis  float64
	MaxBackoffMillis    int
	GrowthFactor        int
}

func NewExponentialBackoffRetryStrategy(props ExponentialBackoffRetryStrategyProps) Strategy {
	eligibilityStrategy := EligibilityStrategy(DefaultEligibilityStrategy{})
	if props.EligibilityStrategy != nil {
		eligibilityStrategy = props.EligibilityStrategy
	}
	var log logger.MomentoLogger
	if props.LoggerFactory != nil {
		log = props.LoggerFactory.GetLogger("exponential-backoff-retry-strategy")
	} else {
		log = logger.NewNoopMomentoLoggerFactory().GetLogger("exponential-backoff-retry-strategy")
	}
	initialDelayMillis := DefaultInitialDelayMs
	if props.InitialDelayMillis != 0 {
		initialDelayMillis = props.InitialDelayMillis
	}
	maxBackoffMillis := DefaultMaxBackoffMs
	if props.MaxBackoffMillis != 0 {
		maxBackoffMillis = props.MaxBackoffMillis
	}
	growthFactor := DefaultGrowthFactor
	if props.GrowthFactor != 0 {
		growthFactor = props.GrowthFactor
	}

	return &exponentialBackoffRetryStrategy{
		eligibilityStrategy: eligibilityStrategy,
		log:                 log,
		initialDelayMillis:  initialDelayMillis,
		maxBackoffMillis:    maxBackoffMillis,
		growthFactor:        growthFactor,
	}
}

func (r *exponentialBackoffRetryStrategy) WithInitialDelayMillis(delay float64) Strategy {
	return &exponentialBackoffRetryStrategy{
		eligibilityStrategy: r.eligibilityStrategy,
		log:                 r.log,
		initialDelayMillis:  delay,
		maxBackoffMillis:    r.maxBackoffMillis,
		growthFactor:        r.growthFactor,
	}
}

func (r *exponentialBackoffRetryStrategy) WithMaxBackoffMillis(backoff int) Strategy {
	return &exponentialBackoffRetryStrategy{
		eligibilityStrategy: r.eligibilityStrategy,
		log:                 r.log,
		initialDelayMillis:  r.initialDelayMillis,
		maxBackoffMillis:    backoff,
		growthFactor:        r.growthFactor,
	}
}

func (r *exponentialBackoffRetryStrategy) WithEligibilityStrategy(strategy EligibilityStrategy) Strategy {
	return &exponentialBackoffRetryStrategy{
		eligibilityStrategy: strategy,
		log:                 r.log,
		initialDelayMillis:  r.initialDelayMillis,
		maxBackoffMillis:    r.maxBackoffMillis,
		growthFactor:        r.growthFactor,
	}
}

func (r *exponentialBackoffRetryStrategy) WithGrowthFactor(growthFactor int) Strategy {
	return &exponentialBackoffRetryStrategy{
		eligibilityStrategy: r.eligibilityStrategy,
		log:                 r.log,
		initialDelayMillis:  r.initialDelayMillis,
		maxBackoffMillis:    r.maxBackoffMillis,
		growthFactor:        growthFactor,
	}
}

func (r *exponentialBackoffRetryStrategy) GetResponseDataReceivedTimeoutMillis() *int {
	return nil
}

func (r *exponentialBackoffRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	// attempt is 0-based, so we subtract 1 to get the correct attempt number
	attempt := props.AttemptNumber - 1
	r.log.Debug(
		"Determining whether request is eligible for retry; status code: %s, "+
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

func (r *exponentialBackoffRetryStrategy) String() string {
	return fmt.Sprintf("ExponentialBackoffRetryStrategy{eligibilityStrategy: %T, initialDelayMillis: %f, maxBackoffMillis: %d, growthFactor: %d}",
		r.eligibilityStrategy, r.initialDelayMillis, r.maxBackoffMillis, r.growthFactor)
}

func (r *exponentialBackoffRetryStrategy) computeBaseDelay(attemptNumber int) int {
	baseDelay := int(math.Pow(DefaultGrowthFactor, float64(attemptNumber)) * float64(r.initialDelayMillis))
	if baseDelay > r.maxBackoffMillis {
		return r.maxBackoffMillis
	}
	return baseDelay
}

func (r *exponentialBackoffRetryStrategy) computePreviousBaseDelay(baseDelay int) int {
	return baseDelay / r.growthFactor
}

func randInRange(min, max int) int {
	if min > max || min == max {
		return min
	}
	return rand.Intn(max-min) + min
}
