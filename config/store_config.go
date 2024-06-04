package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type storeConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
}

type StoreConfigurationProps struct {
	LoggerFactory     logger.MomentoLoggerFactory
	TransportStrategy TransportStrategy
}

type StoreConfiguration interface {
	GetLoggerFactory() logger.MomentoLoggerFactory
	GetTransportStrategy() TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) StoreConfiguration
	GetClientSideTimeout() time.Duration
	WithClientTimeout(clientTimeout time.Duration) StoreConfiguration
}

func NewStoreConfiguration(props *StoreConfigurationProps) StoreConfiguration {
	return &storeConfiguration{
		loggerFactory:     props.LoggerFactory,
		transportStrategy: props.TransportStrategy,
	}
}

func (c *storeConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return c.loggerFactory
}

func (c *storeConfiguration) GetTransportStrategy() TransportStrategy {
	return c.transportStrategy
}

func (c *storeConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) StoreConfiguration {
	return &storeConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: transportStrategy,
	}
}

func (c *storeConfiguration) GetClientSideTimeout() time.Duration {
	return c.transportStrategy.GetClientSideTimeout()
}

func (c *storeConfiguration) WithClientTimeout(clientTimeout time.Duration) StoreConfiguration {
	return &storeConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: c.transportStrategy.WithClientTimeout(clientTimeout),
	}
}
