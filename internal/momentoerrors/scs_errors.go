package momentoerrors

import (
	"fmt"

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

const InternalServerErrorMessage = "Unexpected exception occurred while trying to fulfill the request."
const ClientSdkErrorMessage = "SDK Failed to process the request."
const InvalidArgumentError = "InvalidArgumentError"
const InternalServerError = "InternalServerError"
const ClientSdkError = "ClientSdkError"

func ConvertSvcErr(err error) MomentoSvcErr {
	if grpcStatus, ok := status.FromError(err); ok {
		switch grpcStatus.Code().String() {
		case "InvalidArgument":
			fallthrough
		case "Unimplemented":
			fallthrough
		case "OutOfRange":
			fallthrough
		case "FailedPrecondition":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "Canceled":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "DeadlineExceeded":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "PermissionDenied":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "Unauthenticated":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "ResourceExhausted":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "NotFound":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "AlreadyExists":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message(), err)
		case "Unknown":
			fallthrough
		case "Aborted":
			fallthrough
		case "Internal":
			fallthrough
		case "Unavailable":
			return NewMomentoSvcErr(grpcStatus.Code().String(), fmt.Sprintf("Unable to reach request endpoint. Request failed with %s", grpcStatus.Message()), err)
		case "DataLoss":
			fallthrough
		default:
			return NewMomentoSvcErr(InternalServerError, InternalServerErrorMessage, err)
		}
	}
	return NewMomentoSvcErr(ClientSdkError, ClientSdkErrorMessage, err)
}
