package config

import "time"

type Laptop struct {
	Configuration
}

const defaultMaxSessionMemoryMb = 256

// 4 minutes.  We want to remain comfortably underneath the idle timeout for AWS NLB, which is 350s.
const defaultMaxIdle = 4 * time.Minute
const defaultLoggerName = "message"

func LatestLaptopConfig() *Laptop {
	return &Laptop{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			Logger: NewLogger(&LoggerConfiguration{
				Name: defaultLoggerName,
			}),
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline:           5 * time.Second,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdle: defaultMaxIdle,
			}),
		}),
	}
}

type InRegion struct {
	Configuration
}

func LatestInRegionConfig() *InRegion {
	return &InRegion{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			Logger: NewLogger(&LoggerConfiguration{
				Name: defaultLoggerName,
			}),
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadline:           1100 * time.Millisecond,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdle: defaultMaxIdle,
			}),
		}),
	}
}
