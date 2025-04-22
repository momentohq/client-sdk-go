package config

import (
	"fmt"
	"time"
)

const DEFAULT_MAX_SUBSCRIPTIONS = 400

type TopicsStaticTransportStrategyProps struct {
	// low-level gRPC settings for communication with the Momento server
	GrpcConfiguration TopicsStaticGrpcConfigurationType
}

type TopicsDynamicTransportStrategyProps struct {
	// low-level gRPC settings for communication with the Momento server
	GrpcConfiguration TopicsDynamicGrpcConfigurationType
}

// base interface for both static and dynamic transport strategies
type TopicsTransportStrategy interface {

	// GetClientSideTimeout Gets configuration for client side timeout from transport strategy
	GetClientSideTimeout() time.Duration
}

// static version

type TopicsStaticTransportStrategyType interface {
	TopicsTransportStrategy

	// GetGrpcConfig Configures the low-level gRPC settings for the Momento client's communication
	// with the Momento server.
	GetGrpcConfig() TopicsStaticGrpcConfigurationType

	// WithGrpcConfig Copy constructor for overriding the gRPC configuration. Returns  a new
	// TransportStrategy with the specified gRPC config.
	WithGrpcConfig(grpcConfig TopicsStaticGrpcConfigurationType) TopicsStaticTransportStrategyType

	// GetNumStreamGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for stream operations (i.e. topic subscriptions).
	GetNumStreamGrpcChannels() uint32

	// WithNumStreamGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for stream operations. Each GRPC connection can multiplex 100 concurrent subscriptions.
	// Defaults to 4.
	WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsStaticTransportStrategyType

	// GetClientSideTimeout Gets configuration for client side timeout from transport strategy
	// GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding the client sie timeout. Returns a new
	// TransportStrategy with the specified client side timeout.
	WithClientTimeout(clientTimeout time.Duration) TopicsStaticTransportStrategyType

	// GetNumUnaryGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for unary operations (i.e. topic publishes).
	GetNumUnaryGrpcChannels() uint32

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsStaticTransportStrategyType
}

type TopicsStaticTransportStrategy struct {
	grpcConfig TopicsStaticGrpcConfigurationType
}

func (s *TopicsStaticTransportStrategy) GetClientSideTimeout() time.Duration {
	return s.grpcConfig.GetClientTimeout()
}

func (s *TopicsStaticTransportStrategy) WithClientTimeout(clientTimeout time.Duration) TopicsStaticTransportStrategyType {
	return &TopicsStaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithClientTimeout(clientTimeout),
	}
}

func (s *TopicsStaticTransportStrategy) GetNumStreamGrpcChannels() uint32 {
	return s.grpcConfig.GetNumStreamGrpcChannels()
}

func (s *TopicsStaticTransportStrategy) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsStaticTransportStrategyType {
	return &TopicsStaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithNumStreamGrpcChannels(numStreamGrpcChannels),
	}
}

func (s *TopicsStaticTransportStrategy) GetNumUnaryGrpcChannels() uint32 {
	return s.grpcConfig.GetNumUnaryGrpcChannels()
}

func (s *TopicsStaticTransportStrategy) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsStaticTransportStrategyType {
	return &TopicsStaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithNumUnaryGrpcChannels(numUnaryGrpcChannels),
	}
}

func NewTopicsStaticTransportStrategy(props *TopicsStaticTransportStrategyProps) TopicsStaticTransportStrategyType {
	return &TopicsStaticTransportStrategy{
		grpcConfig: props.GrpcConfiguration,
	}
}

func (s *TopicsStaticTransportStrategy) GetGrpcConfig() TopicsStaticGrpcConfigurationType {
	return s.grpcConfig
}

func (s *TopicsStaticTransportStrategy) WithGrpcConfig(grpcConfig TopicsStaticGrpcConfigurationType) TopicsStaticTransportStrategyType {
	return &TopicsStaticTransportStrategy{
		grpcConfig: grpcConfig,
	}
}

func (s *TopicsStaticTransportStrategy) String() string {
	return fmt.Sprintf("TransportStrategy{grpcConfig=%v}", s.grpcConfig)
}

// dynamic version

type TopicsDynamicTransportStrategyType interface {
	TopicsTransportStrategy

	// GetGrpcConfig Configures the low-level gRPC settings for the Momento client's communication
	// with the Momento server.
	GetGrpcConfig() TopicsDynamicGrpcConfigurationType

	// WithGrpcConfig Copy constructor for overriding the gRPC configuration. Returns  a new
	// TransportStrategy with the specified gRPC config.
	WithGrpcConfig(grpcConfig TopicsDynamicGrpcConfigurationType) TopicsDynamicTransportStrategyType

	// WithMaxSubscriptions sets the maximum number of concurrent subscriptions a TopicClient can support.
	// Defaults to 400 subscriptions.
	WithMaxSubscriptions(maxSubscriptions uint32) TopicsDynamicTransportStrategyType

	// GetClientSideTimeout Gets configuration for client side timeout from transport strategy
	// GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding the client sie timeout. Returns a new
	// TransportStrategy with the specified client side timeout.
	WithClientTimeout(clientTimeout time.Duration) TopicsDynamicTransportStrategyType

	// GetNumUnaryGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for unary operations (i.e. topic publishes).
	GetNumUnaryGrpcChannels() uint32

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsDynamicTransportStrategyType
}

type TopicsDynamicTransportStrategy struct {
	grpcConfig TopicsDynamicGrpcConfigurationType
}

func (s TopicsDynamicTransportStrategy) WithMaxSubscriptions(maxSubscriptions uint32) TopicsDynamicTransportStrategyType {
	return &TopicsDynamicTransportStrategy{
		grpcConfig: s.grpcConfig.WithMaxSubscriptions(maxSubscriptions),
	}
}

func (s TopicsDynamicTransportStrategy) GetClientSideTimeout() time.Duration {
	return s.grpcConfig.GetClientTimeout()
}

func (s TopicsDynamicTransportStrategy) WithClientTimeout(clientTimeout time.Duration) TopicsDynamicTransportStrategyType {
	return &TopicsDynamicTransportStrategy{
		grpcConfig: s.grpcConfig.WithClientTimeout(clientTimeout),
	}
}

func (s TopicsDynamicTransportStrategy) GetNumUnaryGrpcChannels() uint32 {
	return s.grpcConfig.GetNumUnaryGrpcChannels()
}

func (s TopicsDynamicTransportStrategy) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsDynamicTransportStrategyType {
	return &TopicsDynamicTransportStrategy{
		grpcConfig: s.grpcConfig.WithNumUnaryGrpcChannels(numUnaryGrpcChannels),
	}
}

func NewTopicsDynamicTransportStrategy(props *TopicsDynamicTransportStrategyProps) TopicsDynamicTransportStrategyType {
	return &TopicsDynamicTransportStrategy{
		grpcConfig: props.GrpcConfiguration,
	}
}

func (s TopicsDynamicTransportStrategy) GetGrpcConfig() TopicsDynamicGrpcConfigurationType {
	return s.grpcConfig
}

func (s TopicsDynamicTransportStrategy) WithGrpcConfig(grpcConfig TopicsDynamicGrpcConfigurationType) TopicsDynamicTransportStrategyType {
	return &TopicsDynamicTransportStrategy{
		grpcConfig: grpcConfig,
	}
}

func (s *TopicsDynamicTransportStrategy) String() string {
	return fmt.Sprintf("TransportStrategy{grpcConfig=%v}", s.grpcConfig)
}
