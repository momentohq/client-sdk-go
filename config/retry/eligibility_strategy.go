package retry

import (
	"github.com/momentohq/client-sdk-go/momento_rpc_names"
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

var retryableRequestMethods = map[momento_rpc_names.MomentoRPCMethod]bool{
	momento_rpc_names.Get:      true,
	momento_rpc_names.GetBatch: true,
	momento_rpc_names.Set:      true,
	momento_rpc_names.SetBatch: true,
	momento_rpc_names.SetIf:    false,
	// SetIfNotExists is deprecated
	momento_rpc_names.SetIfNotExists: false,
	momento_rpc_names.Delete:         true,
	momento_rpc_names.KeysExist:      true,
	momento_rpc_names.Increment:      false,
	momento_rpc_names.UpdateTtl:      false,
	momento_rpc_names.ItemGetTtl:     true,
	momento_rpc_names.ItemGetType:    true,

	momento_rpc_names.DictionaryGet:       true,
	momento_rpc_names.DictionaryFetch:     true,
	momento_rpc_names.DictionarySet:       true,
	momento_rpc_names.DictionaryIncrement: false,
	momento_rpc_names.DictionaryDelete:    true,
	momento_rpc_names.DictionaryLength:    true,

	momento_rpc_names.SetFetch:      true,
	momento_rpc_names.SetSample:     true,
	momento_rpc_names.SetUnion:      true,
	momento_rpc_names.SetDifference: true,
	momento_rpc_names.SetContains:   true,
	momento_rpc_names.SetLength:     true,
	momento_rpc_names.SetPop:        false,

	momento_rpc_names.ListPushFront: false,
	momento_rpc_names.ListPushBack:  false,
	momento_rpc_names.ListPopFront:  false,
	momento_rpc_names.ListPopBack:   false,
	// Not used, and unknown "/cache_client.Scs/ListErase",
	momento_rpc_names.ListRemove:           true,
	momento_rpc_names.ListFetch:            true,
	momento_rpc_names.ListLength:           true,
	momento_rpc_names.ListConcatenateFront: false,
	momento_rpc_names.ListConcatenateBack:  false,
	momento_rpc_names.ListRetain:           false,

	momento_rpc_names.SortedSetPut:           true,
	momento_rpc_names.SortedSetFetch:         true,
	momento_rpc_names.SortedSetGetScore:      true,
	momento_rpc_names.SortedSetRemove:        true,
	momento_rpc_names.SortedSetIncrement:     false,
	momento_rpc_names.SortedSetGetRank:       true,
	momento_rpc_names.SortedSetLength:        true,
	momento_rpc_names.SortedSetLengthByScore: true,

	momento_rpc_names.TopicSubscribe: true,
}

type DefaultEligibilityStrategy struct{}

func (s DefaultEligibilityStrategy) IsEligibleForRetry(props StrategyProps) bool {
	if !retryableStatusCodes[props.GrpcStatusCode] {
		return false
	}

	if !retryableRequestMethods[momento_rpc_names.MomentoRPCMethod(props.GrpcMethod)] {
		return false
	}
	return true
}
