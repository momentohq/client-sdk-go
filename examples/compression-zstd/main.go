package main

import (
	"bytes"
	"context"
	"log"
	"time"

	zstdMiddleware "github.com/momentohq/client-sdk-go-compression-zstd/zstd_compression"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/compression"
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

	// Initialize Momento client with zstd compression middleware.
	loggerFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE)
	myConfig := config.LaptopLatestWithLogger(loggerFactory).WithMiddleware(
		[]middleware.Middleware{
			zstdMiddleware.NewZstdCompressionMiddleware(zstdMiddleware.ZstdCompressionMiddlewareProps{
				CompressionStrategyProps: compression.CompressionStrategyProps{
					CompressionLevel: compression.CompressionLevelDefault,
					Logger:           loggerFactory.GetLogger("compression-middleware"),
				},
			}),
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
		log.Printf("[SetIfAbsentOrHashEqual] Error: %s\n", err.Error())
	}
	switch r := setResp.(type) {
	case *responses.SetIfAbsentOrHashEqualStored:
		log.Printf("[SetIfAbsentOrHashEqual] Stored value\n")
		storedHashValue = r.HashByte()
	case *responses.SetIfAbsentOrHashEqualNotStored:
		log.Printf("[SetIfAbsentOrHashEqual] Unable to store the value\n")
	}

	getWithHashResp, err := client.GetWithHash(ctx, &momento.GetWithHashRequest{
		CacheName: cacheName,
		Key:       momento.String("key-1"),
	})
	if err != nil {
		log.Printf("[GetWithHash] Error: %s\n", err)
	}
	switch r := getWithHashResp.(type) {
	case *responses.GetWithHashHit:
		log.Printf("[GetWithHash] Lookup resulted in cache HIT\n")
		if bytes.Equal(r.HashByte(), storedHashValue) {
			log.Printf("[GetWithHash] Hash is equal to stored hash value\n")
		} else {
			log.Printf("[GetWithHash] Hash is not equal to stored hash value\n")
		}
	case *responses.GetWithHashMiss:
		log.Printf("[GetWithHash] Lookup resulted in cache MISS\n")
	}

	// Compression example 2: Regular Set and Get

	_, err = client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String("key-2"),
		Value:     momento.String(longStringValue),
	})
	if err != nil {
		log.Printf("[Set] Error: %s\n", err.Error())
	}
	log.Printf("[Set] Stored value\n")

	getResp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String("key-2"),
	})
	if err != nil {
		log.Printf("[Get] Error: %s\n", err.Error())
	}
	switch getResp.(type) {
	case *responses.GetHit:
		log.Printf("[Get] Lookup resulted in cache HIT\n")
	case *responses.GetMiss:
		log.Printf("[Get] Lookup resulted in cache MISS\n")
	}

	// Cleanup: delete the cache
	_, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}
