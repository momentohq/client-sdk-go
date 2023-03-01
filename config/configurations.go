package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type Laptop struct {
	Configuration
}

func LatestLaptopConfig(loggerFactory ...logger.MomentoLoggerFactory) *Laptop {
	defaultLoggerFactory := logger.NewNoopMomentoLoggerFactory()
	if len(loggerFactory) != 0 {
		defaultLoggerFactory = loggerFactory[0]
	}
	return &Laptop{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			LoggerFactory: defaultLoggerFactory,
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline: 5 * time.Second,
				}),
			}),
		}),
	}
}

type InRegion struct {
	Configuration
}

func LatestInRegionConfig(loggerFactory ...logger.MomentoLoggerFactory) *InRegion {
	defaultLoggerFactory := logger.NewNoopMomentoLoggerFactory()
	if len(loggerFactory) != 0 {
		defaultLoggerFactory = loggerFactory[0]
	}
	return &InRegion{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			LoggerFactory: defaultLoggerFactory,
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline: 1100 * time.Millisecond,
				}),
			}),
		}),
	}
}
