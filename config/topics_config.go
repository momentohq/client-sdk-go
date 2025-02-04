package config

import (
	"github.com/momentohq/client-sdk-go/config/logger"
)

type topicsConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	maxSubscriptions  uint32 // DEPRECATED
	numGrpcChannels   uint32 // DEPRECATED
	numStreamChannels uint32
	numUnaryChannels  uint32
}

type TopicsConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory

	// Deprecated: use NumStreamGrpcChannels and NumUnaryGrpcChannels instead.
	MaxSubscriptions uint32

	// Deprecated: use NumStreamGrpcChannels and NumUnaryGrpcChannels instead.
	NumGrpcChannels uint32

	// NumStreamGrpcChannels represents the number of GRPC channels the topic client
	// should open and work with for stream operations (i.e. topic subscriptions).
	NumStreamGrpcChannels uint32

	// NumUnaryGrpcChannels represents the number of GRPC channels the topic client
	// should open and work with for unary operations (i.e. topic publishes).
	NumUnaryGrpcChannels uint32
}

type TopicsConfiguration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetMaxSubscriptions Returns the configuration option for the maximum number of subscriptions
	// a client is allowed
	// Deprecated: Use GetNumGrpcChannels instead.
	GetMaxSubscriptions() uint32

	// Deprecated: please use the WithNumStreamGrpcChannels and WithNumUnaryGrpcChannels overrides
	// instead to tune the number of GRPC channels for stream and unary operations, respectively.
	// Defaults to 4 stream channels and 4 unary channels otherwise.
	//
	// WithMaxSubscriptions is currently implemented to create one GRPC connection for every
	// 100 subscribers. Can result in edge cases where subscribers and publishers are in contention
	// and may bottleneck a large volume of publish requests.
	// One GRPC connection can multiplex 100 subscribers/publishers.
	WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration

	// GetNumGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with.
	GetNumGrpcChannels() uint32

	// Deprecated: please use the WithNumStreamGrpcChannels and WithNumUnaryGrpcChannels overrides
	// instead to tune the number of GRPC channels for stream and unary operations, respectively.
	// Defaults to 4 stream channels and 4 unary channels otherwise.
	//
	// WithNumGrpcChannels is currently implemented to create the specified number of GRPC connections
	// (each GRPC connection can multiplex 100 subscribers/publishers). Defaults to 1.
	WithNumGrpcChannels(numGrpcChannels uint32) TopicsConfiguration

	// GetNumStreamGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for stream operations (i.e. topic subscriptions).
	GetNumStreamGrpcChannels() uint32

	// WithNumStreamGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for stream operations. Each GRPC connection can multiplex 100 concurrent subscriptions.
	// Defaults to 4.
	WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsConfiguration

	// GetNumUnaryGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for unary operations (i.e. topic publishes).
	GetNumUnaryGrpcChannels() uint32

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsConfiguration
}

func NewTopicConfiguration(props *TopicsConfigurationProps) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:     props.LoggerFactory,
		maxSubscriptions:  props.MaxSubscriptions,
		numGrpcChannels:   props.NumGrpcChannels,
		numStreamChannels: props.NumStreamGrpcChannels,
		numUnaryChannels:  props.NumUnaryGrpcChannels,
	}
}

func (s *topicsConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}

func (s *topicsConfiguration) GetMaxSubscriptions() uint32 {
	return s.maxSubscriptions
}

func (s *topicsConfiguration) WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration {
	// numGrpcChannels is mutually exclusive with maxSubscriptions, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		maxSubscriptions:  maxSubscriptions,
		numStreamChannels: s.numStreamChannels,
		numUnaryChannels:  s.numUnaryChannels,
	}
}

func (s *topicsConfiguration) GetNumGrpcChannels() uint32 {
	return s.numGrpcChannels
}

func (s *topicsConfiguration) WithNumGrpcChannels(numGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions is mutually exclusive with numGrpcChannels, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		numGrpcChannels:   numGrpcChannels,
		numStreamChannels: s.numStreamChannels,
		numUnaryChannels:  s.numUnaryChannels,
	}
}

func (s *topicsConfiguration) GetNumStreamGrpcChannels() uint32 {
	return s.numStreamChannels
}

func (s *topicsConfiguration) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		numStreamChannels: numStreamGrpcChannels,
		numUnaryChannels:  s.numUnaryChannels,
	}
}

func (s *topicsConfiguration) GetNumUnaryGrpcChannels() uint32 {
	return s.numUnaryChannels
}

func (s *topicsConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		numStreamChannels: s.numStreamChannels,
		numUnaryChannels:  numUnaryGrpcChannels,
	}
}
