package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	cacheName = "cache"
	listName  = "my-list"
)

func main() {
	// Initialization
	client := getClient()
	ctx := context.Background()
	setupCache(client, ctx)

	resp, err := client.ListFetch(ctx, &incubating.ListFetchRequest{
		CacheName: cacheName,
		ListName:  listName,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *incubating.ListFetchHit:
		fmt.Println(strings.Join(r.ValueListString(), ", "))
	case *incubating.ListFetchMiss:
		fmt.Println("we regret to inform you there is no such list")
		os.Exit(1)
	}

	lenResp, err := client.ListLength(ctx, &incubating.ListLengthRequest{
		CacheName: cacheName,
		ListName:  listName,
	})
	if err != nil {
		panic(err)
	}
	switch r := lenResp.(type) {
	case *incubating.ListLengthSuccess:
		fmt.Printf("list %s is length %d\n", listName, int(r.Length()))
	}
}

func getClient() incubating.ScsClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := incubating.NewScsClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
	})
	if err != nil {
		panic(err)
	}
	return client
}

func setupCache(client momento.ScsClient, ctx context.Context) {
	err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "test-cache",
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.AlreadyExistsError {
				panic(err)
			}
		}
	}
}
