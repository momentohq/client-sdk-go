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
	GrpcConfiguration  config.GrpcConfiguration
}

type DataGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	RetryStrategy      retry.Strategy
	ReadConcern        config.ReadConcern
	GrpcConfiguration  config.GrpcConfiguration
	EagerConnect       bool
}

type TokenGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	GrpcConfiguration  config.GrpcConfiguration
}

type DataStreamGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	GrpcConfiguration  config.GrpcConfiguration
}

type TopicStreamGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	GrpcConfiguration  config.GrpcConfiguration
}

type PingGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	GrpcConfiguration  config.GrpcConfiguration
}

type LeaderboardGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	GrpcConfiguration  config.GrpcConfiguration
}

type StoreGrpcManagerRequest struct {
	CredentialProvider auth.CredentialProvider
	GrpcConfiguration  config.GrpcConfiguration
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

type LeaderboardClientRequest struct {
	Configuration      config.LeaderboardConfiguration
	CredentialProvider auth.CredentialProvider
}

type StorageDataClientRequest struct {
	CredentialProvider auth.CredentialProvider
	Configuration      config.StorageConfiguration
	Log                logger.MomentoLogger
}

type CreateStoreRequest struct {
	StoreName string
}

type DeleteStoreRequest struct {
	StoreName string
}

type ListStoresRequest struct {
	NextToken string
}
