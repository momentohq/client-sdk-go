package momentoerrors

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MomentoSvcErr interface {
	error
	Code() string
	Message() string
	OriginalErr() error
}

// NewMomentoSvcErr returns a new Momento service error.
// Used internally, mostly for invalid argument errors.
// Default to using the code as the message wrapper otherwise.
func NewMomentoSvcErr(code string, message string, originalErr error) MomentoSvcErr {
	if code == InvalidArgumentError {
		return newMomentoSvcErr(code, message, originalErr, "Invalid argument passed to Momento client")
	}
	return newMomentoSvcErr(code, message, originalErr, code)
}

const (
	// InvalidArgumentError occurs when an invalid argument is passed to Momento client.
	InvalidArgumentError = "InvalidArgumentError"
	// InternalServerError occurs when an unexpected error is encountered trying to fulfill the request.
	InternalServerError = "InternalServerError"
	// ClientSdkError occurs when a client side error happens.
	ClientSdkError = "ClientSdkError"
	// BadRequestError occurs when a request was invalid.
	BadRequestError = "BadRequestError"
	// CanceledError occurs when a request was cancelled by the server.
	CanceledError = "CanceledError"
	// TimeoutError occurs when an operation did not complete in time.
	TimeoutError = "TimeoutError"
	// PermissionError occurs when there are insufficient permissions to perform operation.
	PermissionError = "PermissionError"
	// AuthenticationError occurs when invalid authentication credentials to connect to cache service are provided.
	AuthenticationError = "AuthenticationError"
	// LimitExceededError occurs when request rate, bandwidth, or object size exceeded the limits for the account.
	LimitExceededError = "LimitExceededError"
	// NotFoundError occurs when a cache with specified name doesn't exist.
	//
	// Deprecated: Use more specific CacheNotFoundError, StoreNotFoundError, or ItemNotFoundError instead.
	NotFoundError = "NotFoundError"
	// CacheNotFoundError occurs when a cache with specified name doesn't exist.
	CacheNotFoundError = "NotFoundError"
	// StoreNotFoundError occurs when a store with specified name doesn't exist.
	StoreNotFoundError = "StoreNotFoundError"
	// ItemNotFoundError occurs when an item with specified key doesn't exist.
	ItemNotFoundError = "ItemNotFoundError"
	// AlreadyExistsError occurs when a cache with specified name already exists.
	//
	// Deprecated: Use more specific CacheAlreadyExistsError or StoreAlreadyExistsError instead.
	AlreadyExistsError = "AlreadyExistsError"
	// CacheAlreadyExistsError occurs when a cache with specified name already exists.
	CacheAlreadyExistsError = "AlreadyExistsError"
	// StoreAlreadyExistsError occurs when a store with specified name already exists.
	StoreAlreadyExistsError = "StoreAlreadyExistsError"
	// UnknownServiceError occurs when an unknown error has occurred.
	UnknownServiceError = "UnknownServiceError"
	// ServerUnavailableError occurs when the server was unable to handle the request.
	ServerUnavailableError = "ServerUnavailableError"
	// FailedPreconditionError occurs when the system is not in a state required for the operation's execution.
	FailedPreconditionError = "FailedPreconditionError"
	// InternalServerErrorMessage is the message for an unexpected error occurring while trying to fulfill the request.
	InternalServerErrorMessage = "CacheService failed with an internal error"
	// ClientSdkErrorMessage is the message for when SDK Failed to process the request.
	ClientSdkErrorMessage = "SDK Failed to process the request."
	// ConnectionError occurs when there is an error connecting to Momento servers.
	ConnectionError = "ConnectionError"
)

const (
	TopicSubscriptionsLimitExceeded = "Topic subscriptions limit exceeded for this account"
	OperationsRateLimitExceeded     = "Request rate limit exceeded for this account"
	ThroughputRateLimitExceeded     = "Bandwidth limit exceeded for this account"
	RequestSizeLimitExceeded        = "Request size limit exceeded for this account"
	ItemSizeLimitExceeded           = "Item size limit exceeded for this account"
	ElementSizeLimitExceeded        = "Element size limit exceeded for this account"
	UnknownLimitExceeded            = "Limit exceeded for this account"
)

