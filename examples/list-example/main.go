package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	cacheName     = "cache"
	listName      = "my-list"
	itemBatchSize = 3
)

func printList(client incubating.ScsClient, ctx context.Context) {
	fetchResp, err := client.ListFetch(ctx, &incubating.ListFetchRequest{
		CacheName: cacheName,
		ListName:  listName,
	})
	if err != nil {
		panic(err)
	}
	switch r := fetchResp.(type) {
	case *incubating.ListFetchHit:
		fmt.Println(strings.Join(r.ValueListString(), ", "))
	case *incubating.ListFetchMiss:
		fmt.Printf("no such list: %s\n", listName)
		os.Exit(1)
	}
}

func printListLength(client incubating.ScsClient, ctx context.Context) {
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

func pushFront(client incubating.ScsClient, ctx context.Context, value []byte) {
	pushFrontResp, err := client.ListPushFront(ctx, &incubating.ListPushFrontRequest{
		CacheName: cacheName,
		ListName:  listName,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}

	switch r := pushFrontResp.(type) {
	case *incubating.ListPushFrontSuccess:
		fmt.Printf("pushed value %s to list with length %d\n", value, r.ListLength())
	}
}

func pushBack(client incubating.ScsClient, ctx context.Context, value []byte) {
	pushBackResp, err := client.ListPushBack(ctx, &incubating.ListPushBackRequest{
		CacheName: cacheName,
		ListName:  listName,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}

	switch r := pushBackResp.(type) {
	case *incubating.ListPushBackSuccess:
		fmt.Printf("pushed value %s to list with length %d\n", value, r.ListLength())
	}
}

func popFront(client incubating.ScsClient, ctx context.Context) {
	popFrontResp, err := client.ListPopFront(ctx, &incubating.ListPopFrontRequest{
		ListName:  listName,
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	switch r := popFrontResp.(type) {
	case *incubating.ListPopFrontHit:
		fmt.Printf("popped value from front of list: %s\n", r.ValueString())
	case *incubating.ListPopFrontMiss:
		fmt.Println("got a miss response in response to attempt to pop value from front")
	}
}

func popBack(client incubating.ScsClient, ctx context.Context) {
	popBackResp, err := client.ListPopBack(ctx, &incubating.ListPopBackRequest{
		ListName:  listName,
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	switch r := popBackResp.(type) {
	case *incubating.ListPopBackHit:
		fmt.Printf("popped value from front of list: %s\n", r.ValueString())
	case *incubating.ListPopBackMiss:
		fmt.Println("got a miss response in response to attempt to pop value from front")
	}
}

func main() {
	// Initialization
	client := getClient()
	ctx := context.Background()
	setupCache(client, ctx)

	for i := 0; i < itemBatchSize; i++ {
		value := []byte(fmt.Sprintf("push front numero %d!", i+1))
		pushFront(client, ctx, value)
	}

	printList(client, ctx)
	printListLength(client, ctx)

	for i := 0; i < itemBatchSize; i++ {
		value := []byte(fmt.Sprintf("push back numero %d!", i+1))
		pushBack(client, ctx, value)
	}

	printList(client, ctx)
	printListLength(client, ctx)

	for i := 0; i < itemBatchSize; i++ {
		popFront(client, ctx)
		printListLength(client, ctx)
		popBack(client, ctx)
		printListLength(client, ctx)
		printList(client, ctx)
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
