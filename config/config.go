package config

type ConfigurationProps struct {
	TransportStrategy TransportStrategy
}

type Configuration interface {
	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	//with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration

	// GetClientSideTimeoutMillis Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeoutMillis() uint32

	// WithClientTimeoutMillis Copy constructor for overriding TransportStrategy client side timeout. Returns a new Configuration object
	// with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeoutMillis(clientTimeoutMillis uint32) Configuration
}

type SimpleCacheConfiguration struct {
	transportStrategy TransportStrategy
}

func (s *SimpleCacheConfiguration) GetClientSideTimeoutMillis() uint32 {
	return s.transportStrategy.GetClientSideTimeout()
}

func NewSimpleCacheConfiguration(props *ConfigurationProps) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: props.TransportStrategy,
	}
}

func (s *SimpleCacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *SimpleCacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: transportStrategy,
	}
}

func (s *SimpleCacheConfiguration) WithClientTimeoutMillis(clientTimeoutMillis uint32) Configuration {
	return &SimpleCacheConfiguration{
		transportStrategy: s.transportStrategy.WithClientTimeoutMillis(clientTimeoutMillis),
	}
}
