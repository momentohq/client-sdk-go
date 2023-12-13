package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/momentohq/client-sdk-go/utils"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	auth_resp "github.com/momentohq/client-sdk-go/responses/auth"
)

var (
	ctx    context.Context
	client momento.CacheClient
)

func example_API_InstantiateCacheClient() {
	context := context.Background()
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	defaultTtl := time.Duration(9999)
	eagerConnectTimeout := 30 * time.Second

	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	if err != nil {
		panic(err)
	}

	client.Ping(context)
}

func example_API_ListCaches() {
	resp, err := client.ListCaches(ctx, &momento.ListCachesRequest{})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.ListCachesSuccess:
		log.Printf("Found caches %+v", r.Caches())
	}
}

func example_API_CreateCache() {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "cache-name",
	})
	if err != nil {
		panic(err)
	}
}

func example_API_Get() {
	key := uuid.NewString()
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.GetHit:
		log.Printf("Lookup resulted in cache HIT. value=%s\n", r.ValueString())
	case *responses.GetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}
}

func example_API_Set() {
	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err := client.Set(ctx, &momento.SetRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
		Value:     momento.String(value),
		Ttl:       time.Duration(9999),
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.TimeoutError {
				// this would represent a client-side timeout, and you could fall back to your original data source
			} else {
				panic(err)
			}
		}
	}
}

func example_API_Delete() {
	key := uuid.NewString()
	_, err := client.Delete(ctx, &momento.DeleteRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}
}

func example_API_InstantiateTopicClient() {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	topicClient, err := momento.NewTopicClient(
		config.TopicsDefault(),
		credProvider,
	)
	if err != nil {
		panic(err)
	}
}

func example_API_TopicPublish(client momento.TopicClient) {
	_, err := client.Publish(ctx, &momento.TopicPublishRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
		Value:     momento.String("test-message"),
	})
	if err != nil {
		panic(err)
	}
}

func example_API_TopicSubscribe(client momento.TopicClient) {
	// Instantiate subscriber
	sub, subErr := client.Subscribe(ctx, &momento.TopicSubscribeRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
	})
	if subErr != nil {
		panic(subErr)
	}

	time.Sleep(time.Second)
	_, pubErr := client.Publish(ctx, &momento.TopicPublishRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
		Value:     momento.String("test-message"),
	})
	if pubErr != nil {
		panic(pubErr)
	}
	time.Sleep(time.Second)

	item, err := sub.Item(ctx)
	if err != nil {
		panic(err)
	}
	switch msg := item.(type) {
	case momento.String:
		fmt.Printf("received message as string: '%v'\n", msg)
	case momento.Bytes:
		fmt.Printf("received message as bytes: '%v'\n", msg)
	}
}

func example_API_GenerateDisposableToken(client momento.AuthClient) {
	tokenId := "a token id"
	resp, err := client.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
		ExpiresIn: utils.ExpiresInSeconds(10),
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

	switch r := resp.(type) {
	case *auth_resp.GenerateDisposableTokenSuccess:
		log.Printf("Successfully generated a disposable token for endpoint=%s with tokenId=%s\n", r.Endpoint, tokenId)
	}
}
