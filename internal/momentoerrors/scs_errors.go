package momentoerrors

import (
	"google.golang.org/grpc/status"
)

type MomentoSvcErr interface {
	error
	Code() string
	Message() string
	OriginalErr() error
}

func NewMomentoSvcErr(code string, message string, originalErr error) MomentoSvcErr {
	return newMomentoSvcErr(code, message, originalErr)
}

const (
	InvalidArgumentError       = "InvalidArgumentError"
	InternalServerError        = "InternalServerError"
	ClientSdkError             = "ClientSdkError"
	BadRequestError            = "BadRequestError"
	CanceledError              = "CanceledError"
	TimeoutError               = "TimeoutError"
	PermissionError            = "PermissionError"
	AuthenticationError        = "AuthenticationError"
	LimitExceededError         = "LimitExceededError"
	NotFoundError              = "NotFoundError"
	AlreadyExistsError         = "AlreadyExistsError"
	InternalServerErrorMessage = "CacheService failed with an internal error"
	ClientSdkErrorMessage      = "SDK Failed to process the request."
)

func ConvertSvcErr(err error) MomentoSvcErr {
	if grpcStatus, ok := status.FromError(err); ok {
		switch grpcStatus.Code().String() {
		case "InvalidArgument":
			return NewMomentoSvcErr(BadRequestError, grpcStatus.Message(), err)
		case "Unimplemented":
			return NewMomentoSvcErr(BadRequestError, grpcStatus.Message(), err)
		case "OutOfRange":
			return NewMomentoSvcErr(BadRequestError, grpcStatus.Message(), err)
		case "FailedPrecondition":
			return NewMomentoSvcErr(BadRequestError, grpcStatus.Message(), err)
		case "Canceled":
			return NewMomentoSvcErr(CanceledError, grpcStatus.Message(), err)
		case "DeadlineExceeded":
			return NewMomentoSvcErr(TimeoutError, grpcStatus.Message(), err)
		case "PermissionDenied":
			return NewMomentoSvcErr(PermissionError, grpcStatus.Message(), err)
		case "Unauthenticated":
			return NewMomentoSvcErr(AuthenticationError, grpcStatus.Message(), err)
		case "ResourceExhausted":
			return NewMomentoSvcErr(LimitExceededError, grpcStatus.Message(), err)
		case "NotFound":
			return NewMomentoSvcErr(NotFoundError, grpcStatus.Message(), err)
		case "AlreadyExists":
			return NewMomentoSvcErr(AlreadyExistsError, grpcStatus.Message(), err)
		case "Unknown":
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		case "Aborted":
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		case "Internal":
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		case "Unavailable":
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		case "DataLoss":
			return NewMomentoSvcErr(InternalServerError, grpcStatus.Message(), err)
		default:
			return NewMomentoSvcErr(InternalServerError, InternalServerErrorMessage, err)
		}
	}
	return NewMomentoSvcErr(ClientSdkError, ClientSdkErrorMessage, err)
}
