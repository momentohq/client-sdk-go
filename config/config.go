package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/internal/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type ConfigurationProps struct {
	LoggerFactory     logger.MomentoLoggerFactory
	TransportStrategy TransportStrategy
	RetryStrategy     retry.Strategy
}

type Configuration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetRetryStrategy Returns the current configuration options for wire interactions with the Momento service
	GetRetryStrategy() retry.Strategy

	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// WithRetryStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	// with the specified momento.TransportStrategy
	WithRetryStrategy(retryStrategy retry.Strategy) Configuration

	// WithClientTimeout Copy constructor for overriding TransportStrategy client side timeout. Returns a new
	//Configuration object with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeout(clientTimeout time.Duration) Configuration

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	// with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration
}

type SimpleCacheConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
	retryStrategy     retry.Strategy
}

func (s *SimpleCacheConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}

func (s *SimpleCacheConfiguration) GetClientSideTimeout() time.Duration {
	return s.transportStrategy.GetClientSideTimeout()
}

func NewSimpleCacheConfiguration(props *ConfigurationProps) Configuration {
	return &SimpleCacheConfiguration{
		loggerFactory:     props.LoggerFactory,
		transportStrategy: props.TransportStrategy,
		retryStrategy:     props.RetryStrategy,
	}
}

func (s *SimpleCacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *SimpleCacheConfiguration) GetRetryStrategy() retry.Strategy {
	return s.retryStrategy
}

func (s *SimpleCacheConfiguration) WithClientTimeout(clientTimeout time.Duration) Configuration {
	return &SimpleCacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithClientTimeout(clientTimeout),
		retryStrategy:     s.retryStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithRetryStrategy(strategy retry.Strategy) Configuration {
	return &SimpleCacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		retryStrategy:     strategy,
	}
}
func (s *SimpleCacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &SimpleCacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: transportStrategy,
		retryStrategy:     s.retryStrategy,
	}
}
