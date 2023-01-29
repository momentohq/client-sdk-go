package config

import (
	"time"
)

type ConfigurationProps struct {
	Logger            Logger
	TransportStrategy TransportStrategy
}

type Configuration interface {
	// GetLogger Returns the currently configured logger with the Momento service
	GetLogger() Logger

	// WithCustomLogger Copy constructor for overriding Logger returns a new Configuration object
	// with the specified logger name.
	WithCustomLogger(logger Logger) Configuration

	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	//with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// WithClientTimeoutMillis Copy constructor for overriding TransportStrategy client side timeout. Returns a new Configuration object
	// with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeoutMillis(clientTimeoutMillis time.Duration) Configuration
}

type SimpleCacheConfiguration struct {
	logger            Logger
	transportStrategy TransportStrategy
}

func (s *SimpleCacheConfiguration) GetLogger() Logger {
	return s.logger
}

func (s *SimpleCacheConfiguration) GetClientSideTimeout() time.Duration {
	return s.transportStrategy.GetClientSideTimeout()
}

func NewSimpleCacheConfiguration(props *ConfigurationProps) Configuration {
	return &SimpleCacheConfiguration{
		logger:            props.Logger,
		transportStrategy: props.TransportStrategy,
	}
}

func (s *SimpleCacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *SimpleCacheConfiguration) WithCustomLogger(logger Logger) Configuration {
	return &SimpleCacheConfiguration{
		logger:            logger,
		transportStrategy: s.transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &SimpleCacheConfiguration{
		logger:            s.logger,
		transportStrategy: transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithClientTimeoutMillis(clientTimeout time.Duration) Configuration {
	return &SimpleCacheConfiguration{
		logger:            s.logger,
		transportStrategy: s.transportStrategy.WithClientTimeout(clientTimeout),
	}
}
