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

var retryableRequestMethods = map[string]bool{
	string(momento_rpc_names.Get):      true,
	string(momento_rpc_names.GetBatch): true,
	string(momento_rpc_names.Set):      true,
	string(momento_rpc_names.SetBatch): true,
	string(momento_rpc_names.SetIf):    false,
	// SetIfNotExists is deprecated
	string(momento_rpc_names.SetIfNotExists): false,
	string(momento_rpc_names.Delete):         true,
	string(momento_rpc_names.KeysExist):      true,
	string(momento_rpc_names.Increment):      false,
	string(momento_rpc_names.UpdateTtl):      false,
	string(momento_rpc_names.ItemGetTtl):     true,
	string(momento_rpc_names.ItemGetType):    true,

	string(momento_rpc_names.DictionaryGet):       true,
	string(momento_rpc_names.DictionaryFetch):     true,
	string(momento_rpc_names.DictionarySet):       true,
	string(momento_rpc_names.DictionaryIncrement): false,
	string(momento_rpc_names.DictionaryDelete):    true,
	string(momento_rpc_names.DictionaryLength):    true,

	string(momento_rpc_names.SetFetch):      true,
	string(momento_rpc_names.SetSample):     true,
	string(momento_rpc_names.SetUnion):      true,
	string(momento_rpc_names.SetDifference): true,
	string(momento_rpc_names.SetContains):   true,
	string(momento_rpc_names.SetLength):     true,
	string(momento_rpc_names.SetPop):        false,

	string(momento_rpc_names.ListPushFront): false,
	string(momento_rpc_names.ListPushBack):  false,
	string(momento_rpc_names.ListPopFront):  false,
	string(momento_rpc_names.ListPopBack):   false,
	// Not used, and unknown "/cache_client.Scs/ListErase",
	string(momento_rpc_names.ListRemove):           true,
	string(momento_rpc_names.ListFetch):            true,
	string(momento_rpc_names.ListLength):           true,
	string(momento_rpc_names.ListConcatenateFront): false,
	string(momento_rpc_names.ListConcatenateBack):  false,
	string(momento_rpc_names.ListRetain):           false,

	string(momento_rpc_names.SortedSetPut):           true,
	string(momento_rpc_names.SortedSetFetch):         true,
	string(momento_rpc_names.SortedSetGetScore):      true,
	string(momento_rpc_names.SortedSetRemove):        true,
	string(momento_rpc_names.SortedSetIncrement):     false,
	string(momento_rpc_names.SortedSetGetRank):       true,
	string(momento_rpc_names.SortedSetLength):        true,
	string(momento_rpc_names.SortedSetLengthByScore): true,

	string(momento_rpc_names.TopicSubscribe): true,
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
