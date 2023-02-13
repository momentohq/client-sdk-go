package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
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
		err := client.SortedSetPut(ctx, &incubating.SortedSetPutRequest{
			CacheName: cacheName,
			SetName:   setName,
			Elements: []*incubating.SortedSetScoreRequestElement{{
				Name:  momento.StringBytes{Text: fmt.Sprintf("element-%d", i)},
				Score: float64(i),
			}},
		})
		if err != nil {
			panic(err)
		}
	}

	// Fetch sorted set
	fetchResp, err := client.SortedSetFetch(ctx, &incubating.SortedSetFetchRequest{
		CacheName: cacheName,
		SetName:   setName,
	})
	if err != nil {
		panic(err)
	}

	// Display all elements in sorted set
	switch r := fetchResp.(type) {
	case *incubating.SortedSetFetchHit:
		fmt.Println("--------------")
		fmt.Println("Found sorted set with following elements:")
		for _, e := range r.Elements {
			fmt.Println(fmt.Sprintf("setName: %s elementName: %s score: %f", setName, e.Name, e.Score))
		}
	case *incubating.SortedSetFetchMiss:
		fmt.Println("we regret to inform you there is no such set")
		os.Exit(1)
	}

	// Fetch top 5 elements in descending order (high -> low)
	fmt.Println("--------------")
	fmt.Println("\n\nFetching Top 5 elements from sorted set:")
	top5Rsp, err := client.SortedSetFetch(ctx, &incubating.SortedSetFetchRequest{
		CacheName:       cacheName,
		SetName:         setName,
		NumberOfResults: incubating.FetchLimitedElements{Limit: 5},
		Order:           incubating.DESCENDING,
	})
	if err != nil {
		panic(err)
	}

	// Display top 5 elements using the result from SortedSetFetch in descending order
	switch r := top5Rsp.(type) {
	case *incubating.SortedSetFetchHit:
		for _, e := range r.Elements {
			fmt.Println(fmt.Sprintf("setName: %s elementName: %s score: %f", setName, e.Name, e.Score))
		}
		fmt.Println("\n")
	case *incubating.SortedSetFetchMiss:
		fmt.Println("we regret to inform you there is no such set")
		os.Exit(1)
	}

}

func getClient() incubating.ScsClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := incubating.NewScsClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
		DefaultTTL:         60 * time.Second,
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
