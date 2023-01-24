package models

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/auth"
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
	CredentialProvider auth.CredentialProvider
}

type DataClientRequest struct {
	CredentialProvider    auth.CredentialProvider
	DefaultTtlSeconds     uint32
	RequestTimeoutSeconds uint32
}

type PubSubClientRequest struct {
	CredentialProvider auth.CredentialProvider
	// TODO think about timeout settings more
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
