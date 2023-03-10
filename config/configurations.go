package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/internal/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

func LaptopLatest() Configuration {
	return LaptopLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func LaptopLatestWithLogger(loggerFactory logger.MomentoLoggerFactory) Configuration {
	return NewCacheConfiguration(&ConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 5 * time.Second,
			}),
		}),
		RetryStrategy: retry.NewFixedCountRetryStrategy(loggerFactory),
	})
}

func InRegionLatest() Configuration {
	return InRegionLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func InRegionLatestWithLogger(loggerFactory logger.MomentoLoggerFactory) Configuration {
	return NewCacheConfiguration(&ConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 1100 * time.Millisecond,
			}),
		}),
		RetryStrategy: retry.NewFixedCountRetryStrategy(loggerFactory),
	})
}
