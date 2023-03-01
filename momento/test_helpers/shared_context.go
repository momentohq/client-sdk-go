package helpers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

type SharedContext struct {
	Client             momento.CacheClient
	CacheName          string
	CollectionName     string
	Ctx                context.Context
	DefaultTtl         time.Duration
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

func NewSharedContext() SharedContext {
	shared := SharedContext{}

	shared.Ctx = context.Background()
	credentialProvider, _ := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	shared.CredentialProvider = credentialProvider
	shared.Configuration = config.LatestLaptopConfig()
	shared.DefaultTtl = 3 * time.Second

	var err error
	client, err := momento.NewCacheClient(shared.Configuration, shared.CredentialProvider, shared.DefaultTtl)
	if err != nil {
		panic(err)
	}

	shared.Client = client

	shared.CacheName = uuid.NewString()
	shared.CollectionName = uuid.NewString()

	return shared
}

func (shared SharedContext) CreateDefaultCache() {
	_, err := shared.Client.CreateCache(shared.Ctx, &momento.CreateCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}
}

func (shared SharedContext) Close() {
	_, err := shared.Client.DeleteCache(shared.Ctx, &momento.DeleteCacheRequest{CacheName: shared.CacheName})
	if err != nil {
		panic(err)
	}

	shared.Client.Close()
}
