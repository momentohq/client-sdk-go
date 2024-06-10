package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type storageConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
	numGrpcChannels   uint32
}

type StorageConfigurationProps struct {
	LoggerFactory     logger.MomentoLoggerFactory
	TransportStrategy TransportStrategy
	NumGrpcChannels   uint32
}

type StorageConfiguration interface {
	GetLoggerFactory() logger.MomentoLoggerFactory
	GetTransportStrategy() TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) StorageConfiguration
	GetClientSideTimeout() time.Duration
	WithClientTimeout(clientTimeout time.Duration) StorageConfiguration
	GetNumGrpcChannels() uint32
	WithNumGrpcChannels(numGrpcChannels uint32) StorageConfiguration
}

func NewStorageConfiguration(props *StorageConfigurationProps) StorageConfiguration {
	return &storageConfiguration{
		loggerFactory:     props.LoggerFactory,
		transportStrategy: props.TransportStrategy,
		numGrpcChannels:   props.NumGrpcChannels,
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
		numGrpcChannels:   c.numGrpcChannels,
	}
}

func (c *storageConfiguration) GetClientSideTimeout() time.Duration {
	return c.transportStrategy.GetClientSideTimeout()
}

func (c *storageConfiguration) WithClientTimeout(clientTimeout time.Duration) StorageConfiguration {
	return &storageConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: c.transportStrategy.WithClientTimeout(clientTimeout),
		numGrpcChannels:   c.numGrpcChannels,
	}
}

func (c *storageConfiguration) GetNumGrpcChannels() uint32 {
	return c.numGrpcChannels
}

func (c *storageConfiguration) WithNumGrpcChannels(numGrpcChannels uint32) StorageConfiguration {
	return &storageConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: c.transportStrategy,
		numGrpcChannels:   numGrpcChannels,
	}
}
