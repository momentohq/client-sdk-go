package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type LeaderboardConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory

	// TransportStrategy is responsible for configuring network tunables.
	TransportStrategy TransportStrategy
}

type LeaderboardConfiguration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	// with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) LeaderboardConfiguration

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding TransportStrategy client side timeout. Returns a new
	// Configuration object with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeout(clientTimeout time.Duration) LeaderboardConfiguration
}

type leaderboardConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
}

func NewLeaderboardConfiguration(props *LeaderboardConfigurationProps) LeaderboardConfiguration {
	return &leaderboardConfiguration{
		loggerFactory:     props.LoggerFactory,
		transportStrategy: props.TransportStrategy,
	}
}

func (c *leaderboardConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return c.loggerFactory
}

func (c *leaderboardConfiguration) GetTransportStrategy() TransportStrategy {
	return c.transportStrategy
}

func (c *leaderboardConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) LeaderboardConfiguration {
	return &leaderboardConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: transportStrategy,
	}
}

func (c *leaderboardConfiguration) GetClientSideTimeout() time.Duration {
	return c.transportStrategy.GetClientSideTimeout()
}

func (c *leaderboardConfiguration) WithClientTimeout(clientTimeout time.Duration) LeaderboardConfiguration {
	return &leaderboardConfiguration{
		loggerFactory:     c.loggerFactory,
		transportStrategy: c.transportStrategy.WithClientTimeout(clientTimeout),
	}
}
