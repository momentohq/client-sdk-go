package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
)

func StoreDefault() StoreConfiguration {
	return StoreDefaultWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func StoreDefaultWithLogger(loggerFactory logger.MomentoLoggerFactory) StoreConfiguration {
	return NewStoreConfiguration(&StoreConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 5 * time.Second,
			}),
		}),
	})
}
