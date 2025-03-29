package helpers

import "fmt"

type timestampPayload struct {
	cacheName   string
	requestName string
	timestamp   int64
}

type retryMetrics struct {
	data map[string]map[string][]int64
}

type RetryMetricsCollector interface {
	AddTimestamp(cacheName string, requestName string, timestamp int64)
	GetTotalRetryCount(cacheName string, requestName string) (int, error)
	GetAverageTimeBetweenRetries(cacheName string, requestName string) (int64, error)
	GetAllMetrics() map[string]map[string][]int64
}

func NewRetryMetricsCollector() RetryMetricsCollector {
	return &retryMetrics{data: make(map[string]map[string][]int64)}
}

func (r *retryMetrics) AddTimestamp(cacheName string, requestName string, timestamp int64) {
	if _, ok := r.data[cacheName]; !ok {
		r.data[cacheName] = make(map[string][]int64)
	}
	r.data[cacheName][requestName] = append(r.data[cacheName][requestName], timestamp)
}

func (r *retryMetrics) GetTotalRetryCount(cacheName string, requestName string) (int, error) {
	if _, ok := r.data[cacheName]; !ok {
		return 0, fmt.Errorf("cache name '%s' is not valid", cacheName)
	}
	if timestamps, ok := r.data[cacheName][requestName]; ok {
		// The first timestamp is the original request, so we subtract 1
		return len(timestamps) - 1, nil
	}
	return 0, fmt.Errorf("request name '%s' is not valid", requestName)
}

// GetAverageTimeBetweenRetries returns the average time between retries in seconds.
//
//	Limited to second resolution, but I can obviously change that if desired.
//	This tracks with the JS implementation.
func (r *retryMetrics) GetAverageTimeBetweenRetries(cacheName string, requestName string) (int64, error) {
	if _, ok := r.data[cacheName]; !ok {
		return int64(0), fmt.Errorf("cache name '%s' is not valid", cacheName)
	}
	if timestamps, ok := r.data[cacheName][requestName]; ok {
		if len(timestamps) < 2 {
			return 0, nil
		}
		var sum int64
		for i := 1; i < len(timestamps); i++ {
			sum += timestamps[i] - timestamps[i-1]
		}
		return sum / int64(len(timestamps)-1), nil
	}
	return 0, fmt.Errorf("request name '%s' is not valid", requestName)
}

func (r *retryMetrics) GetAllMetrics() map[string]map[string][]int64 {
	return r.data
}
