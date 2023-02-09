package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	const (
		cacheName      = "my-test-cache"
		itemDefaultTTL = time.Duration(time.Second * 60)
	)

	// Initializes Momento
	client, err := momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credentialProvider,
		DefaultTTL:         itemDefaultTTL,
	})
	if err != nil {
		panic(err)
	}

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
	log.Printf("Setting key: %s, value: %s\n", key, value)
	err = client.Set(ctx, &momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       &momento.StringBytes{Text: key},
		Value:     &momento.StringBytes{Text: value},
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(ctx, &momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       &momento.StringBytes{Text: key},
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *momento.CacheGetHit:
		log.Printf("Lookup resulted in cahce HIT. value=%s\n", r.ValueString())
	case *momento.CacheGetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}

	// Permanently delete the cache
	if err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}
