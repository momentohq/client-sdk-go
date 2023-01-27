package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
)

func main() {
	var authToken = os.Getenv("MOMENTO_AUTH_TOKEN")
	const (
		cacheName             = "cache"
		itemDefaultTTLSeconds = 60
	)

	if authToken == "" {
		log.Fatal("Missing required environment variable MOMENTO_AUTH_TOKEN")
	}

	// Initializes Momento
	client, err := momento.NewSimpleCacheClient(authToken, itemDefaultTTLSeconds)
	if err != nil {
		panic(err)
	}

	// Create Cache and check if CacheName exists
	err = client.CreateCache(&momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil && err.Error() != momento.AlreadyExistsError {
		panic(err)
	}
	log.Printf("Cache named %s is created\n", cacheName)

	// List caches
	token := ""
	for {
		listCacheResp, err := client.ListCaches(&momento.ListCachesRequest{NextToken: token})
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
	_, err = client.Set(&momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(&momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Lookup resulted in a : %s\n", resp.Result())
	log.Printf("Looked up value: %s\n", resp.StringValue())

	// Permanently delete the cache
	err = client.DeleteCache(&momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}
