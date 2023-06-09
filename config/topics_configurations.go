// Package config provides pre-built configurations.
package config

import (
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/config/logger"
)

// TopicsDefault provides defaults configuration for a Topic Client
func TopicsDefault() TopicsConfiguration {
	return TopicsDefaultWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func TopicsDefaultWithLogger(loggerFactory logger.MomentoLoggerFactory) TopicsConfiguration {
	return NewTopicConfiguration(&TopicsConfigurationProps{
		LoggerFactory: loggerFactory,
	})
}
