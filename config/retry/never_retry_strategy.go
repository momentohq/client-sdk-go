package retry

type neverRetryStrategy struct{}

func (r *neverRetryStrategy) DetermineWhenToRetry(_ StrategyProps) *int {
	return nil
}

func (r *neverRetryStrategy) String() string {
	return "neverRetryStrategy{}"
}

// NewNeverRetryStrategy is a retry strategy that never retries any request
func NewNeverRetryStrategy() Strategy {
	return &neverRetryStrategy{}
}
