package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/utils"
)

const (
	cacheName             = "my-test-cache"
	listName              = "my-test-list"
	itemDefaultTTLSeconds = 60
)

var (
	ctx    context.Context
	client momento.CacheClient
)

func pushFrontToList(value momento.Value) {
	fmt.Printf("\npushing '%s' to front of list\n", value)
	resp, err := client.ListPushFront(ctx, &momento.ListPushFrontRequest{
		CacheName:          cacheName,
		ListName:           listName,
		Value:              value,
		TruncateBackToSize: 0,
		Ttl: &utils.CollectionTtl{
			Ttl:        5 * time.Second,
			RefreshTtl: true,
		},
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *momento.ListPushFrontSuccess:
		fmt.Printf("pushed with 5 sec TTL to front of list whose length is now %d\n", r.ListLength())
	}
}

func pushBackToList(value momento.Value) {
	fmt.Printf("\npushing '%s' to back of list\n", value)
	resp, err := client.ListPushBack(ctx, &momento.ListPushBackRequest{
		CacheName:           cacheName,
		ListName:            listName,
		Value:               value,
		TruncateFrontToSize: 0,
		Ttl: &utils.CollectionTtl{
			Ttl:        5 * time.Second,
			RefreshTtl: true,
		},
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *momento.ListPushBackSuccess:
		fmt.Printf("pushed with 5 sec TTL to back of list whose length is now %d\n", r.ListLength())
	}
}

func printList() {
	resp, err := client.ListFetch(ctx, &momento.ListFetchRequest{
		CacheName: cacheName,
		ListName:  listName,
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *momento.ListFetchHit:
		fmt.Printf("\nlist fetch returned:\n\n\t%s\n", strings.Join(r.ValueListString(), "\n\t"))
	case *momento.ListFetchMiss:
		fmt.Println("\nlist fetch returned a MISS")
	}
}

func printListLength() {
	resp, err := client.ListLength(ctx, &momento.ListLengthRequest{
		CacheName: cacheName,
		ListName:  listName,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *momento.ListLengthMiss:
		fmt.Println("\nlist length returned a MISS")
	case *momento.ListLengthHit:
		fmt.Printf("\ngot list length: %d", r.Length())
	}
}

func concatFront(values []momento.Value) {
	resp, err := client.ListConcatenateFront(ctx, &momento.ListConcatenateFrontRequest{
		CacheName: cacheName,
		ListName:  listName,
		Values:    values,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *momento.ListConcatenateFrontSuccess:
		fmt.Printf("\nconcatenated values to front. list is now length %d\n", r.ListLength())
	}
}

func concatBack(values []momento.Value) {
	resp, err := client.ListConcatenateBack(ctx, &momento.ListConcatenateBackRequest{
		CacheName: cacheName,
		ListName:  listName,
		Values:    values,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *momento.ListConcatenateBackSuccess:
		fmt.Printf("\nconcatenated values to back. list is now length %d\n", r.ListLength())
	}
}

func removeValue(value momento.Value) {
	_, err := client.ListRemoveValue(ctx, &momento.ListRemoveValueRequest{
		CacheName: cacheName,
		ListName:  listName,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx = context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	// Initializes Momento
	client, err = momento.NewCacheClient(
		config.LatestLaptopConfig(),
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
	)
	if err != nil {
		panic(err)
	}

	// Create Cache and check if CacheName exists
	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}

	printListLength()

	for i := 0; i < 5; i++ {
		stringValue := fmt.Sprintf("hello #%d", i+1)
		pushFrontToList(momento.String(stringValue))
	}

	printListLength()
	printList()

	time.Sleep(time.Second * 5)

	for i := 0; i < 5; i++ {
		stringValue := fmt.Sprintf("hello #%d", i+1)
		pushBackToList(momento.String(stringValue))
	}

	printListLength()
	printList()

	time.Sleep(time.Second * 5)

	printListLength()
	printList()

	for i := 0; i < 5; i++ {
		stringValue := fmt.Sprintf("hello #%d", i+1)
		pushFrontToList(momento.String(stringValue))
	}
	for i := 0; i < 5; i++ {
		resp, err := client.ListPopFront(ctx, &momento.ListPopFrontRequest{
			CacheName: cacheName,
			ListName:  listName,
		})
		if err != nil {
			panic(err)
		}
		switch r := resp.(type) {
		case *momento.ListPopFrontHit:
			fmt.Printf("\npopped value '%s'\n", r.ValueString())
		case *momento.ListPopFrontMiss:
			fmt.Println("\npop from front returned MISS")
		}
	}
	printListLength()

	pushFrontToList(momento.String("list seed"))

	var values []momento.Value
	for i := 0; i < 5; i++ {
		values = append(values, momento.String(fmt.Sprintf("concat front %d", i)))
	}
	concatFront(values)
	printList()

	values = nil
	for i := 0; i < 5; i++ {
		values = append(values, momento.String(fmt.Sprintf("concat back %d", i)))
	}
	concatBack(values)
	printList()

	_, err = client.Delete(ctx, &momento.DeleteRequest{
		CacheName: cacheName,
		Key:       momento.String(listName),
	})
	if err != nil {
		panic(err)
	}

	for i := 1; i < 11; i++ {
		if i%2 != 0 {
			pushBackToList(momento.String("odd"))
		} else {
			pushBackToList(momento.String("even"))
		}
	}
	printList()
	value := "even"
	removeValue(momento.String(value))
	fmt.Printf("\nremoved '%s' from list\n", value)
	printList()

	// Delete the cache
	if _, err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	fmt.Printf("\ndeleted cache\n")
}
