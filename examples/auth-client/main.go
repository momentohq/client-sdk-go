package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	auth_responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
)

const (
	cacheName           = "go-auth-example"
	defaultTtl          = 60 * time.Second
	eagerConnectTimeout = 30 * time.Second
	timeUntilExpiry     = 10 * time.Second
)

func generateRefreshApiKey(authClient momento.AuthClient, ctx context.Context) {
	//generate an api key
	resp, err := authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInSeconds(10),
		Scope:     momento.AllDataReadWrite,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\ngenerated new API Key")

	//use it to create cache client and make some cache calls
	success := resp.(*auth_responses.GenerateApiKeySuccess)
	credentialProvider, err := auth.FromString(success.ApiKey)
	if err != nil {
		panic(err)
	}
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	if err != nil {
		panic(err)
	}
	key := uuid.NewString()
	value := uuid.NewString()
	_, err = client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(value),
	})
	if err != nil {
		panic(err)
	}

	// refresh the api key before it expires
	time.Sleep(2 * time.Second)
	refreshAuthClient, err := momento.NewAuthClient(config.AuthDefault(), credentialProvider)
	if err != nil {
		panic(err)
	}
	refreshResp, err := refreshAuthClient.RefreshApiKey(ctx, &momento.RefreshApiKeyRequest{
		RefreshToken: success.RefreshToken,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nrefreshed API Key")
	refreshSuccess := refreshResp.(*auth_responses.RefreshApiKeySuccess)
	credentialProvider, err = auth.FromString(refreshSuccess.ApiKey)
	if err != nil {
		panic(err)
	}

	//using the old key won't work after its expiration, but creating a new cache client using the new key returned by refreshApiKey will work
	time.Sleep(timeUntilExpiry)
	_, err = client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err == nil {
		fmt.Println("\nAPI Key has not expired.")
	} else {
		fmt.Println("\nFailed to get due to expired API key")
	}
	refreshClient, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	_, err = refreshClient.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nGet Request succeeded with new API key")
}

func generateDisposableToken(authClient momento.AuthClient, ctx context.Context) {
	//generate disposable token
	tokenId := "a token id"
	tokenResp, err := authClient.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
		ExpiresIn: utils.ExpiresInSeconds(9),
		Scope:     momento.AllDataReadWrite,
		Props: momento.DisposableTokenProps{
			TokenId: &tokenId,
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("\ngenerated new Disposable Token")

	//use it to create cache client and make some cache calls
	tokenSuccess := tokenResp.(*auth_responses.GenerateDisposableTokenSuccess)
	credentialProvider, err := auth.FromDisposableToken(tokenSuccess.ApiKey)
	if err != nil {
		panic(err)
	}
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	if err != nil {
		panic(err)
	}
	key := uuid.NewString()
	value := uuid.NewString()
	_, err = client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(value),
	})
	if err != nil {
		panic(err)
	}

	//generate new disposable ticket before the first expires
	tokenId = "a token id"
	tokenResp, err = authClient.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
		ExpiresIn: utils.ExpiresInSeconds(20),
		Scope:     momento.AllDataReadWrite,
		Props: momento.DisposableTokenProps{
			TokenId: &tokenId,
		},
	})
	fmt.Println("\ngenerated new Disposable Token")
	tokenSuccess = tokenResp.(*auth_responses.GenerateDisposableTokenSuccess)
	credentialProvider, err = auth.FromDisposableToken(tokenSuccess.ApiKey)
	if err != nil {
		panic(err)
	}

	//using the old token won't work after its expiration, but the new token will work
	time.Sleep(timeUntilExpiry)
	_, err = client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err == nil {
		fmt.Println("\nAPI Key has not expired.")
	} else {
		fmt.Println("\nFailed to get due to expired token")
	}
	refreshClient, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	_, err = refreshClient.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nGet Request succeeded with new token")
}

func main() {
	ctx := context.Background()
	// must be a super user API Key
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("V1_API_KEY")
	if err != nil {
		panic(err)
	}
	var authClient momento.AuthClient
	authClient, err = momento.NewAuthClient(config.AuthDefault(), credentialProvider)
	if err != nil {
		panic(err)
	}
	// Initializes Momento
	cacheClient, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}
	// Create Cache
	_, err = cacheClient.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	//Generate API Key and Refresh API Key
	generateRefreshApiKey(authClient, ctx)

	//Generate New Disposable Token Before First Expires
	generateDisposableToken(authClient, ctx)
}
