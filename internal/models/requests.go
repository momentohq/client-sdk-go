package models

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ControlGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type DataGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type LocalDataGrpcManagerRequest struct {
	Endpoint string
}

type ControlClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

type DataClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTtlSeconds  uint32
}

type PubSubClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

type NewLocalPubSubClientRequest struct {
	Port int
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