// ConvertSvcErr converts gRPC error to MomentoSvcErr.
func ConvertSvcErr(err error, metadata ...metadata.MD) MomentoSvcErr {
	if grpcStatus, ok := status.FromError(err); ok {
		switch grpcStatus.Code() {
		case codes.InvalidArgument:
			return newMomentoSvcErr(
				InvalidArgumentError,
				grpcStatus.Message(),
				err,
				"Invalid argument passed to Momento client",
			)
		case codes.Unimplemented:
			return newMomentoSvcErr(
				BadRequestError,
				grpcStatus.Message(),
				err,
				"The request was invalid; please contact us at support@momentohq.com",
			)
		case codes.OutOfRange:
			return newMomentoSvcErr(
				BadRequestError,
				grpcStatus.Message(),
				err,
				"The request was invalid; please contact us at support@momentohq.com",
			)
		case codes.FailedPrecondition:
			return newMomentoSvcErr(
				FailedPreconditionError,
				grpcStatus.Message(),
				err,
				"System is not in a state required for the operation's execution",
			)
		case codes.Canceled:
			return newMomentoSvcErr(
				CanceledError,
				grpcStatus.Message(),
				err,
				"The request was cancelled by the server; please contact us at support@momentohq.com",
			)
		case codes.DeadlineExceeded:
			return newMomentoSvcErr(
				TimeoutError,
				grpcStatus.Message(),
				err,
				"The client's configured timeout was exceeded; you may need to use a Configuration with more lenient timeouts",
			)
		case codes.PermissionDenied:
			return newMomentoSvcErr(
				PermissionError,
				grpcStatus.Message(),
				err,
				"Insufficient permissions to perform operation",
			)
		case codes.Unauthenticated:
			return newMomentoSvcErr(
				AuthenticationError,
				grpcStatus.Message(),
				err,
				"Invalid authentication credentials to connect to Momento service",
			)
		case codes.ResourceExhausted:
			// By default, use the generic limit exceeded message wrapper.
			messageWrapper := UnknownLimitExceeded

			// If available, use metadata to determine cause of resource exhausted error.
			if len(metadata) > 1 {
				errData, ok := metadata[1]["err"]
				if ok && errData[0] != "" {
					switch errData[0] {
					case "topic_subscriptions_limit_exceeded":
						messageWrapper = TopicSubscriptionsLimitExceeded
					case "operations_rate_limit_exceeded":
						messageWrapper = OperationsRateLimitExceeded
					case "throughput_rate_limit_exceeded":
						messageWrapper = ThroughputRateLimitExceeded
					case "request_size_limit_exceeded":
						messageWrapper = RequestSizeLimitExceeded
					case "item_size_limit_exceeded":
						messageWrapper = ItemSizeLimitExceeded
					case "element_size_limit_exceeded":
						messageWrapper = ElementSizeLimitExceeded
					}
				} else {
					// If err metadata is not available, try string matching on the
					// error details to return the most specific message wrapper.
					lowerCasedMessage := strings.ToLower(grpcStatus.Message())
					if strings.Contains(lowerCasedMessage, "subscribers") {
						messageWrapper = TopicSubscriptionsLimitExceeded
					} else if strings.Contains(lowerCasedMessage, "operations") {
						messageWrapper = OperationsRateLimitExceeded
					} else if strings.Contains(lowerCasedMessage, "throughput") {
						messageWrapper = ThroughputRateLimitExceeded
					} else if strings.Contains(lowerCasedMessage, "request limit") {
						messageWrapper = RequestSizeLimitExceeded
					} else if strings.Contains(lowerCasedMessage, "item size") {
						messageWrapper = ItemSizeLimitExceeded
					} else if strings.Contains(lowerCasedMessage, "element size") {
						messageWrapper = ElementSizeLimitExceeded
					}
				}
			}
			return newMomentoSvcErr(LimitExceededError, grpcStatus.Message(), err, messageWrapper)
		case codes.NotFound:
			cacheMessageWrapper := "A cache with the specified name does not exist.  To resolve this error, make sure you have created the cache before attempting to use it"
			storeMessageWrapper := "A store with the specified name does not exist.  To resolve this error, make sure you have created the store before attempting to use it"
			// Use metadata to determine cause of not found error
			if len(metadata) > 1 {
				errData, ok := metadata[1]["err"]
				// In the absence of error metadata, return CacheNotFoundError.
				// This is currently re-mapped to a StoreNotFoundError in the PreviewStorageClient"s
				// DeleteStore method.
				if !ok {
					return newMomentoSvcErr(CacheNotFoundError, grpcStatus.Message(), err, cacheMessageWrapper)
				}
				errCause := errData[0]
				if errCause == "item_not_found" {
					return newMomentoSvcErr(
						ItemNotFoundError,
						grpcStatus.Message(),
						err,
						"An item with the specified key does not exist",
					)
				} else if errCause == "store_not_found" {
					return newMomentoSvcErr(StoreNotFoundError, grpcStatus.Message(), err, storeMessageWrapper)
				}
			}
			return newMomentoSvcErr(CacheNotFoundError, grpcStatus.Message(), err, cacheMessageWrapper)
		case codes.AlreadyExists:
			cacheMessageWrapper := "A cache with the specified name already exists.  To resolve this error, either delete the existing cache and make a new one, or use a different name"
			storeMessageWrapper := "A store with the specified name already exists.  To resolve this error, either delete the existing store and make a new one, or use a different name"
			if len(metadata) > 1 {
				errData, ok := metadata[1]["err"]
				// In the absence of error metadata, return CacheAlreadyExistsError.
				if !ok {
					return newMomentoSvcErr(CacheAlreadyExistsError, grpcStatus.Message(), err, cacheMessageWrapper)
				}
				errCause := errData[0]
				switch errCause {
				case "store_already_exists":
					return newMomentoSvcErr(StoreAlreadyExistsError, grpcStatus.Message(), err, storeMessageWrapper)
				default:
					return newMomentoSvcErr(CacheAlreadyExistsError, grpcStatus.Message(), err, cacheMessageWrapper)
				}
			}
			// If no metadata is available, return CacheAlreadyExistsError by default.
			return newMomentoSvcErr(CacheAlreadyExistsError, grpcStatus.Message(), err, cacheMessageWrapper)
		case codes.Unknown:
			return newMomentoSvcErr(
				UnknownServiceError,
				grpcStatus.Message(),
				err,
				"Service returned an unknown response; please contact us at support@momentohq.com",
			)
		case codes.Aborted:
			return newMomentoSvcErr(
				InternalServerError,
				grpcStatus.Message(),
				err,
				"Unexpected error encountered while trying to fulfill the request; please contact us at support@momentohq.com",
			)
		case codes.Internal:
			return newMomentoSvcErr(
				InternalServerError,
				grpcStatus.Message(),
				err,
				"Unexpected error encountered while trying to fulfill the request; please contact us at support@momentohq.com",
			)
		case codes.Unavailable:
			return newMomentoSvcErr(
				ServerUnavailableError,
				grpcStatus.Message(),
				err,
				"The server was unable to handle the request; consider retrying.  If the error persists, please contact us at support@momentohq.com",
			)
		case codes.DataLoss:
			return newMomentoSvcErr(
				InternalServerError,
				grpcStatus.Message(),
				err,
				"Unexpected error encountered while trying to fulfill the request; please contact us at support@momentohq.com",
			)
		default:
			return newMomentoSvcErr(
				UnknownServiceError,
				InternalServerErrorMessage,
				err,
				"Service returned an unknown response; please contact us at support@momentohq.com",
			)
		}
	}
	return NewMomentoSvcErr(ClientSdkError, ClientSdkErrorMessage, err)
}

func NewConnectionError(err error) MomentoSvcErr {
	return NewMomentoSvcErr(ConnectionError, "Connection is in an unexpected state", err)
}
