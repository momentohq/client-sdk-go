package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	const (
		cacheName             = "cache"
		itemDefaultTTLSeconds = 60
	)

	// Initializes Momento
	client, err := momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credentialProvider,
		DefaultTTLSeconds:  itemDefaultTTLSeconds,
	})
	if err != nil {
		panic(err)
	}
	logger := config.LatestLaptopConfig().GetLogger()
	// Create Cache and check if CacheName exists
	logger.Info("Creating cache")
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
			logger.Info(fmt.Sprintf("%s", cacheInfo.Name()))
		}
		token = listCacheResp.NextToken()
		if token == "" {
			break
		}
	}

	// Sets key with default TTL and gets value with that key
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	logger.Info(fmt.Sprintf("Setting key: %s, value: %s", key, value))
	err = client.Set(ctx, &momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Getting key: %s", key))
	resp, err := client.Get(ctx, &momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}
	if resp.IsHit() {
		logger.Info(fmt.Sprintf("Lookup resulted in cahce HIT. value=%s", resp.AsHit().ValueString()))
	} else {
		logger.Info(fmt.Sprintf("Look up did not find a value key=%s", key))
	}

	// Permanently delete the cache
	err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("Cache named %s is deleted", cacheName))
}
