package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
)

type topicsConfiguration struct {
	loggerFactory logger.MomentoLoggerFactory
}

type TopicsConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory
}

type TopicsConfiguration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory
}

func NewTopicConfiguration(props *TopicsConfigurationProps) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory: props.LoggerFactory,
	}
}

func (s *topicsConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}
