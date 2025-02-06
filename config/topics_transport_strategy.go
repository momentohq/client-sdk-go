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

type TopicsStaticGrpcConfiguration struct {
	client_timeout              time.Duration
	keepAlivePermitWithoutCalls bool
	keepAliveTimeout            time.Duration
	keepAliveTime               time.Duration
	maxSendMessageLength        int
	maxReceiveMessageLength     int
	numStreamGrpcChannels       uint32
	numUnaryGrpcChannels        uint32
}

// NewStaticGrpcConfiguration constructs a new TopicsGrpcConfiguration to tune lower-level grpc settings.
// Note: keepalive settings are enabled by default, use WithKeepAliveDisabled() to disable all of them,
// or use the appropriate copy constructor to override individual settings.
func NewTopicsStaticGrpcConfiguration(grpcConfiguration *TopicsGrpcConfigurationProps) *TopicsStaticGrpcConfiguration {
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

	return &TopicsStaticGrpcConfiguration{
		client_timeout:              grpcConfiguration.client_timeout,
		keepAlivePermitWithoutCalls: DEFAULT_KEEPALIVE_WITHOUT_STREAM,
		keepAliveTimeout:            DEFAULT_KEEPALIVE_TIMEOUT,
		keepAliveTime:               DEFAULT_KEEPALIVE_TIME,
		maxSendMessageLength:        maxSendLength,
		maxReceiveMessageLength:     maxReceiveLength,
		numStreamGrpcChannels:       grpcConfiguration.numStreamGrpcChannels,
		numUnaryGrpcChannels:        grpcConfiguration.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) GetClientTimeout() time.Duration {
	return s.client_timeout
}

func (s *TopicsStaticGrpcConfiguration) GetKeepAlivePermitWithoutCalls() bool {
	return s.keepAlivePermitWithoutCalls
}

func (s *TopicsStaticGrpcConfiguration) GetKeepAliveTimeout() time.Duration {
	return s.keepAliveTimeout
}

func (s *TopicsStaticGrpcConfiguration) GetKeepAliveTime() time.Duration {
	return s.keepAliveTime
}

func (s *TopicsStaticGrpcConfiguration) GetMaxSendMessageLength() int {
	return s.maxSendMessageLength
}

func (s *TopicsStaticGrpcConfiguration) GetMaxReceiveMessageLength() int {
	return s.maxReceiveMessageLength
}

func (s *TopicsStaticGrpcConfiguration) GetNumStreamGrpcChannels() uint32 {
	return s.numStreamGrpcChannels
}

func (s *TopicsStaticGrpcConfiguration) GetNumUnaryGrpcChannels() uint32 {
	return s.numUnaryGrpcChannels
}

func (s *TopicsStaticGrpcConfiguration) WithClientTimeout(client_timeout time.Duration) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAlivePermitWithoutCalls(keepAlivePermitWithoutCalls bool) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveTimeout(keepAliveTimeout time.Duration) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveTime(keepAliveTime time.Duration) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithKeepAliveDisabled() TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: false,
		keepAliveTimeout:            0,
		keepAliveTime:               0,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithNumStreamGrpcChannels(numStreamGrpcChannels uint32) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       numStreamGrpcChannels,
		numUnaryGrpcChannels:        s.numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) WithNumUnaryGrpcChannels(numUnaryGrpcChannels uint32) TopicsGrpcConfiguration {
	return &TopicsStaticGrpcConfiguration{
		client_timeout:              s.client_timeout,
		keepAlivePermitWithoutCalls: s.keepAlivePermitWithoutCalls,
		keepAliveTimeout:            s.keepAliveTimeout,
		keepAliveTime:               s.keepAliveTime,
		maxSendMessageLength:        s.maxSendMessageLength,
		maxReceiveMessageLength:     s.maxReceiveMessageLength,
		numStreamGrpcChannels:       s.numStreamGrpcChannels,
		numUnaryGrpcChannels:        numUnaryGrpcChannels,
	}
}

func (s *TopicsStaticGrpcConfiguration) String() string {
	return fmt.Sprintf(
		"TopicsGrpcConfiguration{client_timeout=%v, keepAlivePermitWithoutCalls=%v, keepAliveTimeout=%v, keepAliveTime=%v, maxSendMessageLength=%v, maxReceiveMessageLength=%v, numStreamGrpcChannels=%v, numUnaryGrpcChannels=%v}",
		s.client_timeout, s.keepAlivePermitWithoutCalls, s.keepAliveTimeout, s.keepAliveTime, s.maxSendMessageLength, s.maxReceiveMessageLength, s.numStreamGrpcChannels, s.numUnaryGrpcChannels,
	)
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
