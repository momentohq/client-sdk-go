package config

import (
	"google.golang.org/grpc/codes"
	"time"
)

type RetryStrategyProps struct {
	// Which status codes to retry requests on
	RetryableRequestStatuses []codes.Code
	// Max number of times to retry
	MaxRetries uint
	// How long to wait before retrying a request
	PerRetryTimeout time.Duration
}

type RetryStrategy interface {
	// GetRetryableRequestStatuses Which status codes to retry requests on
	GetRetryableRequestStatuses() []codes.Code

	// GetMaxRetries max number of times to retry
	GetMaxRetries() uint

	// GetPerRetryTimeout How long to wait before retrying a request
	GetPerRetryTimeout() time.Duration

	// WithRetryableRequestStatuses CopyConstructor for setting which response codes to try on
	WithRetryableRequestStatuses(codes []codes.Code) RetryStrategy

	// WithMaxRetry CopyConstructor for setting max number of times to retry
	WithMaxRetry(maxRetries uint) RetryStrategy

	// WithPerRequestTimeout CopyConstructor for setting how long to wait for a request before retrying
	WithPerRequestTimeout(timeout time.Duration) RetryStrategy
}

type StaticRetryStrategy struct {
	retryableRequestStatuses []codes.Code
	maxRetries               uint
	perRetryTimeout          time.Duration
}

func NewStaticRetryStrategy(props *RetryStrategyProps) *StaticRetryStrategy {
	return &StaticRetryStrategy{
		retryableRequestStatuses: props.RetryableRequestStatuses,
		maxRetries:               props.MaxRetries,
		perRetryTimeout:          props.PerRetryTimeout,
	}
}

func (s StaticRetryStrategy) GetRetryableRequestStatuses() []codes.Code {
	return s.retryableRequestStatuses
}

func (s StaticRetryStrategy) GetMaxRetries() uint {
	return s.maxRetries
}

func (s StaticRetryStrategy) GetPerRetryTimeout() time.Duration {
	return s.perRetryTimeout
}

func (s StaticRetryStrategy) WithRetryableRequestStatuses(codes []codes.Code) RetryStrategy {
	return &StaticRetryStrategy{
		retryableRequestStatuses: codes,
		maxRetries:               s.maxRetries,
		perRetryTimeout:          s.perRetryTimeout,
	}
}

func (s StaticRetryStrategy) WithMaxRetry(maxRetries uint) RetryStrategy {
	return &StaticRetryStrategy{
		retryableRequestStatuses: s.retryableRequestStatuses,
		maxRetries:               maxRetries,
		perRetryTimeout:          s.perRetryTimeout,
	}
}

func (s StaticRetryStrategy) WithPerRequestTimeout(timeout time.Duration) RetryStrategy {
	return &StaticRetryStrategy{
		retryableRequestStatuses: s.retryableRequestStatuses,
		maxRetries:               s.maxRetries,
		perRetryTimeout:          timeout,
	}
}
