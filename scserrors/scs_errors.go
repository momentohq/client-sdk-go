package scserrors

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

func GrpcErrorConverter(grpcErr error) error {
	if grpcStatus, ok := status.FromError(grpcErr); ok {
		switch grpcStatus.Code().String() {
		case "AlreadyExists":
			return fmt.Errorf(fmt.Sprintf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message()))
		case "InvalidArgument":
			return fmt.Errorf(fmt.Sprintf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message()))
		case "NotFound":
			return fmt.Errorf(fmt.Sprintf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message()))
		case "PermissionDenied":
			return fmt.Errorf(fmt.Sprintf("%s: %s", grpcStatus.Code().String(), grpcStatus.Message()))
		}
	}
	return InternalServerError("CacheService failed with an internal error'")
}
