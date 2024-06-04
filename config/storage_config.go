package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type storageConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
}

type StorageConfigurationProps struct {
	LoggerFactory     logger.MomentoLoggerFactory
	TransportStrategy TransportStrategy
}

type StorageConfiguration interface {
	GetLoggerFactory() logger.MomentoLoggerFactory
	GetTransportStrategy() TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) StorageConfiguration
	GetClientSideTimeout() time.Duration
	WithClientTimeout(clientTimeout time.Duration) StorageConfiguration
}

func NewStorageConfiguration(props *StorageConfigurationProps) StorageConfiguration {
	return &storageConfiguration{
		loggerFactory:     props.LoggerFactory,
		transportStrategy: props.TransportStrategy,
	}
}

func (c *storageConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return c.loggerFactory
}

func (c *storageConfiguration) GetTransportStrategy() TransportStrategy {
	return c.transportStrategy
}

func (c *storageConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) StorageConfiguration {
	return &storageConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: transportStrategy,
	}
}

func (c *storageConfiguration) GetClientSideTimeout() time.Duration {
	return c.transportStrategy.GetClientSideTimeout()
}

func (c *storageConfiguration) WithClientTimeout(clientTimeout time.Duration) StorageConfiguration {
	return &storageConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: c.transportStrategy.WithClientTimeout(clientTimeout),
	}
}
