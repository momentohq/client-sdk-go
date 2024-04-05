// Package config provides pre-built configurations.
package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/config/logger"
)

// LeaderboardDefault provides defaults configuration for a Leaderboard Client
func LeaderboardDefault() LeaderboardConfiguration {
	return LeaderboardDefaultWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func LeaderboardDefaultWithLogger(loggerFactory logger.MomentoLoggerFactory) LeaderboardConfiguration {
	return NewLeaderboardConfiguration(&LeaderboardConfigurationProps{
		LoggerFactory: loggerFactory,
		TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
			GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
				deadline: 5 * time.Second,
			}),
		}),
	})
}
