package retry

import (
	"google.golang.org/grpc/codes"
)

type EligibilityStrategy interface {
	// IsEligibleForRetry Determines whether a grpc call is able to be retried. The determination is based on the result
	// of the last invocation of the call and whether the call is idempotent.
	IsEligibleForRetry(props StrategyProps) bool
}

var retryableStatusCodes = map[codes.Code]bool{
	codes.Internal:    true,
	codes.Unavailable: true,
	// this code is retryable in other SDKs, but because the client can generate this error code
	// by cancelling the context, we do not retry it here.
	codes.Canceled: false,
}

var retryableRequestMethods = map[string]bool{
	"/cache_client.Scs/Get":      true,
	"/cache_client.Scs/GetBatch": true,
	"/cache_client.Scs/Set":      true,
	"/cache_client.Scs/SetBatch": true,
	"/cache_client.Scs/SetIf":    false,
	// SetIfNotExists is deprecated
	"/cache_client.Scs/SetIfNotExists": false,
	"/cache_client.Scs/Delete":         true,
	"/cache_client.Scs/KeysExist":      true,
	"/cache_client.Scs/Increment":      false,
	"/cache_client.Scs/UpdateTtl":      false,
	"/cache_client.Scs/ItemGetTtl":     true,
	"/cache_client.Scs/ItemGetType":    true,

	"/cache_client.Scs/DictionaryGet":       true,
	"/cache_client.Scs/DictionaryFetch":     true,
	"/cache_client.Scs/DictionarySet":       true,
	"/cache_client.Scs/DictionaryIncrement": false,
	"/cache_client.Scs/DictionaryDelete":    true,
	"/cache_client.Scs/DictionaryLength":    true,

	"/cache_client.Scs/SetFetch":      true,
	"/cache_client.Scs/SetSample":     true,
	"/cache_client.Scs/SetUnion":      true,
	"/cache_client.Scs/SetDifference": true,
	"/cache_client.Scs/SetContains":   true,
	"/cache_client.Scs/SetLength":     true,
	"/cache_client.Scs/SetPop":        false,

	"/cache_client.Scs/ListPushFront": false,
	"/cache_client.Scs/ListPushBack":  false,
	"/cache_client.Scs/ListPopFront":  false,
	"/cache_client.Scs/ListPopBack":   false,
	// Not used, and unknown "/cache_client.Scs/ListErase",
	"/cache_client.Scs/ListRemove":           true,
	"/cache_client.Scs/ListFetch":            true,
	"/cache_client.Scs/ListLength":           true,
	"/cache_client.Scs/ListConcatenateFront": false,
	"/cache_client.Scs/ListConcatenateBack":  false,
	"/cache_client.Scs/ListRetain":           false,

	"/cache_client.Scs/SortedSetPut":           true,
	"/cache_client.Scs/SortedSetFetch":         true,
	"/cache_client.Scs/SortedSetGetScore":      true,
	"/cache_client.Scs/SortedSetRemove":        true,
	"/cache_client.Scs/SortedSetIncrement":     false,
	"/cache_client.Scs/SortedSetGetRank":       true,
	"/cache_client.Scs/SortedSetLength":        true,
	"/cache_client.Scs/SortedSetLengthByScore": true,

	"/cache_client.pubsub.Pubsub/Subscribe": true,
}

// DefaultEligibilityStrategy is the default strategy for determining if a request is eligible for retry.
type DefaultEligibilityStrategy struct{}

func (s DefaultEligibilityStrategy) IsEligibleForRetry(props StrategyProps) bool {
	if !retryableStatusCodes[props.GrpcStatusCode] {
		return false
	}

	if !retryableRequestMethods[props.GrpcMethod] {
		return false
	}
	return true
}

// TimeoutAwareEligibilityStrategy is an eligibility strategy that treats DeadlineExceeded as retryable.
// It is used in conjunction with the TimeoutAwareFixedCountRetryStrategy.
//
// Typically, a DeadlineExceeded would indicate that the overall client timeout has been reached and no further
// retries should be attempted. However, there are some cases where timeouts may occur due to the client being
// overloaded or experiencing transient network issues. In these cases, it may be beneficial to retry the request
// even if the overall timeout has been reached.
type TimeoutAwareEligibilityStrategy struct{}

var timeoutAwareRetryableStatusCodes = map[codes.Code]bool{
	codes.Internal:         true,
	codes.Unavailable:      true,
	codes.DeadlineExceeded: true,
	// this code is retryable in other SDKs, but because the client can generate this error code
	// by cancelling the context, we do not retry it here.
	codes.Canceled: false,
}

func (s TimeoutAwareEligibilityStrategy) IsEligibleForRetry(props StrategyProps) bool {
	return timeoutAwareRetryableStatusCodes[props.GrpcStatusCode] && retryableRequestMethods[props.GrpcMethod]
}
