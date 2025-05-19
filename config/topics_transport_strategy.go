package config

import (
	"fmt"
	"time"
)

type TopicsTransportStrategyProps struct {
	// low-level gRPC settings for communication with the Momento server
	GrpcConfiguration TopicsGrpcConfiguration
}

type TopicsTransportStrategy interface {

	// GetGrpcConfig Configures the low-level gRPC settings for the Momento client's communication
	// with the Momento server.
	GetGrpcConfig() TopicsGrpcConfiguration

	// WithGrpcConfig Copy constructor for overriding the gRPC configuration. Returns  a new
	// TransportStrategy with the specified gRPC config.
	WithGrpcConfig(grpcConfig TopicsGrpcConfiguration) TopicsTransportStrategy

	// GetClientSideTimeout Gets configuration for client side timeout from transport strategy
	GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding the client sie timeout. Returns a new
	// TransportStrategy with the specified client side timeout.
	WithClientTimeout(clientTimeout time.Duration) TopicsTransportStrategy

	// GetNumStreamGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for stream operations (i.e. topic subscriptions).
	GetNumStreamGrpcChannels() uint32

	// WithNumStreamGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for stream operations. Each GRPC connection can multiplex 100 concurrent subscriptions.
	// Defaults to 4.
	WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsTransportStrategy

	// GetNumUnaryGrpcChannels Returns the configuration option for the number of GRPC channels
	// the topic client should open and work with for unary operations (i.e. topic publishes).
	GetNumUnaryGrpcChannels() uint32

	// WithNumUnaryGrpcChannels is currently implemented to create the specified number of GRPC connections
	// for unary operations. Each GRPC connection can multiplex 100 concurrent publish requests.
	// Defaults to 4.
	WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsTransportStrategy
}

type TopicsStaticTransportStrategy struct {
	grpcConfig TopicsGrpcConfiguration
}

func (s *TopicsStaticTransportStrategy) GetClientSideTimeout() time.Duration {
	return s.grpcConfig.GetClientTimeout()
}

func (s *TopicsStaticTransportStrategy) WithClientTimeout(clientTimeout time.Duration) TopicsTransportStrategy {
	return &TopicsStaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithClientTimeout(clientTimeout),
	}
}

func (s *TopicsStaticTransportStrategy) GetNumStreamGrpcChannels() uint32 {
	return s.grpcConfig.GetNumStreamGrpcChannels()
}

func (s *TopicsStaticTransportStrategy) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsTransportStrategy {
	return &TopicsStaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithNumStreamGrpcChannels(numStreamGrpcChannels),
	}
}

func (s *TopicsStaticTransportStrategy) GetNumUnaryGrpcChannels() uint32 {
	return s.grpcConfig.GetNumUnaryGrpcChannels()
}

func (s *TopicsStaticTransportStrategy) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsTransportStrategy {
	return &TopicsStaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithNumUnaryGrpcChannels(numUnaryGrpcChannels),
	}
}

func NewTopicsStaticTransportStrategy(props *TopicsTransportStrategyProps) TopicsTransportStrategy {
	return &TopicsStaticTransportStrategy{
		grpcConfig: props.GrpcConfiguration,
	}
}

func (s *TopicsStaticTransportStrategy) GetGrpcConfig() TopicsGrpcConfiguration {
	return s.grpcConfig
}

func (s *TopicsStaticTransportStrategy) WithGrpcConfig(grpcConfig TopicsGrpcConfiguration) TopicsTransportStrategy {
	return &TopicsStaticTransportStrategy{
		grpcConfig: grpcConfig,
	}
}

func (s *TopicsStaticTransportStrategy) String() string {
	return fmt.Sprintf("TransportStrategy{grpcConfig=%v}", s.grpcConfig)
}
