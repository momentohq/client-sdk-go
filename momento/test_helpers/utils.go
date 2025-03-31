package helpers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/momento_rpc_names"
)

func NewRandomString() string {
	return uuid.NewString()
}

func NewRandomMomentoString() momento.String {
	return momento.String(NewRandomString())
}

func NewRandomMomentoBytes() momento.Bytes {
	return momento.Bytes([]byte(uuid.NewString()))
}

func ConvertErrorCodeToMomentoLocalErrorCode(code string) string {
	switch code {
	case momentoerrors.InvalidArgumentError:
		return "invalid-argument"
	case momentoerrors.CacheAlreadyExistsError:
		return "already-exists"
	case momentoerrors.StoreAlreadyExistsError:
		return "already-exists"
	case momentoerrors.CacheNotFoundError:
		return "not-found"
	case momentoerrors.StoreNotFoundError:
		return "not-found"
	case momentoerrors.ItemNotFoundError:
		return "not-found"
	case momentoerrors.InternalServerError:
		return "internal"
	case momentoerrors.PermissionError:
		return "permission-denied"
	case momentoerrors.AuthenticationError:
		return "unauthenticated"
	case momentoerrors.CanceledError:
		return "cancelled"
	case momentoerrors.ConnectionError:
		return "unavailable"
	case momentoerrors.LimitExceededError:
		return "resource-exhausted"
	case momentoerrors.BadRequestError:
		return "invalid-argument"
	case momentoerrors.TimeoutError:
		return "deadline-exceeded"
	case momentoerrors.ServerUnavailableError:
		return "unavailable"
	case momentoerrors.FailedPreconditionError:
		return "failed-precondition"
	case momentoerrors.UnknownServiceError:
		return "unknown"
	default:
		// This is only used for testing, so we can panic here.
		panic(fmt.Sprintf("Unknown error code: %s", code))
	}
}

func ConvertRpcNameToMomentoLocalRpcName(rpcName momento_rpc_names.MomentoRPCMethod) string {
	switch rpcName {
	case momento_rpc_names.Get:
		return "get"
	case momento_rpc_names.GetWithHash:
		return "get-with-hash"
	case momento_rpc_names.Set:
		return "set"
	case momento_rpc_names.SetIfHash:
		return "set-if-hash"
	case momento_rpc_names.Delete:
		return "delete"
	case momento_rpc_names.Increment:
		return "increment"
	case momento_rpc_names.SetIf:
		return "set-if"
	case momento_rpc_names.GetBatch:
		return "get-batch"
	case momento_rpc_names.SetBatch:
		return "set-batch"
	case momento_rpc_names.KeysExist:
		return "keys-exist"
	case momento_rpc_names.UpdateTtl:
		return "update-ttl"
	case momento_rpc_names.ItemGetTtl:
		return "item-get-ttl"
	case momento_rpc_names.ItemGetType:
		return "item-get-type"
	case momento_rpc_names.DictionarySet:
		return "dictionary-set"
	case momento_rpc_names.DictionaryGet:
		return "dictionary-get"
	case momento_rpc_names.DictionaryFetch:
		return "dictionary-fetch"
	case momento_rpc_names.DictionaryIncrement:
		return "dictionary-increment"
	case momento_rpc_names.DictionaryDelete:
		return "dictionary-delete"
	case momento_rpc_names.DictionaryLength:
		return "dictionary-length"
	case momento_rpc_names.SetFetch:
		return "set-fetch"
	case momento_rpc_names.SetSample:
		return "set-sample"
	case momento_rpc_names.SetUnion:
		return "set-union"
	case momento_rpc_names.SetDifference:
		return "set-difference"
	case momento_rpc_names.SetContains:
		return "set-contains"
	case momento_rpc_names.SetLength:
		return "set-length"
	case momento_rpc_names.SetPop:
		return "set-pop"
	case momento_rpc_names.ListPushFront:
		return "list-push-front"
	case momento_rpc_names.ListPushBack:
		return "list-push-back"
	case momento_rpc_names.ListPopFront:
		return "list-pop-front"
	case momento_rpc_names.ListPopBack:
		return "list-pop-back"
	case momento_rpc_names.ListErase:
		return "list-remove"
	case momento_rpc_names.ListFetch:
		return "list-fetch"
	case momento_rpc_names.ListLength:
		return "list-length"
	case momento_rpc_names.ListConcatenateFront:
		return "list-concatenate-front"
	case momento_rpc_names.ListConcatenateBack:
		return "list-concatenate-back"
	case momento_rpc_names.ListRetain:
		return "list-retain"
	case momento_rpc_names.SortedSetPut:
		return "sorted-set-put"
	case momento_rpc_names.SortedSetFetch:
		return "sorted-set-fetch"
	case momento_rpc_names.SortedSetGetScore:
		return "sorted-set-get-score"
	case momento_rpc_names.SortedSetRemove:
		return "sorted-set-remove"
	case momento_rpc_names.SortedSetIncrement:
		return "sorted-set-increment"
	case momento_rpc_names.SortedSetGetRank:
		return "sorted-set-get-rank"
	case momento_rpc_names.SortedSetLength:
		return "sorted-set-length"
	case momento_rpc_names.SortedSetLengthByScore:
		return "sorted-set-length-by-score"
	case momento_rpc_names.TopicPublish:
		return "topic-publish"
	case momento_rpc_names.TopicSubscribe:
		return "topic-subscribe"
	default:
		// This is only used for testing, so we can panic here.
		panic(fmt.Sprintf("Unknown RPC method: %s", rpcName))
	}
}
