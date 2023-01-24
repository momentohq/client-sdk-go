package config

type Laptop struct {
	Configuration
}

const defaultMaxSessionMemoryMb = 256

// 4 minutes.  We want to remain comfortably underneath the idle timeout for AWS NLB, which is 350s.
const defaultMaxIdleMillis = 4 * 60 * 1_000

func LatestLaptopConfig() *Laptop {
	return &Laptop{
		Configuration: NewSimpleCacheConfiguration(&ConfigurationProps{
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadlineMillis:     5000,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdleMillis: defaultMaxIdleMillis,
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
			TransportStrategy: NewStaticTransportStrategy(&TransportStrategyProps{
				GrpcConfiguration: NewStaticGrpcConfiguration(&GrpcConfigurationProps{
					deadlineMillis:     1100,
					maxSessionMemoryMb: defaultMaxSessionMemoryMb,
				}),
				MaxIdleMillis: defaultMaxIdleMillis,
			}),
		}),
	}
}
