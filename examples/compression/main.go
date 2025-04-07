package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const (
	cacheName             = "my-test-cache"
	itemDefaultTTLSeconds = 60
	longStringValue       = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	// Initialize Momento client with compression middleware
	loggerFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)
	myConfig := config.LaptopLatestWithLogger(loggerFactory).WithMiddleware(
		[]middleware.Middleware{
			NewCompressionMiddleware(middleware.Props{Logger: loggerFactory.GetLogger("compression-middleware")}),
		},
	)
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		myConfig,
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}

	// Create cache
	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	// Compression example 1: SetIfAbsentOrHashEqual and GetWithHash

	var storedHashValue []byte
	setResp, err := client.SetIfAbsentOrHashEqual(ctx, &momento.SetIfAbsentOrHashEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key-1"),
		Value:     momento.String(longStringValue),
		HashEqual: momento.Bytes("set-if-absent-or-hash-equal"),
	})
	if err != nil {
		fmt.Printf("[SetIfAbsentOrHashEqual] Error: %s\n", err.Error())
	}
	switch r := setResp.(type) {
	case *responses.SetIfAbsentOrHashEqualStored:
		fmt.Printf("[SetIfAbsentOrHashEqual] Stored value\n")
		storedHashValue = r.HashByte()
	case *responses.SetIfAbsentOrHashEqualNotStored:
		fmt.Printf("[SetIfAbsentOrHashEqual] Unable to store the value\n")
	}

	getWithHashResp, err := client.GetWithHash(ctx, &momento.GetWithHashRequest{
		CacheName: cacheName,
		Key:       momento.String("key-1"),
	})
	if err != nil {
		// TODO: failed to decompress message isn't propagated?
		momentoError, ok := err.(momento.MomentoError)
		if ok {
			fmt.Printf("[GetWithHash] Error: %s\n", momentoError.Message())
		} else {
			fmt.Printf("[GetWithHash] Error: %s\n", err.Error())
		}
	}
	switch r := getWithHashResp.(type) {
	case *responses.GetWithHashHit:
		fmt.Printf("[GetWithHash] Lookup resulted in cache HIT\n")
		if bytes.Equal(r.HashByte(), storedHashValue) {
			fmt.Printf("[GetWithHash] Hash is equal to stored hash value\n")
		} else {
			fmt.Printf("[GetWithHash] Hash is not equal to stored hash value\n")
		}
	case *responses.GetWithHashMiss:
		fmt.Printf("[GetWithHash] Lookup resulted in cache MISS\n")
	}

	// Compression example 2: Regular Set and Get

	_, err = client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String("key-2"),
		Value:     momento.String(longStringValue),
	})
	if err != nil {
		fmt.Printf("[Set] Error: %s\n", err.Error())
	}
	fmt.Printf("[Set] Stored value\n")

	getResp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String("key-2"),
	})
	if err != nil {
		fmt.Printf("[Get] Error: %s\n", err.Error())
	}
	switch getResp.(type) {
	case *responses.GetHit:
		fmt.Printf("[Get] Lookup resulted in cache HIT\n")
	case *responses.GetMiss:
		fmt.Printf("[Get] Lookup resulted in cache MISS\n")
	}

	// Cleanup: delete the cache
	_, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}
