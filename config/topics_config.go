package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
)

type topicsConfiguration struct {
	loggerFactory    logger.MomentoLoggerFactory
	maxSubscriptions uint32
}

type TopicsConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory

	MaxSubscriptions uint32
}

type TopicsConfiguration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetMaxSubscriptions Returns the configuration option for the maximum number of subscriptions
	// a client is allowed
	GetMaxSubscriptions() uint32

	WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration
}

func NewTopicConfiguration(props *TopicsConfigurationProps) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:    props.LoggerFactory,
		maxSubscriptions: props.MaxSubscriptions,
	}
}

func (s *topicsConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}

func (s *topicsConfiguration) GetMaxSubscriptions() uint32 {
	return s.maxSubscriptions
}

func (s *topicsConfiguration) WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:    s.loggerFactory,
		maxSubscriptions: maxSubscriptions,
	}
}
