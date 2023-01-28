package config

import "time"
import "google.golang.org/grpc/codes"

type Laptop struct {
	Configuration
}

const defaultMaxSessionMemoryMb = 256

// 4 minutes.  We want to remain comfortably underneath the idle timeout for AWS NLB, which is 350s.
const defaultMaxIdle = 4 * time.Minute
const defaultMaxRetries = 3

var defaultRetryableStatusCodes = []codes.Code{
	codes.Internal,
	codes.Unavailable,
}

//var defaultRetryStrategy =

func LatestLaptopConfig() *Laptop {
	const overallRequestTimeout = 5 * time.Second
	return &Laptop{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline:           overallRequestTimeout,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdle: defaultMaxIdle,
			}),
			RetryStrategy: NewStaticRetryStrategy(&RetryStrategyProps{
				RetryableRequestStatuses: defaultRetryableStatusCodes,
				MaxRetries:               defaultMaxRetries,
				PerRetryTimeout:          overallRequestTimeout / defaultMaxRetries,
			}),
		}),
	}
}

type InRegion struct {
	Configuration
}

func LatestInRegionConfig() *InRegion {
	const overallRequestTimeout = 1100 * time.Millisecond
	return &InRegion{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline:           overallRequestTimeout,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdle: defaultMaxIdle,
			}),
			RetryStrategy: NewStaticRetryStrategy(&RetryStrategyProps{
				RetryableRequestStatuses: defaultRetryableStatusCodes,
				MaxRetries:               defaultMaxRetries,
				PerRetryTimeout:          overallRequestTimeout / defaultMaxRetries,
			}),
		}),
	}
}
