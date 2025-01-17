package config

import (
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/internal/retry"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type ReadConcern string

const (
	// BALANCED is the default read concern for the cache client.
	BALANCED ReadConcern = "balanced"
	// CONSISTENT read concern guarantees read after write consistency.
	CONSISTENT ReadConcern = "consistent"
)

type ConfigurationProps struct {
	// LoggerFactory represents a type used to configure the Momento logging system.
	LoggerFactory logger.MomentoLoggerFactory
	// TransportStrategy is responsible for configuring network tunables.
	TransportStrategy TransportStrategy
	// RetryStrategy defines a contract for how and when to retry a request.
	RetryStrategy retry.Strategy
	// NumGrpcChannels is the number of GRPC channels the client should open and work with.
	NumGrpcChannels uint32
	// ReadConcern is the read concern for the cache client.
	ReadConcern ReadConcern
}

type Configuration interface {
	// GetLoggerFactory Returns the current configuration options for logging verbosity and format
	GetLoggerFactory() logger.MomentoLoggerFactory

	// GetRetryStrategy Returns the current configuration options for wire interactions with the Momento service
	GetRetryStrategy() retry.Strategy

	// GetTransportStrategy Returns the current configuration options for wire interactions with the Momento service
	GetTransportStrategy() TransportStrategy

	// GetClientSideTimeout Returns the current configuration options for client side timeout with the Momento service
	GetClientSideTimeout() time.Duration

	// GetNumGrpcChannels Returns the configuration option for the number of GRPC channels
	// the cache client should open and work with.
	GetNumGrpcChannels() uint32

	GetReadConcern() ReadConcern

	// WithRetryStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	// with the specified momento.TransportStrategy
	WithRetryStrategy(retryStrategy retry.Strategy) Configuration

	// WithClientTimeout Copy constructor for overriding TransportStrategy client side timeout. Returns a new
	// Configuration object with the specified momento.TransportStrategy using passed client side timeout.
	WithClientTimeout(clientTimeout time.Duration) Configuration

	// WithTransportStrategy Copy constructor for overriding TransportStrategy returns a new Configuration object
	// with the specified momento.TransportStrategy
	WithTransportStrategy(transportStrategy TransportStrategy) Configuration

	// WithNumGrpcChannels Copy constructor for overriding NumGrpcChannels returns a new Configuration object
	// with the specified NumGrpcChannels
	WithNumGrpcChannels(numGrpcChannels uint32) Configuration

	// WithReadConcern Copy constructor for overriding ReadConcern returns a new Configuration object
	// with the specified ReadConcern
	WithReadConcern(readConcern ReadConcern) Configuration
}

type cacheConfiguration struct {
	loggerFactory     logger.MomentoLoggerFactory
	transportStrategy TransportStrategy
	retryStrategy     retry.Strategy
	numGrpcChannels   uint32
	readConcern       ReadConcern
}

func (s *cacheConfiguration) GetLoggerFactory() logger.MomentoLoggerFactory {
	return s.loggerFactory
}

func (s *cacheConfiguration) GetClientSideTimeout() time.Duration {
	return s.transportStrategy.GetClientSideTimeout()
}

func NewCacheConfiguration(props *ConfigurationProps) Configuration {
	return &cacheConfiguration{
		loggerFactory:     props.LoggerFactory,
		transportStrategy: props.TransportStrategy,
		retryStrategy:     props.RetryStrategy,
		numGrpcChannels:   props.NumGrpcChannels,
		readConcern:       props.ReadConcern,
	}
}

func (s *cacheConfiguration) GetTransportStrategy() TransportStrategy {
	return s.transportStrategy
}

func (s *cacheConfiguration) GetRetryStrategy() retry.Strategy {
	return s.retryStrategy
}

func (s *cacheConfiguration) GetNumGrpcChannels() uint32 {
	return s.numGrpcChannels
}

func (s *cacheConfiguration) GetReadConcern() ReadConcern {
	return s.readConcern
}

func (s *cacheConfiguration) WithClientTimeout(clientTimeout time.Duration) Configuration {
	return &cacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy.WithClientTimeout(clientTimeout),
		retryStrategy:     s.retryStrategy,
		numGrpcChannels:   s.numGrpcChannels,
		readConcern:       s.readConcern,
	}
}

func (s *cacheConfiguration) WithRetryStrategy(strategy retry.Strategy) Configuration {
	return &cacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		retryStrategy:     strategy,
		numGrpcChannels:   s.numGrpcChannels,
		readConcern:       s.readConcern,
	}
}

func (s *cacheConfiguration) WithTransportStrategy(transportStrategy TransportStrategy) Configuration {
	return &cacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: transportStrategy,
		retryStrategy:     s.retryStrategy,
		numGrpcChannels:   s.numGrpcChannels,
		readConcern:       s.readConcern,
	}
}

func (s *cacheConfiguration) WithNumGrpcChannels(numGrpcChannels uint32) Configuration {
	return &cacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		retryStrategy:     s.retryStrategy,
		numGrpcChannels:   numGrpcChannels,
		readConcern:       s.readConcern,
	}
}

func (s *cacheConfiguration) WithReadConcern(readConcern ReadConcern) Configuration {
	return &cacheConfiguration{
		loggerFactory:     s.loggerFactory,
		transportStrategy: s.transportStrategy,
		retryStrategy:     s.retryStrategy,
		numGrpcChannels:   s.numGrpcChannels,
		readConcern:       readConcern,
	}
}

func (s *cacheConfiguration) String() string {
	return fmt.Sprintf(
		"Configuration{loggerFactory=%v, transportStrategy=%v, retryStrategy=%v, numGrpcChannels=%v, readConcern=%v}",
		s.loggerFactory,
		s.transportStrategy,
		s.retryStrategy,
		s.numGrpcChannels,
		s.readConcern,
	)
}
