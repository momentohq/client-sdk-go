package models

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
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

type TokenGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type DataStreamGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type TopicStreamGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
}

type PingGrpcManagerRequest struct {
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
	TopicsConfiguration config.TopicsConfiguration
	CredentialProvider  auth.CredentialProvider
	Log                 logger.MomentoLogger
}

type TopicClientRequest struct {
	CredentialProvider auth.CredentialProvider
	Log                logger.MomentoLogger
}

type TokenClientRequest struct {
	CredentialProvider auth.CredentialProvider
	Log                logger.MomentoLogger
}

type PingClientRequest struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}
