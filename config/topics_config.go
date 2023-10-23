package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
)

type topicsConfiguration struct {
	loggerFactory    logger.MomentoLoggerFactory
	maxSubscriptions uint32
	numGrpcChannels  uint32
}

type TopicsConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory

	MaxSubscriptions uint32

	NumGrpcChannels uint32
}

type TopicsConfiguration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetMaxSubscriptions Returns the configuration option for the maximum number of subscriptions
	// a client is allowed
	// Deprecated: Use GetNumGrpcChannels instead.
	GetMaxSubscriptions() uint32

	// Deprecated: maxSubscriptions is currently implemented to create one GRPC connection for every
	// 100 subscribers. Can result in edge cases where subscribers and publishers are in contention
	// and may bottleneck a large volume of publish requests.
	//
	// Please use WithNumGrpcChannels instead as per your use case.
	// One GRPC connection can multiplex 100 subscribers/publishers.
	WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration

	// GetNumGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with.
	GetNumGrpcChannels() uint32

	WithNumGrpcChannels(numGrpcChannels uint32) TopicsConfiguration
}

func NewTopicConfiguration(props *TopicsConfigurationProps) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:    props.LoggerFactory,
		maxSubscriptions: props.MaxSubscriptions,
		numGrpcChannels:  props.NumGrpcChannels,
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

func (s *topicsConfiguration) GetNumGrpcChannels() uint32 {
	return s.numGrpcChannels
}

func (s *topicsConfiguration) WithNumGrpcChannels(numGrpcChannels uint32) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:   s.loggerFactory,
		numGrpcChannels: numGrpcChannels,
	}
}
