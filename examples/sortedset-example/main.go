package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const (
	cacheName = "test-cache"
	setName   = "my-set"
)

func main() {
	// Create Momento client
	client := getClient()
	ctx := context.Background()

	// Create cache
	setupCache(client, ctx)

	// Put score for each element to set
	// Using counter, element N has score N
	for i := 1; i < 11; i++ {
		_, err := client.SortedSetPut(ctx, &momento.SortedSetPutRequest{
			CacheName: cacheName,
			SetName:   setName,
			Elements: []*momento.SortedSetPutElement{{
				Value: momento.String(fmt.Sprintf("element-%d", i)),
				Score: float64(i),
			}},
		})
		if err != nil {
			panic(err)
		}
	}

	// Fetch sorted set
	fmt.Println("\n\nFetching all elements from sorted set:")
	fmt.Println("--------------")
	fetchResp, err := client.SortedSetFetch(ctx, &momento.SortedSetFetchRequest{
		CacheName: cacheName,
		SetName:   setName,
	})
	if err != nil {
		panic(err)
	}

	displayElements(setName, fetchResp)

	// Fetch elements in descending order (high -> low)
	fmt.Println("\n\nFetching all elements from sorted set in descending order:")
	fmt.Println("--------------")
	descendingResp, err := client.SortedSetFetch(ctx, &momento.SortedSetFetchRequest{
		CacheName: cacheName,
		SetName:   setName,
		Order:     momento.DESCENDING,
	})
	if err != nil {
		panic(err)
	}

	displayElements(setName, descendingResp)
}

func getClient() momento.CacheClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := momento.NewCacheClient(
		config.LaptopLatest(),
		credProvider,
		60*time.Second,
	)
	if err != nil {
		panic(err)
	}
	return client
}

func setupCache(client momento.CacheClient, ctx context.Context) {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "test-cache",
	})
	if err != nil {
		panic(err)
	}
}

func displayElements(setName string, resp responses.SortedSetFetchResponse) {
	switch r := resp.(type) {
	case *responses.SortedSetFetchHit:
		for _, e := range r.ValueStringElements() {
			fmt.Printf("setName: %s, value: %s, score: %f\n", setName, e.Value, e.Score)
		}
		fmt.Println("")
	case *responses.SortedSetFetchMiss:
		fmt.Println("we regret to inform you there is no such set")
		os.Exit(1)
	}
}
