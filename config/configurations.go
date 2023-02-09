package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
	"time"
)

type Laptop struct {
	Configuration
}

const defaultMaxSessionMemoryMb = 256

// 4 minutes.  We want to remain comfortably underneath the idle timeout for AWS NLB, which is 350s.
const defaultMaxIdle = 4 * time.Minute

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
					deadline:           5 * time.Second,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdle: defaultMaxIdle,
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
					deadline:           1100 * time.Millisecond,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdle: defaultMaxIdle,
			}),
		}),
	}
}
