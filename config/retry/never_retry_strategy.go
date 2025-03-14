package retry

type NeverRetryStrategy struct{}

func (r NeverRetryStrategy) DetermineWhenToRetry(props StrategyProps) *int {
	return nil
}

func (r NeverRetryStrategy) String() string {
	return "NeverRetryStrategy{}"
}

// NewNeverRetryStrategy is a retry strategy that never retries any request
func NewNeverRetryStrategy() Strategy {
	return NeverRetryStrategy{}
}

