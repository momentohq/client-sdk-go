package models

import (
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
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
type CacheGetRequest struct {
	CacheName string
	Key       interface{}
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
