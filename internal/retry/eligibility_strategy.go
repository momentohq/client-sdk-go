package retry

import "google.golang.org/grpc/codes"

type EligibilityStrategy interface {
	// IsEligibleForRetry Determines whether a grpc call is able to be retried. The determination is based on the result
	// of the last invocation of the call and whether the call is idempotent.
	IsEligibleForRetry(props StrategyProps) bool
}

var retryableStatusCodes = map[codes.Code]bool{
	codes.Canceled:    true,
	codes.Internal:    true,
	codes.Unavailable: true,
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
	// UpdateTtl is idempotent on the server but return values can be different if the first call is successful.
	"/cache_client.Scs/UpdateTtl":   true,
	"/cache_client.Scs/ItemGetTtl":  true,
	"/cache_client.Scs/ItemGetType": true,

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
	// Warning: in the future, ListRemove may not be idempotent
	// Currently it supports removing all occurrences of a value.
	// In the future, we may also add "the first/last N occurrences of a value".
	// In the latter case it is not idempotent.
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
