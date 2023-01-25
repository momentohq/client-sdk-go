package config

import "time"

type TransportStrategyProps struct {
	// low-level gRPC settings for communication with the Momento server
	GrpcConfiguration GrpcConfiguration

	// The maximum duration for which a connection may remain idle before being replaced.  This
	// setting can be used to force re-connection of a client if it has been idle for too long.
	// In environments such as AWS lambda, if the lambda is suspended for too long the connection
	// may be closed by the load balancer, resulting in an error on the subsequent request.  If
	// this setting is set to a duration less than the load balancer timeout, we can ensure that
	// the connection will be refreshed to avoid errors.
	MaxIdle time.Duration
}

type TransportStrategy interface {

	// GetGrpcConfig Configures the low-level gRPC settings for the Momento client's communication
	// with the Momento server.
	GetGrpcConfig() GrpcConfiguration

	// WithGrpcConfig Copy constructor for overriding the gRPC configuration. Returns  a new
	// TransportStrategy with the specified gRPC config.
	WithGrpcConfig(grpcConfig GrpcConfiguration) TransportStrategy

	// GetClientSideTimeout Gets configuration for client side timeout from transport strategy
	GetClientSideTimeout() time.Duration

	// WithClientTimeout Copy constructor for overriding the client sie timeout. Returns a new
	// TransportStrategy with the specified client side timeout.
	WithClientTimeout(clientTimeout time.Duration) TransportStrategy

	// WithMaxIdle Copy constructor for overriding the max idle connection timeout. Returns a new
	// TransportStrategy with the specified client side idle connection timeout.
	WithMaxIdle(maxIdle time.Duration) TransportStrategy
}

type StaticGrpcConfiguration struct {
	deadline           time.Duration
	maxSessionMemoryMb uint32
}

func NewStaticGrpcConfiguration(grpcConfiguration *GrpcConfigurationProps) *StaticGrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:           grpcConfiguration.deadline,
		maxSessionMemoryMb: grpcConfiguration.maxSessionMemoryMb,
	}
}

func (s *StaticGrpcConfiguration) GetDeadline() time.Duration {
	return s.deadline
}

func (s *StaticGrpcConfiguration) GetMaxSessionMemoryMb() uint32 {
	return s.maxSessionMemoryMb
}

func (s *StaticGrpcConfiguration) WithMaxSessionMb(maxSessionMemoryMb uint32) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:           s.deadline,
		maxSessionMemoryMb: maxSessionMemoryMb,
	}
}

func (s *StaticGrpcConfiguration) WithDeadline(deadline time.Duration) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:           deadline,
		maxSessionMemoryMb: s.maxSessionMemoryMb,
	}
}

type StaticTransportStrategy struct {
	grpcConfig GrpcConfiguration
	maxIdle    time.Duration
}

func (s *StaticTransportStrategy) GetClientSideTimeout() time.Duration {
	return s.grpcConfig.GetDeadline()
}

func (s *StaticTransportStrategy) WithClientTimeout(clientTimeout time.Duration) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig: s.grpcConfig.WithDeadline(clientTimeout),
		maxIdle:    s.maxIdle,
	}
}

func NewStaticTransportStrategy(props *TransportStrategyProps) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig: props.GrpcConfiguration,
		maxIdle:    props.MaxIdle,
	}
}

func (s *StaticTransportStrategy) GetGrpcConfig() GrpcConfiguration {
	return s.grpcConfig
}

func (s *StaticTransportStrategy) GetMaxIdle() time.Duration {
	return s.maxIdle
}

func (s *StaticTransportStrategy) WithGrpcConfig(grpcConfig GrpcConfiguration) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig: grpcConfig,
		maxIdle:    s.maxIdle,
	}
}

func (s *StaticTransportStrategy) WithMaxIdle(maxIdle time.Duration) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig: s.grpcConfig,
		maxIdle:    maxIdle,
	}
}
