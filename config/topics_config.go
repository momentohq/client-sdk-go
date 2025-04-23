package config

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

const MAX_CONCURRENT_STREAMS_PER_CHANNEL uint32 = 100

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

	// RetryStrategy is responsible for configuring the strategy for reconnecting a subscription to a
	// topic that has been interrupted. It is not applicable to publish requests.
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

	// GetRetryStrategy Returns the current strategy for topic subscription reconnection
	GetRetryStrategy() retry.Strategy

	// GetMiddleware Returns the list of middleware to be used by the topic client.
	GetMiddleware() []middleware.TopicMiddleware

	// WithMiddleware Copy constructor for overriding Middleware returns a new Configuration object
	// with the specified Middleware
	WithMiddleware(middleware []middleware.TopicMiddleware) TopicsConfiguration

	// AddMiddleware Copy constructor for adding Middleware returns a new Configuration object.
	AddMiddleware(m middleware.TopicMiddleware) TopicsConfiguration

	WithClientTimeout(clientTimeout time.Duration) TopicsConfiguration

	// GetGrpcConfig Configures the low-level gRPC settings for the Momento client's communication
	// with the Momento server.
	GetGrpcConfig() TopicsGrpcConfiguration
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
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return strategy.GetNumStreamGrpcChannels() * MAX_CONCURRENT_STREAMS_PER_CHANNEL
	case *TopicsDynamicTransportStrategy:
		return strategy.grpcConfig.GetMaxSubscriptions()
	}
	return s.maxSubscriptions
}

func (s *topicsConfiguration) WithMaxSubscriptions(maxSubscriptions uint32) TopicsConfiguration {
	// s.loggerFactory.GetLogger("TopicsConfiguration").Warn("WithMaxSubscriptions is deprecated, please use WithNumStreamGrpcChannels and WithNumUnaryGrpcChannels instead")
	// // If this deprecated method is used, we'll use the default 4 unary channels and set
	// // the number of stream channels to accommodate the specified number of subscriptions.
	// numStreamChannels := uint32(math.Ceil(float64(maxSubscriptions) / 100.0))
	// return &topicsConfiguration{
	// 	loggerFactory:     s.loggerFactory,
	// 	transportStrategy: s.transportStrategy.WithNumStreamGrpcChannels(numStreamChannels),
	// 	middleware:        s.middleware,
	// 	retryStrategy:     s.retryStrategy,
	// }
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithNumStreamGrpcChannels(maxSubscriptions / MAX_CONCURRENT_STREAMS_PER_CHANNEL),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	case *TopicsDynamicTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithMaxSubscriptions(maxSubscriptions),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	}
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
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
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithNumStreamGrpcChannels(numGrpcChannels),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	case *TopicsDynamicTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithMaxSubscriptions(numGrpcChannels * MAX_CONCURRENT_STREAMS_PER_CHANNEL),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	}
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetNumStreamGrpcChannels() uint32 {
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return strategy.GetNumStreamGrpcChannels()
	case *TopicsDynamicTransportStrategy:
		return strategy.grpcConfig.GetMaxSubscriptions() / MAX_CONCURRENT_STREAMS_PER_CHANNEL
	}
	return 0
}

func (s *topicsConfiguration) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithNumStreamGrpcChannels(numStreamGrpcChannels),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	case *TopicsDynamicTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithMaxSubscriptions(numStreamGrpcChannels * MAX_CONCURRENT_STREAMS_PER_CHANNEL),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	}
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
}

func (s *topicsConfiguration) GetNumUnaryGrpcChannels() uint32 {
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return strategy.GetNumUnaryGrpcChannels()
	case *TopicsDynamicTransportStrategy:
		return strategy.grpcConfig.GetNumUnaryGrpcChannels()
	}
	return 0
}

func (s *topicsConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsConfiguration {
	// maxSubscriptions and numGrpcChannels are deprecated, not included in the override
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithNumUnaryGrpcChannels(numUnaryGrpcChannels),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	case *TopicsDynamicTransportStrategy:
		return &topicsConfiguration{
			loggerFactory:     s.loggerFactory,
			transportStrategy: strategy.WithNumUnaryGrpcChannels(numUnaryGrpcChannels),
			middleware:        s.middleware,
			retryStrategy:     s.retryStrategy,
		}
	}
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
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

func (s *topicsConfiguration) WithClientTimeout(clientTimeout time.Duration) TopicsConfiguration {
	var updatedStrategy TopicsTransportStrategy
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		updatedStrategy = strategy.WithClientTimeout(clientTimeout)
	case *TopicsDynamicTransportStrategy:
		updatedStrategy = strategy.WithClientTimeout(clientTimeout)
	}
	return &topicsConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: updatedStrategy,
		middleware:        s.middleware,
		retryStrategy:     s.retryStrategy,
	}
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

func (s *topicsConfiguration) GetGrpcConfig() TopicsGrpcConfiguration {
	switch strategy := s.transportStrategy.(type) {
	case *TopicsStaticTransportStrategy:
		return strategy.GetGrpcConfig()
	case *TopicsDynamicTransportStrategy:
		return strategy.GetGrpcConfig()
	}
	return nil
}
