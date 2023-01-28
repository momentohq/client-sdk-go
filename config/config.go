package config

import "time"

type ConfigurationProps struct {
	TransportStrategy TransportStrategy
	RetryStrategy     RetryStrategy
}

type Configuration interface {
	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// GetRetryStrategy Returns the configuration options for the retry strategy to use with the Momento service
	GetRetryStrategy() RetryStrategy

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	//with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration

	// WithRetryStrategy Copy constructor for overriding RetryStrategy returns a new Configuration object
	WithRetryStrategy(retryStrategy RetryStrategy) Configuration

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding TransportStrategy client side timeout. Returns a new
	//Configuration object with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeout(clientTimeoutMillis time.Duration) Configuration
}

type SimpleCacheConfiguration struct {
	transportStrategy TransportStrategy
	retryStrategy     RetryStrategy
}

func NewSimpleCacheConfiguration(props *ConfigurationProps) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: props.TransportStrategy,
		retryStrategy:     props.RetryStrategy,
	}
}

func (s *SimpleCacheConfiguration) GetRetryStrategy() RetryStrategy {
	return s.retryStrategy
}

func (s *SimpleCacheConfiguration) WithRetryStrategy(retryStrategy RetryStrategy) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: s.transportStrategy,
		retryStrategy:     retryStrategy,
	}
}

func (s *SimpleCacheConfiguration) GetClientSideTimeout() time.Duration {
	return s.transportStrategy.GetClientSideTimeout()
}

func (s *SimpleCacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *SimpleCacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithClientTimeout(clientTimeout time.Duration) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: s.transportStrategy.WithClientTimeout(clientTimeout),
	}
}
