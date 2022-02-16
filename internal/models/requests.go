package models

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ResolveRequest struct {
	AuthToken        string
	EndpointOverride string
}

type ControlGrpcManagerRequest struct {
	AuthToken string
	Endpoint  string
}

type DataGrpcManagerRequest struct {
	AuthToken string
	Endpoint  string
}

type ControlClientRequest struct {
	AuthToken string
	Endpoint  string
}

type DataClientRequest struct {
	AuthToken         string
	Endpoint          string
	DefaultTtlSeconds uint32
	RequestTimeout    uint32
}

type ConvertEcacheResultRequest struct {
	ECacheResult pb.ECacheResult
	Message      string
	OpName       string
}

func ConvertEcacheResult(request ConvertEcacheResultRequest) momentoerrors.MomentoSvcErr {
	return momentoerrors.NewMomentoSvcErr(
		momentoerrors.InternalServerError,
		fmt.Sprintf(
			"CacheService returned an unexpected result: %v for operation: %s with message: %s",
			request.ECacheResult, request.OpName, request.Message,
		),
		nil,
	)
}
