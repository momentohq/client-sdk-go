package momentoerrors

import (
	"google.golang.org/grpc/codes"
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
	UnknownServiceError        = "UnknownServiceError"
	ServerUnavailableError     = "ServerUnavailableError"
	FailedPreconditionError    = "FailedPreconditionError"
	InternalServerErrorMessage = "CacheService failed with an internal error"
	ClientSdkErrorMessage      = "SDK Failed to process the request."
)

func ConvertSvcErr(err error) MomentoSvcErr {
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
			return NewMomentoSvcErr(NotFoundError, grpcStatus.Message(), err)
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
