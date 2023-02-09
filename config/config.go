package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
	"time"
)

type ConfigurationProps struct {
	LoggerFactory     logger.MomentoLoggerFactory
	TransportStrategy TransportStrategy
}

type Configuration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	//with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding TransportStrategy client side timeout. Returns a new
	//Configuration object with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeout(clientTimeout time.Duration) Configuration
}

type SimpleCacheConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
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
	}
}

func (s *SimpleCacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *SimpleCacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &SimpleCacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithClientTimeout(clientTimeout time.Duration) Configuration {
	return &SimpleCacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithClientTimeout(clientTimeout),
	}
}
