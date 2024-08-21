package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	DefaultClient             = "defaultClient"
	WithDefaultCache          = "withDefaultCache"
	WithConsistentReadConcern = "withConsistentReadConcern"
)

type SharedContext struct {
	Client                          momento.CacheClient
	ClientWithDefaultCacheName      momento.CacheClient
	ClientWithConsistentReadConcern momento.CacheClient
	DefaultCacheName                string
	TopicClient                     momento.TopicClient
	CacheName                       string
	StoreName                       string
	CleanupStore                    bool
	Ctx                             context.Context
	DefaultTtl                      time.Duration
	Configuration                   config.Configuration
	TopicConfiguration              config.TopicsConfiguration
	CredentialProvider              auth.CredentialProvider
	AuthClient                      momento.AuthClient
	AuthConfiguration               config.AuthConfiguration
	LeaderboardClient               momento.PreviewLeaderboardClient
	LeaderboardConfiguration        config.LeaderboardConfiguration
	StorageClient                   momento.PreviewStorageClient
	StorageConfiguration            config.StorageConfiguration
}

func NewSharedContext() SharedContext {
	shared := SharedContext{}

	shared.Ctx = context.Background()
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	shared.CredentialProvider = credentialProvider
	shared.Configuration = config.LaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory()).WithClientTimeout(15 * time.Second)
	shared.TopicConfiguration = config.TopicsDefaultWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.AuthConfiguration = config.AuthDefaultWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.LeaderboardConfiguration = config.LeaderboardDefaultWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.StorageConfiguration = config.StorageLaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory())
	shared.DefaultTtl = 3 * time.Second

	client, err := momento.NewCacheClient(shared.Configuration, shared.CredentialProvider, shared.DefaultTtl)
	if err != nil {
		panic(err)
	}

	defaultCacheName := fmt.Sprintf("golang-default-%s", uuid.NewString())
	clientDefaultCacheName, err := momento.NewCacheClientWithDefaultCache(
		shared.Configuration, shared.CredentialProvider, shared.DefaultTtl, defaultCacheName,
	)
	if err != nil {
		panic(err)
	}

	consistentReadConcernClient, err := momento.NewCacheClientWithDefaultCache(
		shared.Configuration.WithReadConcern(config.CONSISTENT), shared.CredentialProvider, shared.DefaultTtl, defaultCacheName,
	)
	if err != nil {
		panic(err)
	}

	topicClient, err := momento.NewTopicClient(shared.TopicConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	authClient, err := momento.NewAuthClient(shared.AuthConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	leaderboardClient, err := momento.NewPreviewLeaderboardClient(shared.LeaderboardConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	storageClient, err := momento.NewPreviewStorageClient(shared.StorageConfiguration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	shared.Client = client
	shared.ClientWithDefaultCacheName = clientDefaultCacheName
	shared.ClientWithConsistentReadConcern = consistentReadConcernClient
	shared.DefaultCacheName = defaultCacheName
	shared.TopicClient = topicClient
	shared.AuthClient = authClient
	shared.LeaderboardClient = leaderboardClient
	shared.StorageClient = storageClient

	shared.CacheName = fmt.Sprintf("golang-%s", uuid.NewString())
	shared.StoreName = fmt.Sprintf("golang-%s", uuid.NewString())

	return shared
}

func (shared SharedContext) GetClientPrereqsForType(clientType string) (momento.CacheClient, string) {
	var client momento.CacheClient
	var cacheName string
	if clientType == WithDefaultCache {
		client = shared.ClientWithDefaultCacheName
		cacheName = ""
	} else if clientType == DefaultClient {
		client = shared.Client
		cacheName = shared.CacheName
	} else if clientType == WithConsistentReadConcern {
		client = shared.ClientWithConsistentReadConcern
		cacheName = ""
	} else {
		panic("invalid client type")
	}
	return client, cacheName
}

func (shared SharedContext) CreateDefaultCaches() {
	_, err := shared.Client.CreateCache(shared.Ctx, &momento.CreateCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}
	_, err = shared.Client.CreateCache(shared.Ctx, &momento.CreateCacheRequest{CacheName: shared.DefaultCacheName})
	if err != nil {
		panic(err)
	}
}

func (shared *SharedContext) CreateDefaultStores() {
	_, err := shared.StorageClient.CreateStore(shared.Ctx, &momento.CreateStoreRequest{StoreName: shared.StoreName})
	if err != nil {
		panic(err)
	}
	shared.CleanupStore = true
}

func (shared SharedContext) Close() {
	_, err := shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}
	_, err = shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.DefaultCacheName})
	if err != nil {
		panic(err)
	}
	if shared.CleanupStore {
		_, err = shared.StorageClient.DeleteStore(shared.Ctx, &momento.DeleteStoreRequest{StoreName: shared.StoreName})
		if err != nil {
			if err.(momento.MomentoError).Code() != momento.StoreNotFoundError {
				panic(err)
			}
		}
	}

	shared.Client.Close()
	shared.TopicClient.Close()
	shared.AuthClient.Close()
	shared.LeaderboardClient.Close()
	shared.StorageClient.Close()
}
