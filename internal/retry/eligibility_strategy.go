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
	"/cache_client.Scs/Get":    true,
	// NEW and idempotent "/cache_client.Scs/GetBatch"
	"/cache_client.Scs/Set":    true,
	// NEW and idempotent "/cache_client.Scs/SetBatch"
	// NEW and _NOT_ idempotent "/cache_client.Scs/SetIf"

	// deprecated and not idempotent "/cache_client.Scs/SetIfNotExists"
	"/cache_client.Scs/Delete": true,
	// NEW and idempotent "/cache_client.Scs/KeysExist"
	// not idempotent "/cache_client.Scs/Increment"
	// NEW and server-idempotent but return values can be different "/cache_client.Scs/UpdateTtl"
	// NEW and idempotent "/cache_client.Scs/ItemGetTtl"
	// NEW and idempotent "/cache_client.Scs/ItemGetType"

	"/cache_client.Scs/DictionaryGet":    true,
	"/cache_client.Scs/DictionaryFetch":  true,
	"/cache_client.Scs/DictionarySet": true,
	// not idempotent: "/cache_client.Scs/DictionaryIncrement",
	"/cache_client.Scs/DictionaryDelete": true,
	// NEW and idempotent "/cache_client.Scs/DictionaryLength"

	"/cache_client.Scs/SetFetch":         true,
	// NEW and idempotent "/cache_client.Scs/SetSample"
	"/cache_client.Scs/SetUnion":         true,
	"/cache_client.Scs/SetDifference":    true,
	// NEW and idempotent "/cache_client.Scs/SetContains"
	// NEW and idempotent "/cache_client.Scs/SetLength"
	// NEW and _NOT_ idempotent "/cache_client.Scs/SetPop"

	// not idempotent: "/cache_client.Scs/ListPushFront",
	// not idempotent: "/cache_client.Scs/ListPushBack",
	// not idempotent: "/cache_client.Scs/ListPopFront",
	// not idempotent: "/cache_client.Scs/ListPopBack",
	// NEW, not used, and unknown "/cache_client.Scs/ListErase",
	// Warning: in the future, this may not be idempotent
	// Currently it supports removing all occurrences of a value.
	// In the future, we may also add "the first/last N occurrences of a value".
	// In the latter case it is not idempotent.
	"/cache_client.Scs/ListRemove": true,
	"/cache_client.Scs/ListFetch": true,
	"/cache_client.Scs/ListLength": true,
	// not idempotent: "/cache_client.Scs/ListConcatenateFront",
	// not idempotent: "/cache_client.Scs/ListConcatenateBack"
	// NEW and _NOT_ idempotent "/cache_client.Scs/ListRetain"

	// NEW and idempotent "/cache_client.Scs/SortedSetPut"
	// NEW and idempotent "/cache_client.Scs/SortedSetFetch"
	// NEW and idempotent "/cache_client.Scs/SortedSetGetScore"
	// NEW and idempotent "/cache_client.Scs/SortedSetRemove"
	// NEW and _NOT_ idempotent "/cache_client.Scs/SortedSetIncrement"
	// NEW and idempotent "/cache_client.Scs/SortedSetGetRank"
	// NEW and idempotent "/cache_client.Scs/SortedSetLength"
	// NEW and idempotent "/cache_client.Scs/SortedSetLengthByScore"

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
