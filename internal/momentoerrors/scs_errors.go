package momentoerrors

import (
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
func NewMomentoSvcErr(code string, message string, originalErr error) MomentoSvcErr {
	return newMomentoSvcErr(code, message, originalErr)
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
	AlreadyExistsError = "AlreadyExistsError"
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

// ConvertSvcErr converts gRPC error to MomentoSvcErr.
func ConvertSvcErr(err error, metadata ...metadata.MD) MomentoSvcErr {
	if grpcStatus, ok := status.FromError(err); ok {
		switch grpcStatus.Code() {
		case codes.InvalidArgument:
			return NewMomentoSvcErr(InvalidArgumentError, grpcStatus.Message(), err)
		case codes.Unimplemented:
			return NewMomentoSvcErr(BadRequestError, grpcStatus.Message(), err)
		case codes.OutOfRange:
			return NewMomentoSvcErr(BadRequestError, grpcStatus.Message(), err)
		case codes.FailedPrecondition:
			return NewMomentoSvcErr(FailedPreconditionError, grpcStatus.Message(), err)
		case codes.Canceled:
			return NewMomentoSvcErr(CanceledError, grpcStatus.Message(), err)
		case codes.DeadlineExceeded:
			return NewMomentoSvcErr(TimeoutError, grpcStatus.Message(), err)
		case codes.PermissionDenied:
			return NewMomentoSvcErr(PermissionError, grpcStatus.Message(), err)
		case codes.Unauthenticated:
			return NewMomentoSvcErr(AuthenticationError, grpcStatus.Message(), err)
		case codes.ResourceExhausted:
			return NewMomentoSvcErr(LimitExceededError, grpcStatus.Message(), err)
		case codes.NotFound:
			// Use metadata to determine cause of not found error
			if len(metadata) > 1 {
				errData, ok := metadata[1]["err"]
				// In the absence of error metadata, return CacheNotFoundError.
				// This is currently re-mapped to a StoreNotFoundError in the PreviewStorageClient's
				// DeleteStore method.
				if !ok {
					return NewMomentoSvcErr(CacheNotFoundError, grpcStatus.Message(), err)
				}
				errCause := errData[0]
				if errCause == "item_not_found" {
					return NewMomentoSvcErr(ItemNotFoundError, grpcStatus.Message(), err)
				} else if errCause == "store_not_found" {
					return NewMomentoSvcErr(StoreNotFoundError, grpcStatus.Message(), err)
				}
			}
			return NewMomentoSvcErr(CacheNotFoundError, grpcStatus.Message(), err)
		case codes.AlreadyExists:
			return NewMomentoSvcErr(AlreadyExistsError, grpcStatus.Message(), err)
		case codes.Unknown:
			return NewMomentoSvcErr(UnknownServiceError, grpcStatus.Message(), err)
		case codes.Aborted:
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		case codes.Internal:
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		case codes.Unavailable:
			return NewMomentoSvcErr(ServerUnavailableError, grpcStatus.Message(), err)
		case codes.DataLoss:
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		default:
			return NewMomentoSvcErr(UnknownServiceError, InternalServerErrorMessage, err)
		}
	}
	return NewMomentoSvcErr(ClientSdkError, ClientSdkErrorMessage, err)
}

func NewConnectionError(err error) MomentoSvcErr {
	return NewMomentoSvcErr(ConnectionError, "Connection is in an unexpected state", err)
}
