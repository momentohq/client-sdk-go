package main

import (
	"context"
	"errors"
	"log"

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
		cacheName             = "cache"
		itemDefaultTtlSeconds = 60
	)

	// Initializes Momento
	client, err := momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credentialProvider,
		DefaultTtlSeconds:  itemDefaultTtlSeconds,
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

	// List caches
	token := ""
	for {
		listCacheResp, err := client.ListCaches(ctx, &momento.ListCachesRequest{NextToken: token})
		if err != nil {
			panic(err)
		}
		for _, cacheInfo := range listCacheResp.Caches() {
			log.Printf("%s\n", cacheInfo.Name())
		}
		token = listCacheResp.NextToken()
		if token == "" {
			break
		}
	}

	// Sets key with default TTL and gets value with that key
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err = client.Set(ctx, &momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(ctx, &momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Lookup resulted in a : %s\n", resp.Result())
	log.Printf("Looked up value: %s\n", resp.StringValue())

	// Permanently delete the cache
	err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}
