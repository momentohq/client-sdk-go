package config

import (
	"time"
)

type TransportStrategyProps struct {
	// low-level gRPC settings for communication with the Momento server
	GrpcConfiguration GrpcConfiguration
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
}

type StaticGrpcConfiguration struct {
	deadline                    time.Duration
	keepAlivePermitWithoutCalls bool
	keepAliveTimeout            time.Duration
	keepAliveTime               time.Duration
	maxSendMessageLength        int
	maxReceiveMessageLength     int
}

func NewStaticGrpcConfiguration(grpcConfiguration *GrpcConfigurationProps) *StaticGrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:                    grpcConfiguration.deadline,
		keepAlivePermitWithoutCalls: grpcConfiguration.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            grpcConfiguration.keepAliveTimeout,
		keepAliveTime:               grpcConfiguration.keepAliveTime,
		maxSendMessageLength:        grpcConfiguration.maxSendMessageLength,
		maxReceiveMessageLength:     grpcConfiguration.maxReceiveMessageLength,
	}
}

func (s *StaticGrpcConfiguration) GetDeadline() time.Duration {
	return s.deadline
}

func (s *StaticGrpcConfiguration) GetKeepAlivePermitWithoutCalls() *bool {
	return &s.keepAlivePermitWithoutCalls
}

func (s *StaticGrpcConfiguration) GetKeepAliveTimeout() *time.Duration {
	return &s.keepAliveTimeout
}

func (s *StaticGrpcConfiguration) GetKeepAliveTime() *time.Duration {
	return &s.keepAliveTime
}

func (s *StaticGrpcConfiguration) GetMaxSendMessageLength() *int {
	return &s.maxSendMessageLength
}

func (s *StaticGrpcConfiguration) GetMaxReceiveMessageLength() *int {
	return &s.maxReceiveMessageLength
}

func (s *StaticGrpcConfiguration) WithDeadline(deadline time.Duration) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:                    deadline,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
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
