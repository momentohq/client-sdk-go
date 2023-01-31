package config

import (
	"time"
)

type ConfigurationProps struct {
	LoggerOptions     LoggerOptions
	TransportStrategy TransportStrategy
}

type Configuration interface {
	// GetLogger Returns the currently configured logger with the Momento service
	GetLogger(loggerType loggerType) Logger

	// GetLoggerOptions Returns the currently configured logger options with the Momento service
	GetLoggerOptions() LoggerOptions

	// WithCustomLoggerOptions Copy constructor for overriding LoggerOptions returns a new Configuration object
	// with the specified logger options.
	WithCustomLoggerOptions(loggerOptions LoggerOptions) Configuration

	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	//with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding TransportStrategy client side timeout. Returns a new
	//Configuration object with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeout(clientTimeoutMillis time.Duration) Configuration
}

type SimpleCacheConfiguration struct {
	loggerOptions     LoggerOptions
	transportStrategy TransportStrategy
}

func (s *SimpleCacheConfiguration) GetLoggerOptions() LoggerOptions {
	return s.loggerOptions
}

func (s *SimpleCacheConfiguration) GetClientSideTimeout() time.Duration {
	return s.transportStrategy.GetClientSideTimeout()
}

func NewSimpleCacheConfiguration(props *ConfigurationProps) Configuration {
	return &SimpleCacheConfiguration{
		loggerOptions:     props.LoggerOptions,
		transportStrategy: props.TransportStrategy,
	}
}

func (s *SimpleCacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *SimpleCacheConfiguration) WithCustomLoggerOptions(loggerOptions LoggerOptions) Configuration {
	return &SimpleCacheConfiguration{
		loggerOptions:     loggerOptions,
		transportStrategy: s.transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &SimpleCacheConfiguration{
		loggerOptions:     s.loggerOptions,
		transportStrategy: transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithClientTimeout(clientTimeout time.Duration) Configuration {
	return &SimpleCacheConfiguration{
		loggerOptions:     s.loggerOptions,
		transportStrategy: s.transportStrategy.WithClientTimeout(clientTimeout),
	}
}

func (s *SimpleCacheConfiguration) GetLogger(loggerType loggerType) Logger {
	if loggerType == builtin {
		return NewBuiltInLogger(&s.loggerOptions)
	}
	return NewBuiltInLogger(&s.loggerOptions)
}
