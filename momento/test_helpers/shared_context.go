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
	DefaultClient    = "defaultClient"
	WithDefaultCache = "withDefaultCache"
)

type SharedContext struct {
	Client                     momento.CacheClient
	ClientWithDefaultCacheName momento.CacheClient
	DefaultCacheName           string
	TopicClient                momento.TopicClient
	CacheName                  string
	CollectionName             string
	Ctx                        context.Context
	DefaultTtl                 time.Duration
	Configuration              config.Configuration
	CredentialProvider         auth.CredentialProvider
}

func NewSharedContext() SharedContext {
	shared := SharedContext{}

	shared.Ctx = context.Background()
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	shared.CredentialProvider = credentialProvider
	shared.Configuration = config.LaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory()).WithClientTimeout(15 * time.Second)
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

	topicClient, err := momento.NewTopicClient(shared.Configuration, shared.CredentialProvider)
	if err != nil {
		panic(err)
	}

	shared.Client = client
	shared.ClientWithDefaultCacheName = clientDefaultCacheName
	shared.DefaultCacheName = defaultCacheName
	shared.TopicClient = topicClient

	shared.CacheName = fmt.Sprintf("golang-%s", uuid.NewString())
	shared.CollectionName = uuid.NewString()

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

func (shared SharedContext) Close() {
	_, err := shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}
	_, err = shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.DefaultCacheName})
	if err != nil {
		panic(err)
	}

	shared.Client.Close()
}
