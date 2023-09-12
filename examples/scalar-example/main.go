package main

import (
	"context"
	"log"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	const (
		cacheName             = "my-test-cache"
		itemDefaultTTLSeconds = 60
	)

	// Initializes Momento
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}

	// Create Cache
	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	// Sets key with default TTL and gets value with that key
	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err = client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(value),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.GetHit:
		log.Printf("Lookup resulted in cache HIT. value=%s\n", r.ValueString())
	case *responses.GetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}

	// Permanently delete the cache
	if _, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}
