package main

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	auth_responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
)

const (
	cacheName             = "go-perf-test"
	itemDefaultTTLSeconds = 60
	requestTimeoutSeconds = 600 // 10 minutes
	maxRequestsPerSecond  = 10_000
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	var authClient momento.AuthClient
	authClient, err = momento.NewAuthClient(config.AuthDefault(), credentialProvider)
	if err != nil {
		panic(err)
	}
	//Generate API Key and Refresh API Key
	resp, err := authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(2),
		Scope:     momento.AllDataReadWrite,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\ngenerated new API Key")
	success := resp.(*auth_responses.GenerateApiKeySuccess)
	time.Sleep(5 * time.Second)

	_, err = authClient.RefreshApiKey(ctx, &momento.RefreshApiKeyRequest{
		RefreshToken: success.RefreshToken,
	})
	fmt.Println("\nrefreshed API Key")

	//Generate Disposable Token
	tokenId := "a token id"
	tokenResp, err := authClient.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
		ExpiresIn: utils.ExpiresInSeconds(5),
		Scope: momento.TopicSubscribeOnly(
			momento.CacheName{Name: "a cache"},
			momento.TopicName{Name: "a topic"},
		),
		Props: momento.DisposableTokenProps{
			TokenId: &tokenId,
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("\ngenerated new Disposable Token")
	tokenSuccess := tokenResp.(*auth_responses.GenerateDisposableTokenSuccess)
	_, err = auth.FromString(tokenSuccess.ApiKey)
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	//Generate New Disposable Token Before First Expires
	tokenId = "a token id"
	tokenResp, err = authClient.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
		ExpiresIn: utils.ExpiresInSeconds(10),
		Scope: momento.TopicSubscribeOnly(
			momento.CacheName{Name: "a cache"},
			momento.TopicName{Name: "a topic"},
		),
		Props: momento.DisposableTokenProps{
			TokenId: &tokenId,
		},
	})
	fmt.Println("\ngenerated new Disposable Token")
	tokenSuccess = tokenResp.(*auth_responses.GenerateDisposableTokenSuccess)
	_, err = auth.FromString(tokenSuccess.ApiKey)
	if err != nil {
		panic(err)
	}
}
