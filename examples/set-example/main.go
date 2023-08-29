package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const (
	cacheName             = "my-test-cache"
	setName               = "my-test-set"
	itemDefaultTTLSeconds = 60
)

var (
	ctx    context.Context
	client momento.CacheClient
)

func setup() {
	ctx = context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	// Initializes Momento
	client, err = momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}

	// Create Cache
	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}

func addElement(element momento.Value) {
	_, err := client.SetAddElement(ctx, &momento.SetAddElementRequest{
		CacheName: cacheName,
		SetName:   setName,
		Element:   element,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nadded element to set")
}

func addElements(elements []momento.Value) {
	_, err := client.SetAddElements(ctx, &momento.SetAddElementsRequest{
		CacheName: cacheName,
		SetName:   setName,
		Elements:  elements,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nadded elements to set")
}

func removeElement(element momento.Value) {
	_, err := client.SetRemoveElement(ctx, &momento.SetRemoveElementRequest{
		CacheName: cacheName,
		SetName:   setName,
		Element:   element,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nelement removed from set")
}

func removeElements(elements []momento.Value) {
	_, err := client.SetRemoveElements(ctx, &momento.SetRemoveElementsRequest{
		CacheName: cacheName,
		SetName:   setName,
		Elements:  elements,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nelements removed from set")
}

func printSet() {
	resp, err := client.SetFetch(ctx, &momento.SetFetchRequest{
		CacheName: cacheName,
		SetName:   setName,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.SetFetchHit:
		fmt.Printf("\nprinting set elements:\n\t%s\n", strings.Join(r.ValueString(), "\n\t"))
	case *responses.SetFetchMiss:
		fmt.Println("set fetch returned a MISS")
	}
}

func cleanUp() {
	_, err := client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func main() {
	setup()

	elements := []momento.Value{
		momento.String("element the first"),
		momento.String("element the second"),
		momento.String("element the third"),
	}

	addElements(elements)
	printSet()

	addElement(momento.String("one"))
	addElement(momento.String("at"))
	addElement(momento.String("a"))
	addElement(momento.String("time"))
	printSet()

	removeElement(momento.String("at"))
	removeElement(momento.String("a"))
	printSet()

	removeElements(elements)
	printSet()

	cleanUp()
}
