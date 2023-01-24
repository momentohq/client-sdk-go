package config

type TransportStrategyProps struct {
	// low-level gRPC settings for communication with the Momento server
	GrpcConfiguration GrpcConfiguration

	// The maximum duration for which a connection may remain idle before being replaced.  This
	// setting can be used to force re-connection of a client if it has been idle for too long.
	// In environments such as AWS lambda, if the lambda is suspended for too long the connection
	// may be closed by the load balancer, resulting in an error on the subsequent request.  If
	// this setting is set to a duration less than the load balancer timeout, we can ensure that
	// the connection will be refreshed to avoid errors.
	MaxIdleMillis uint32
}

type TransportStrategy interface {

	// GetGrpcConfig Configures the low-level gRPC settings for the Momento client's communication
	// with the Momento server.
	GetGrpcConfig() GrpcConfiguration

	// WithGrpcConfig Copy constructor for overriding the gRPC configuration. Returns  a new
	// TransportStrategy with the specified gRPC config.
	WithGrpcConfig(grpcConfig GrpcConfiguration) TransportStrategy

	// GetClientSideTimeout Gets configuration for client side timeout from transport strategy
	GetClientSideTimeout() uint32

	// WithClientTimeoutMillis Copy constructor for overriding the client sie timeout. Returns a new
	// TransportStrategy with the specified client side timeout.
	WithClientTimeoutMillis(clientTimeoutMillis uint32) TransportStrategy

	// WithMaxIdleMillis Copy constructor for overriding the max idle connection timeout. Returns a new
	// TransportStrategy with the specified client side idle connection timeout.
	WithMaxIdleMillis(maxIdleMillis uint32) TransportStrategy
}

type StaticGrpcConfiguration struct {
	deadlineMillis     uint32
	maxSessionMemoryMb uint32
}

func NewStaticGrpcConfiguration(grpcConfiguration *GrpcConfigurationProps) *StaticGrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadlineMillis:     grpcConfiguration.deadlineMillis,
		maxSessionMemoryMb: grpcConfiguration.maxSessionMemoryMb,
	}
}

func (s *StaticGrpcConfiguration) GetDeadlineMillis() uint32 {
	return s.deadlineMillis
}

func (s *StaticGrpcConfiguration) GetMaxSessionMemoryMb() uint32 {
	return s.maxSessionMemoryMb
}

func (s *StaticGrpcConfiguration) WithMaxSessionMb(maxSessionMemoryMb uint32) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadlineMillis:     s.deadlineMillis,
		maxSessionMemoryMb: maxSessionMemoryMb,
	}
}

func (s *StaticGrpcConfiguration) WithDeadlineMillis(deadlineMillis uint32) GrpcConfiguration {
	return &StaticGrpcConfiguration{
		deadlineMillis:     deadlineMillis,
		maxSessionMemoryMb: s.maxSessionMemoryMb,
	}
}

type StaticTransportStrategy struct {
	grpcConfig    GrpcConfiguration
	maxIdleMillis uint32
}

func (s *StaticTransportStrategy) GetClientSideTimeout() uint32 {
	return s.grpcConfig.GetDeadlineMillis()
}

func (s *StaticTransportStrategy) WithClientTimeoutMillis(clientTimeoutMillis uint32) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig:    s.grpcConfig.WithDeadlineMillis(clientTimeoutMillis),
		maxIdleMillis: s.maxIdleMillis,
	}
}

func NewStaticTransportStrategy(props *TransportStrategyProps) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig:    props.GrpcConfiguration,
		maxIdleMillis: props.MaxIdleMillis,
	}
}

func (s *StaticTransportStrategy) GetGrpcConfig() GrpcConfiguration {
	return s.grpcConfig
}

func (s *StaticTransportStrategy) GetMaxIdleMillis() uint32 {
	return s.maxIdleMillis
}

func (s *StaticTransportStrategy) WithGrpcConfig(grpcConfig GrpcConfiguration) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig:    grpcConfig,
		maxIdleMillis: s.maxIdleMillis,
	}
}

func (s *StaticTransportStrategy) WithMaxIdleMillis(maxIdleMillis uint32) TransportStrategy {
	return &StaticTransportStrategy{
		grpcConfig:    s.grpcConfig,
		maxIdleMillis: maxIdleMillis,
	}
}
