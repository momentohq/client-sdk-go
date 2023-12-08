// Package config provides pre-built configurations.
package config

import (
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/config/logger"
)

// AuthDefault provides defaults configuration for a Auth Client
func AuthDefault() AuthConfiguration {
	return AuthDefaultWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO))
}

func AuthDefaultWithLogger(loggerFactory logger.MomentoLoggerFactory) AuthConfiguration {
	return NewAuthConfiguration(&AuthConfigurationProps{
		LoggerFactory: loggerFactory,
	})
}
