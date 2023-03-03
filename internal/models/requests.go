package models

import (
	"time"

	"github.com/momentohq/client-sdk-go/internal/retry"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
)

type ControlGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	RetryStrategy      retry.Strategy
}

type DataGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	RetryStrategy      retry.Strategy
}

type DataStreamGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type LocalDataGrpcManagerRequest struct {
	Endpoint string
}

type CreateCacheRequest struct {
	CacheName string
}

type DeleteCacheRequest struct {
	CacheName string
}

type ListCachesRequest struct {
	NextToken string
}

type ControlClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

type DataClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTtl         time.Duration
}

type PubSubClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}
