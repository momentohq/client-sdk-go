package models

import pb "github.com/momentohq/client-sdk-go/internal/protos"

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
}

type ConvertEcacheResultRequest struct {
	ECacheResult pb.ECacheResult
	Message      string
	OpName       string
}
