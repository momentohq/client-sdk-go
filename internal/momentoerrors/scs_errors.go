package momentoerrors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/status"
)

func InvalidInputError(errMessage string) error {
	return errors.New("InvalidInputError: " + errMessage)
}

func InternalServerError(errMessage string) error {
	return errors.New("InternalServerError: " + errMessage)
}

func ClientSdkError(errMessage string) error {
	return errors.New("ClientSdkError: " + errMessage)
}

const (
	AlreadyExists    = "AlreadyExists"
	InvalidArgument  = "InvalidArgument"
	NotFound         = "NotFound"
	PermissionDenied = "PermissionDenied"
)

func GrpcErrorConverter(grpcErr error) error {
	if grpcStatus, ok := status.FromError(grpcErr); ok {
		switch grpcStatus.Code().String() {
		case AlreadyExists:
			return fmt.Errorf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message())
		case InvalidArgument:
			return fmt.Errorf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message())
		case NotFound:
			return fmt.Errorf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message())
		case PermissionDenied:
			return fmt.Errorf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message())
		}
	} else {
		return InternalServerError("CacheService failed with an internal error")
	}
	return ClientSdkError(fmt.Sprintf("Operation failed with error: %s", grpcErr.Error()))
}
