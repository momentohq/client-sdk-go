package momentoerrors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type MomentoBaseError interface {
	error
	Code() string
	Message() string
}

func NewMomentoBaseError(code string, message string) MomentoBaseError {
	return newMomentoBaseError(code, message)
}

const InternalServerErrorMessage = "Unexpected exception occurred while trying to fulfill the request."
const ClientSdkErrorMessage = "SDK Failed to process the request."

func ConvertError(err error) MomentoBaseError {
	if grpcStatus, ok := status.FromError(err); ok {
		switch grpcStatus.Code().String() {
		case "InvalidArgument":
			fallthrough
		case "Unimplemented":
			fallthrough
		case "OutOfRange":
			fallthrough
		case "FailedPrecondition":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "Canceled":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "DeadlineExceeded":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "PermissionDenied":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "Unauthenticated":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "ResourceExhausted":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "NotFound":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "AlreadyExists":
			return NewMomentoBaseError(grpcStatus.Code().String(), grpcStatus.Message())
		case "Unknown":
			fallthrough
		case "Aborted":
			fallthrough
		case "Internal":
			fallthrough
		case "Unavailable":
			return NewMomentoBaseError(grpcStatus.Code().String(), fmt.Sprintf("Unable to reach request endpoint. Request failed with %s", grpcStatus.Message()))
		case "DataLoss":
			fallthrough
		default:
			return NewMomentoBaseError("InternalServerError", InternalServerErrorMessage)
		}
	}
	return NewMomentoBaseError("ClientSdkError", ClientSdkErrorMessage)
}
