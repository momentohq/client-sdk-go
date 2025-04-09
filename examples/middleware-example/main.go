package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"

	"github.com/google/uuid"
)

const (
	cacheName             = "my-test-cache"
	itemDefaultTTLSeconds = 60
)

func doWork(ctx context.Context, client momento.CacheClient, index int) {
	// Sets key with default TTL and gets value with that key
	key := uuid.NewString()
	value := fmt.Sprintf("%d", index)
	log.Printf("#%d setting key: %s, value: %s\n", index, key, value)
	_, err := client.Set(ctx, &momento.SetRequest{
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
}

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	loggerFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)
	myConfig := config.LaptopLatest().WithMiddleware(
		[]middleware.Middleware{
			// This is a middleware bundled with the SDK that we access through the `middleware` package.
			// It is configured to process only GetRequest and SetRequest types, ignoring all other request types.
			middleware.NewInFlightRequestCountMiddleware(middleware.Props{
				Logger:       loggerFactory.GetLogger("inflight-request-count"),
				IncludeTypes: []interface{}{momento.GetRequest{}, momento.SetRequest{}},
			}),
			// These are custom middleware built using the `middleware.Middleware` interface
			NewTimingMiddleware(middleware.Props{Logger: loggerFactory.GetLogger("timing")}),
			NewLoggingMiddleware(middleware.Props{Logger: loggerFactory.GetLogger("logging-middleware")}),
		},
	)

	// Initializes Momento
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		myConfig,
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

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		// avoid reuse of the same i value in each closure
		i := i
		go func() {
			defer wg.Done()
			doWork(ctx, client, i)
		}()
	}

	wg.Wait()

	// Permanently delete the cache
	if _, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}
