package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
)

type authConfiguration struct {
	loggerFactory logger.MomentoLoggerFactory
}

type AuthConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory
}

type AuthConfiguration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory
}

func NewAuthConfiguration(props *AuthConfigurationProps) AuthConfiguration {
	return &authConfiguration{
		loggerFactory: props.LoggerFactory,
	}
}

func (s *authConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}
