// main.go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const (
	CACHE_NAME  = "cache"
	CACHE_KEY   = "key"
	CACHE_VALUE = "value"
)

var (
	cachedAuthToken    string = ""
	cachedEndpoint     string = ""
	secretsClient, _          = secretcache.New()
	momentoCacheClient momento.CacheClient
)

func handler() (string, error) {
	cacheClient, err := getCacheClient()
	if err != nil {
		return "Something went wrong getting the cache client", err
	}
	ctx := context.Background()

	_, err = cacheClient.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: CACHE_NAME,
	})
	if err != nil {
		panic(err)
	}

	_, err = cacheClient.Set(ctx, &momento.SetRequest{
		CacheName: CACHE_NAME,
		Key:       momento.String(CACHE_KEY),
		Value:     momento.String(CACHE_VALUE),
	})
	if err != nil {
		panic(err)
	}

	resp, err := cacheClient.Get(ctx, &momento.GetRequest{
		CacheName: CACHE_NAME,
		Key:       momento.String(CACHE_KEY),
	})
	if err != nil {
		panic(err)
	}

	switch resp.(type) {
	case *responses.GetHit:
		fmt.Printf("Cache Hit!\n")
	case *responses.GetMiss:
		fmt.Printf("Cache Miss!\n")
	}

	return "Success", nil
}

func getApiKeySecret(secretName string) (string, error) {
	if cachedAuthToken == "" {
		getSecretName, ok := os.LookupEnv(secretName)
		if !ok {
			fmt.Printf("Missing required env var '%s'\n", secretName)
			return "Missing required env var", nil
		}
		secret, err := secretsClient.GetSecretString(getSecretName)
		if err != nil {
			fmt.Printf("Unable to get secret '%s'\n", getSecretName)
			return "Error", err
		}
		cachedAuthToken = secret
	}

	return cachedAuthToken, nil
}

func getEndpointSecret(secretName string) (string, error) {
	if cachedEndpoint == "" {
		getSecretName, ok := os.LookupEnv(secretName)
		if !ok {
			fmt.Printf("Missing required env var '%s'\n", secretName)
			return "Missing required env var", nil
		}
		secret, err := secretsClient.GetSecretString(getSecretName)
		if err != nil {
			fmt.Printf("Unable to get secret '%s'\n", getSecretName)
			return "Error", err
		}
		cachedEndpoint = secret
	}

	return cachedEndpoint, nil
}

func getCacheClient() (momento.CacheClient, error) {
	if momentoCacheClient != nil {
		fmt.Println("Using cached Momento cache client")
		return momentoCacheClient, nil
	}

	authToken, secretErr := getApiKeySecret("MOMENTO_API_KEY_SECRET_NAME")
	if secretErr != nil {
		panic(secretErr)
	}

	endpoint, secretErr := getEndpointSecret("MOMENTO_ENDPOINT_SECRET_NAME")
	if secretErr != nil {
		panic(secretErr)
	}

	props := auth.ApiKeyV2Props{ApiKey: authToken, Endpoint: endpoint}
	credentialProvider, err := auth.NewApiKeyV2TokenProvider(props)
	if err != nil {
		panic(err)
	}

	newCacheClient, initErr := momento.NewCacheClientWithEagerConnectTimeout(
		config.LambdaLatest(),
		credentialProvider,
		60*time.Second,
		30*time.Second,
	)
	if initErr != nil {
		panic(initErr)
	}

	fmt.Println("New cache client created")
	momentoCacheClient = newCacheClient
	return momentoCacheClient, nil
}

// To measure the latency of 100 GET requests to your Momento cache,
// just call this code in the handler function
func measureLatency(cacheClient momento.CacheClient, ctx context.Context) {
	fmt.Printf("\nMeasuring GET request latency:\n")
	for i := 0; i < 100; i++ {
		start := time.Now()
		resp, _ := cacheClient.Get(ctx, &momento.GetRequest{
			CacheName: CACHE_NAME,
			Key:       momento.String(CACHE_KEY),
		})
		timeTaken := time.Since(start)

		switch resp.(type) {
		case *responses.GetHit:
			fmt.Printf("response %d: Hit! | time: %v\n", i, timeTaken)
		case *responses.GetMiss:
			fmt.Printf("response %d: Miss! | time: %v\n", i, timeTaken)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
