package config

import (
	"math"
	"time"

	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type topicsConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	maxSubscriptions  uint32 // DEPRECATED
	numGrpcChannels   uint32 // DEPRECATED
	transportStrategy TopicsTransportStrategy
	retryStrategy     retry.Strategy
	middleware        []middleware.TopicMiddleware
}

type TopicsConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory

	// Deprecated: use NumStreamGrpcChannels and NumUnaryGrpcChannels instead.
	MaxSubscriptions uint32

	// Deprecated: use NumStreamGrpcChannels and NumUnaryGrpcChannels instead.
	NumGrpcChannels uint32

	// TransportStrategy is responsible for configuring network tunables for the GRPC client,
	// including the number of stream and unary grpc channels that should be used.
	TransportStrategy TopicsTransportStrategy

	RetryStrategy retry.Strategy

	// Middleware is a list of middleware to be used by the topic client.
	Middleware []middleware.TopicMiddleware
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
	// Using WithMaxSubscriptions now will default to creating 4 unary channels and as many stream
	// channels as needed to accommodate the maximum number of subscriptions.
	//
	// WithMaxSubscriptions creates one GRPC connection for every 100 subscribers.
	// Can result in edge cases where subscribers and publishers are in contention
	// and may bottleneck a large volume of publish requests.
	// One GRPC connection can multiplex 100 subscribers/publishers.
	WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration

	// GetNumGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with.
	GetNumGrpcChannels() uint32

	// Deprecated: please use the WithNumStreamGrpcChannels and WithNumUnaryGrpcChannels overrides
	// instead to tune the number of GRPC channels for stream and unary operations, respectively.
	// Using WithNumGrpcChannels now will default creating 4 unary channels and `numGrpcChannels`
	// number of stream channels.
	//
	// WithNumGrpcChannels creates the specified number of GRPC connections
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

	GetTransportStrategy() TopicsTransportStrategy

	WithTransportStrategy(transportStrategy TopicsTransportStrategy) TopicsConfiguration

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	WithRetryStrategy(retryStrategy retry.Strategy) TopicsConfiguration

	// GetRetryStrategy Returns the current configuration options for wire interactions with the Momento service
	GetRetryStrategy() retry.Strategy

	// GetMiddleware Returns the list of middleware to be used by the topic client.
	GetMiddleware() []middleware.TopicMiddleware

	// WithMiddleware Copy constructor for overriding Middleware returns a new Configuration object
	// with the specified Middleware
	WithMiddleware(middleware []middleware.TopicMiddleware) TopicsConfiguration

	// AddMiddleware Copy constructor for adding Middleware returns a new Configuration object.
	AddMiddleware(m middleware.TopicMiddleware) TopicsConfiguration
}

func NewTopicConfiguration(props *TopicsConfigurationProps) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:     props.LoggerFactory,
		maxSubscriptions:  props.MaxSubscriptions,
		numGrpcChannels:   props.NumGrpcChannels,
		transportStrategy: props.TransportStrategy,
		middleware:        props.Middleware,
		retryStrategy:     props.RetryStrategy,
	}
}

func (s *topicsConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}

func (s *topicsConfiguration) GetMaxSubscriptions() uint32 {
	return s.maxSubscriptions
}

func (s *topicsConfiguration) WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration {
	s.loggerFactory.GetLogger("TopicsConfiguration").Warn("WithMaxSubscriptions is deprecated, please use WithNumStreamGrpcChannels and WithNumUnaryGrpcChannels instead")
	// If this deprecated method is used, we'll use the default 4 unary channels and set
	// the number of stream channels to accommodate the specified number of subscriptions.
	numStreamChannels := uint32(math.Ceil(float64(maxSubscriptions) / 100.0))
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithNumStreamGrpcChannels(numStreamChannels),
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetNumGrpcChannels() uint32 {
	return s.numGrpcChannels
}

func (s *topicsConfiguration) WithNumGrpcChannels(numGrpcChannels uint32) TopicsConfiguration {
	s.loggerFactory.GetLogger("TopicsConfiguration").Warn("WithNumGrpcChannels is deprecated, please use WithNumStreamGrpcChannels and WithNumUnaryGrpcChannels instead")
	// If this deprecated method is used, we'll use the default 4 unary channels
	// and set the number of stream channels to the specified number.
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithNumStreamGrpcChannels(numGrpcChannels),
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetNumStreamGrpcChannels() uint32 {
	return s.transportStrategy.GetNumStreamGrpcChannels()
}

func (s *topicsConfiguration) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithNumStreamGrpcChannels(numStreamGrpcChannels),
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetNumUnaryGrpcChannels() uint32 {
	return s.transportStrategy.GetNumUnaryGrpcChannels()
}

func (s *topicsConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithNumUnaryGrpcChannels(numUnaryGrpcChannels),
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetTransportStrategy() TopicsTransportStrategy {
	return s.transportStrategy
}

func (s *topicsConfiguration) WithTransportStrategy(transportStrategy TopicsTransportStrategy) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: transportStrategy,
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetClientSideTimeout() time.Duration {
	return s.transportStrategy.GetClientSideTimeout()
}

func (s *topicsConfiguration) GetMiddleware() []middleware.TopicMiddleware {
	return s.middleware
}

func (s *topicsConfiguration) WithMiddleware(middleware []middleware.TopicMiddleware) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		middleware:        middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) AddMiddleware(m middleware.TopicMiddleware) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		middleware:        append(s.middleware, m),
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) WithRetryStrategy(retryStrategy retry.Strategy) TopicsConfiguration {
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		middleware:        s.middleware,
		retryStrategy:     retryStrategy,
	}
}

func (s *topicsConfiguration) GetRetryStrategy() retry.Strategy {
	return s.retryStrategy
}
