package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
)

func StorageLaptopLatest() StorageConfiguration {
	return StorageLaptopLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func StorageLaptopLatestWithLogger(loggerFactory logger.MomentoLoggerFactory) StorageConfiguration {
	return NewStorageConfiguration(&StorageConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 15 * time.Second,
			}),
		}),
		NumGrpcChannels: 1,
	})
}
