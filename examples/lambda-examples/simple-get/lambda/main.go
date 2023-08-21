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
	cachedAuthToken     string = ""
	cachedMomentoClient momento.CacheClient
	secretsClient, _    = secretcache.New()
)

func handler() (string, error) {
	cacheClient, err := getCacheClient()
	if err != nil {
		return "Something went wrong getting the cache client", err
	}
	cachedMomentoClient = cacheClient
	ctx := context.Background()

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

func getSecret(secretName string) (string, error) {
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

func getCacheClient() (momento.CacheClient, error) {
	if cachedMomentoClient == nil {
		authToken, secretErr := getSecret("MOMENTO_AUTH_TOKEN_SECRET_NAME")
		if secretErr != nil {
			panic(secretErr)
		}

		credentialProvider, err := auth.NewStringMomentoTokenProvider(authToken)
		if err != nil {
			panic(err)
		}

		_cacheClient, initErr := momento.NewCacheClient(
			config.LaptopLatest(),
			credentialProvider,
			60*time.Second,
		)
		if initErr != nil {
			panic(initErr)
		}

		cachedMomentoClient = _cacheClient
	}
	return cachedMomentoClient, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
