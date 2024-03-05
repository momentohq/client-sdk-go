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

// The default value for max_send_message_length is 4mb.  We need to increase this to 5mb in order to
// support cases where users have requested a limit increase up to our maximum item size of 5mb.
const DEFAULT_MAX_MESSAGE_SIZE = 5_243_000

const DEFAULT_KEEPALIVE_WITHOUT_STREAM = true
const DEFAULT_KEEPALIVE_TIME = 5000 * time.Millisecond
const DEFAULT_KEEPALIVE_TIMEOUT = 1000 * time.Millisecond

type StaticGrpcConfiguration struct {
	deadline                    time.Duration
	keepAlivePermitWithoutCalls bool
	keepAliveTimeout            time.Duration
	keepAliveTime               time.Duration
	maxSendMessageLength        int
	maxReceiveMessageLength     int
}

// Constructs new GrpcConfiguration to tune lower-level grpc settings.
// Note: keepalive settings are enabled by default, use WithKeepAliveDisabled() to disable all of them,
// or use the appropriate copy constructor to override individual settings.
func NewStaticGrpcConfiguration(grpcConfiguration *GrpcConfigurationProps) *StaticGrpcConfiguration {
	// We set keepalive values to defaults because we can't tell if users set them to zero values
	// or if the settings were just omitted from the props struct, thus defaulting them to zero values.

	maxSendLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.maxSendMessageLength > 0 {
		maxSendLength = grpcConfiguration.maxSendMessageLength
	}

	maxReceiveLength := DEFAULT_MAX_MESSAGE_SIZE
	if grpcConfiguration.maxReceiveMessageLength > 0 {
		maxReceiveLength = grpcConfiguration.maxReceiveMessageLength
	}

	return &StaticGrpcConfiguration{
		deadline:                    grpcConfiguration.deadline,
		keepAlivePermitWithoutCalls: DEFAULT_KEEPALIVE_WITHOUT_STREAM,
		keepAliveTimeout:            DEFAULT_KEEPALIVE_TIMEOUT,
		keepAliveTime:               DEFAULT_KEEPALIVE_TIME,
		maxSendMessageLength:        maxSendLength,
		maxReceiveMessageLength:     maxReceiveLength,
	}
}

func (s *StaticGrpcConfiguration) GetDeadline() time.Duration {
	return s.deadline
}

func (s *StaticGrpcConfiguration) GetKeepAlivePermitWithoutCalls() bool {
	return s.keepAlivePermitWithoutCalls
}

func (s *StaticGrpcConfiguration) GetKeepAliveTimeout() time.Duration {
	return s.keepAliveTimeout
}

func (s *StaticGrpcConfiguration) GetKeepAliveTime() time.Duration {
	return s.keepAliveTime
}

func (s *StaticGrpcConfiguration) GetMaxSendMessageLength() int {
	return s.maxSendMessageLength
}

func (s *StaticGrpcConfiguration) GetMaxReceiveMessageLength() int {
	return s.maxReceiveMessageLength
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

func (s *StaticGrpcConfiguration) WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:                    s.deadline,
		keepAlivePermitWithoutCalls: keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
	}
}

func (s *StaticGrpcConfiguration) WithKeepAliveTimeout(keepAliveTimeout time.Duration) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:                    s.deadline,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
	}
}

func (s *StaticGrpcConfiguration) WithKeepAliveTime(keepAliveTime time.Duration) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:                    s.deadline,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
	}
}

func (s *StaticGrpcConfiguration) WithKeepAliveDisabled() GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadline:                    s.deadline,
		keepAlivePermitWithoutCalls: false,
		keepAliveTimeout:            0,
		keepAliveTime:               0,
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
