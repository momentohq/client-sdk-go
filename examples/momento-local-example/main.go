package main

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/middleware/impl"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const (
	cacheName  = "my-test-cache"
	defaultTTL = 60 * time.Second
)

func main() {
	ctx := context.Background()

	// Use the MomentoLocalProvider to create a credential provider to connect to a momento-local server.
	var credentialProvider, err = auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{
		Hostname: "127.0.0.1",
		Port:     8080,
	})
	if err != nil {
		panic(err)
	}

	loggerFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG)

	// Each test case uses a different momento-local configuration and cache client.
	TestGetSetWithDelays(ctx, credentialProvider, loggerFactory)
	TestGetSetWithErrors(ctx, credentialProvider, loggerFactory)
}

func CreateMomentoLocalCacheClient(
	ctx context.Context,
	credentialProvider auth.CredentialProvider,
	momentoLocalMetadata impl.MomentoLocalMetadataMiddlewareMetadataProps,
	loggerFactory logger.MomentoLoggerFactory,
) momento.CacheClient {
	// Create an instance of the momento-local middleware to attach to the cache client.
	momentoLocalMiddleware := impl.NewMomentoLocalMetadataMiddleware(impl.MomentoLocalMetadataMiddlewareProps{
		Props: middleware.Props{
			Logger: loggerFactory.GetLogger("momento-local-middleware"),
		},
		MomentoLocalMetadataMiddlewareMetadataProps: momentoLocalMetadata,
	})

	// Use WithMiddleware to attach the middleware to the cache client.
	cacheClient, err := momento.NewCacheClient(
		config.LaptopLatestWithLogger(loggerFactory).WithMiddleware([]middleware.Middleware{momentoLocalMiddleware}),
		credentialProvider,
		defaultTTL,
	)
	if err != nil {
		panic(err)
	}

	// Create cache if it does not yet exist.
	_, err = cacheClient.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	return cacheClient
}

func TestGetSetWithDelays(ctx context.Context, credentialProvider auth.CredentialProvider, loggerFactory logger.MomentoLoggerFactory) {
	delayMillis := 1000 // how long to delay the response
	delayCount := 1     // how many times to delay the response
	momentoLocalMetadata := impl.MomentoLocalMetadataMiddlewareMetadataProps{
		DelayRpcList: &[]string{"set"},
		DelayMillis:  &delayMillis,
		DelayCount:   &delayCount,
	}
	cacheClient := CreateMomentoLocalCacheClient(ctx, credentialProvider, momentoLocalMetadata, loggerFactory)

	// delayCount is 1, meaning we expect only the response to this Set request to be delayed by delayMillis.
	startTime := time.Now()
	_, err := cacheClient.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
	})
	if err != nil {
		panic(err)
	}
	duration := time.Since(startTime)
	fmt.Printf("First set request took %s\n", duration)

	// Another Set request should not be delayed since delayCount is 1 and we've already delayed a Set response once.
	startTime = time.Now()
	_, err = cacheClient.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
	})
	if err != nil {
		panic(err)
	}
	duration = time.Since(startTime)
	fmt.Printf("Second set request took %s\n", duration)

	// This get request should not be delayed as it's not in the delayRpcList.
	startTime = time.Now()
	resp, err := cacheClient.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
	})
	if err != nil {
		panic(err)
	}
	duration = time.Since(startTime)
	fmt.Printf("Get request took %s\n", duration)
	switch r := resp.(type) {
	case *responses.GetHit:
		fmt.Printf("Get hit: %s\n", r.ValueString())
	case *responses.GetMiss:
		fmt.Printf("Get miss: %v\n", r)
	}
}

func TestGetSetWithErrors(ctx context.Context, credentialProvider auth.CredentialProvider, loggerFactory logger.MomentoLoggerFactory) {
	errorCount := 2
	errorStatus := "unavailable"
	momentoLocalMetadata := impl.MomentoLocalMetadataMiddlewareMetadataProps{
		ErrorRpcList: &[]string{"set"},
		ErrorCount:   &errorCount,
		ReturnError:  &errorStatus,
	}
	cacheClient := CreateMomentoLocalCacheClient(ctx, credentialProvider, momentoLocalMetadata, loggerFactory)

	// This set request should return an error twice, then succeed. The logs should show the default retry strategy,
	// FixedCountRetryStrategy, being used with a maximum of 3 attempts.
	_, err := cacheClient.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
	})
	if err != nil {
		fmt.Printf("Set request returned error: %v\n", err)
	}
}
