// Package config provides pre-built configurations.
package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/internal/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

// LaptopLatest provides defaults suitable for a medium-to-high-latency dev environment.
// Permissive timeouts, retries, and relaxed latency and throughput targets.
// It enables keep-alive pings by default because they are very important for long-lived server
// environments where there may be periods of time when the connection is idle.
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
		RetryStrategy:   retry.NewFixedCountRetryStrategy(loggerFactory),
		NumGrpcChannels: 1,
		ReadConcern:     BALANCED,
	})
}

// InRegionLatest provides defaults suitable for an environment where your client is running in the same region as the Momento service.
// It has more aggressive timeouts and retry behavior than the Laptop config.
// It enables keep-alive pings by default because they are very important for long-lived server
// environments where there may be periods of time when the connection is idle.

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
		RetryStrategy:   retry.NewFixedCountRetryStrategy(loggerFactory),
		NumGrpcChannels: 1,
		ReadConcern:     BALANCED,
	})
}

// LambdaLatest provides defaults suitable for an environment where your client is running in
// a serverless environment like AWS Lambda.
// NOTE: keep-alives are very important for long-lived server environments where there may be periods of time
// when the connection is idle. However, they are very problematic for lambda environments where the lambda
// runtime is continuously frozen and unfrozen, because the lambda may be frozen before the "ACK" is received
// from the server. This can cause the keep-alive to timeout even though the connection is completely healthy.
// Therefore, keep-alives should be disabled in lambda and similar environments.

func LambdaLatest() Configuration {
	return LambdaLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func LambdaLatestWithLogger(loggerFactory logger.MomentoLoggerFactory) Configuration {
	return NewCacheConfiguration(&ConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 1100 * time.Millisecond,
			}).WithKeepAliveDisabled(),
		}),
		RetryStrategy:   retry.NewFixedCountRetryStrategy(loggerFactory),
		NumGrpcChannels: 1,
		ReadConcern:     BALANCED,
	})
}
