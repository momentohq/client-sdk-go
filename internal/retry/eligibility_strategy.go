package retry

import "google.golang.org/grpc/codes"

type EligibilityStrategy interface {
	// IsEligibleForRetry Determines whether a grpc call is able to be retried. The determination is based on the result
	// of the last invocation of the call and whether the call is idempotent.
	IsEligibleForRetry(props StrategyProps) bool
}

var retryableStatusCodes = map[codes.Code]bool{
	codes.Internal:    true,
	codes.Unavailable: true,
}

var retryableRequestMethods = map[string]bool{
	"/cache_client.Scs/Set":    true,
	"/cache_client.Scs/Get":    true,
	"/cache_client.Scs/Delete": true,
	// not idempotent "/cache_client.Scs/Increment"
	"/cache_client.Scs/DictionarySet": true,
	// not idempotent: "/cache_client.Scs/DictionaryIncrement",
	"/cache_client.Scs/DictionaryGet":    true,
	"/cache_client.Scs/DictionaryFetch":  true,
	"/cache_client.Scs/DictionaryDelete": true,
	"/cache_client.Scs/SetUnion":         true,
	"/cache_client.Scs/SetDifference":    true,
	"/cache_client.Scs/SetFetch":         true,
	// not idempotent: "/cache_client.Scs/SetIfNotExists"
	// not idempotent: "/cache_client.Scs/ListPushFront",
	// not idempotent: "/cache_client.Scs/ListPushBack",
	// not idempotent: "/cache_client.Scs/ListPopFront",
	// not idempotent: "/cache_client.Scs/ListPopBack",
	"/cache_client.Scs/ListFetch": true,
	// Warning: in the future, this may not be idempotent
	// Currently it supports removing all occurrences of a value.
	// In the future, we may also add "the first/last N occurrences of a value".
	// In the latter case it is not idempotent.
	"/cache_client.Scs/ListRemove": true,
	"/cache_client.Scs/ListLength": true,
	// not idempotent: "/cache_client.Scs/ListConcatenateFront",
	// not idempotent: "/cache_client.Scs/ListConcatenateBack"
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
