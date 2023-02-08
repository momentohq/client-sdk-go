package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/momento"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	const (
		cacheName             = "my-test-cache"
		itemDefaultTTLSeconds = 60
	)

	// Initializes Momento
	configuration := config.LatestLaptopConfig(&logger.BuiltinMomentoLoggerFactory{})
	client, err := momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      configuration,
		CredentialProvider: credentialProvider,
		DefaultTTLSeconds:  itemDefaultTTLSeconds,
	})
	if err != nil {
		panic(err)
	}

	mLogger := configuration.GetLoggerFactory().GetLogger("examples-main")
	// Create Cache and check if CacheName exists
	err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.AlreadyExistsError {
				panic(err)
			}
		}
	}

	// Sets key with default TTL and gets value with that key
	key := uuid.NewString()
	value := uuid.NewString()
	mLogger.Info("Setting key/value", "key", key, "value", value)
	err = client.Set(ctx, &momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       &momento.StringBytes{Text: key},
		Value:     &momento.StringBytes{Text: value},
	})
	if err != nil {
		panic(err)
	}

	mLogger.Info("Getting key", "key", key)
	resp, err := client.Get(ctx, &momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       &momento.StringBytes{Text: key},
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *momento.CacheGetHit:
		mLogger.Info("Lookup resulted in cache HIT", "key", key, "value", r.ValueString())
	case *momento.CacheGetMiss:
		mLogger.Info("Look up did not find a value", "key", key)
	}

	// Permanently delete the cache
	if err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	mLogger.Info(fmt.Sprintf("Cache named %s is deleted\n", cacheName))
}
