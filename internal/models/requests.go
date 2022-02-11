package models

import (
	"fmt"

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
	DataCtxTimeout    *uint32
}

type ConvertEcacheResultRequest struct {
	ECacheResult pb.ECacheResult
	Message      string
	OpName       string
}

func ConvertEcacheResult(request ConvertEcacheResultRequest) error {
	return fmt.Errorf("CacheService returned an unexpected result: %v for operation: %s with message: %s", request.ECacheResult, request.OpName, request.Message)
}
