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

	data := map[string]float64{
		"element1":  1.0,
		"element2":  2.0,
		"element3":  3.0,
		"element4":  4.0,
		"element5":  5.0,
		"element6":  6.0,
		"element7":  7.0,
		"element8":  8.0,
		"element9":  9.0,
		"element10": 10.0,
	}

	// Put score for each element to set
	// Using counter, element N has score N

	_, err := client.SortedSetPutElements(ctx, &momento.SortedSetPutElementsRequest{
		CacheName: cacheName,
		SetName:   setName,
		Elements:  momento.SortedSetElementsFromMap(data),
	})
	if err != nil {
		panic(err)
	}

	// Fetch sorted set
	fmt.Println("\n\nFetching all elements from sorted set:")
	fmt.Println("--------------")
	fetchResp, err := client.SortedSetFetchByRank(ctx, &momento.SortedSetFetchByRankRequest{
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
	descendingResp, err := client.SortedSetFetchByRank(ctx, &momento.SortedSetFetchByRankRequest{
		CacheName: cacheName,
		SetName:   setName,
		Order:     momento.DESCENDING,
	})
	if err != nil {
		panic(err)
	}

	displayElements(setName, descendingResp)

	// Clean up the cache
	_, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func getClient() momento.CacheClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credProvider,
		60*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}
	return client
}

func setupCache(client momento.CacheClient, ctx context.Context) {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
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
