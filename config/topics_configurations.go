// Package config provides pre-built configurations.
package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/retry"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/config/logger"
)

// TopicsDefault provides defaults configuration for a Topic Client
func TopicsDefault() TopicsConfiguration {
	return TopicsDefaultWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func TopicsDefaultWithLogger(loggerFactory logger.MomentoLoggerFactory) TopicsConfiguration {
	reconnectMs := 500
	return NewTopicConfiguration(&TopicsConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewTopicsStaticTransportStrategy(&TopicsTransportStrategyProps{
			GrpcConfiguration: NewTopicsStaticGrpcConfiguration(&TopicsGrpcConfigurationProps{
				client_timeout: 5 * time.Second,
			}),
		}),
		RetryStrategy: retry.NewLegacyTopicSubscriptionRetryStrategy(retry.LegacyTopicSubscriptionRetryStrategyProps{
			LoggerFactory: loggerFactory,
			RetryMs:       &reconnectMs,
		}),
	})
}
