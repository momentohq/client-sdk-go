package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

func LaptopLatest(loggerFactory ...logger.MomentoLoggerFactory) Configuration {
	defaultLoggerFactory := logger.NewNoopMomentoLoggerFactory()
	if len(loggerFactory) != 0 {
		defaultLoggerFactory = loggerFactory[0]
	}
	return NewCacheConfiguration(&ConfigurationProps{
		LoggerFactory: defaultLoggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 5 * time.Second,
			}),
		}),
	})
}

func InRegionLatest(loggerFactory ...logger.MomentoLoggerFactory) Configuration {
	defaultLoggerFactory := logger.NewNoopMomentoLoggerFactory()
	if len(loggerFactory) != 0 {
		defaultLoggerFactory = loggerFactory[0]
	}
	return NewCacheConfiguration(
		&ConfigurationProps{
			LoggerFactory: defaultLoggerFactory,
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline: 1100 * time.Millisecond,
				}),
			}),
		},
	)
}
