package momentoerrors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type MomentoSvcErr interface {
	error
	Code() string
	Message() string
}

func NewMomentoSvcErr(code string, message string) MomentoSvcErr {
	return newMomentoSvcErr(code, message)
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
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "Canceled":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "DeadlineExceeded":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "PermissionDenied":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "Unauthenticated":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "ResourceExhausted":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "NotFound":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "AlreadyExists":
			return NewMomentoSvcErr(grpcStatus.Code().String(), grpcStatus.Message())
		case "Unknown":
			fallthrough
		case "Aborted":
			fallthrough
		case "Internal":
			fallthrough
		case "Unavailable":
			return NewMomentoSvcErr(grpcStatus.Code().String(), fmt.Sprintf("Unable to reach request endpoint. Request failed with %s", grpcStatus.Message()))
		case "DataLoss":
			fallthrough
		default:
			return NewMomentoSvcErr(InternalServerError, InternalServerErrorMessage)
		}
	}
	return NewMomentoSvcErr(ClientSdkError, ClientSdkErrorMessage)
}
