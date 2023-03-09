// Package config provides pre-built configurations.
package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/internal/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

// LaptopLatest provides defaults suitable for a medium-to-high-latency dev environment.
// Permissive timeouts, retries, and relaxed latency and throughput targets.
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
		RetryStrategy: retry.NewFixedCountRetryStrategy(defaultLoggerFactory),
	})
}

// InRegionLatest provides defaults suitable for an environment where your client is running in the same region as the Momento service.
// It has more agressive timeouts and retry behavior than the Laptop config.
func InRegionLatest(loggerFactory ...logger.MomentoLoggerFactory) Configuration {
	defaultLoggerFactory := logger.NewNoopMomentoLoggerFactory()
	if len(loggerFactory) != 0 {
		defaultLoggerFactory = loggerFactory[0]
	}
	return NewCacheConfiguration(&ConfigurationProps{
		LoggerFactory: defaultLoggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 1100 * time.Millisecond,
			}),
		}),
		RetryStrategy: retry.NewFixedCountRetryStrategy(defaultLoggerFactory),
	})
}
