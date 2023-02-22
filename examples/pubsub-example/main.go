package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	cacheName = "test-cache"
	topicName = "test-topic"
)

func main() {
	// Initialization
	client := getClient()
	ctx := context.Background()
	setupCache(client, ctx)

	// Instantiate subscriber
	sub, err := client.TopicSubscribe(ctx, &momento.TopicSubscribeRequest{
		CacheName: cacheName,
		TopicName: topicName,
	})
	if err != nil {
		panic(err)
	}

	// Receive and print messages in a goroutine
	go func() { pollForMessages(ctx, sub) }()
	time.Sleep(time.Second)

	// Publish messages for the subscriber
	publishMessages(client, ctx)
}

func pollForMessages(ctx context.Context, sub momento.TopicSubscription) {
	for {
		item, err := sub.Item(ctx)
		if err != nil {
			panic(err)
		}
		switch msg := item.(type) {
		case *momento.TopicValueString:
			fmt.Printf("received message: '%s'\n", msg.Text)
		case *momento.TopicValueBytes:
			fmt.Printf("received message: '%s'\n", msg.Bytes)
		}
	}
}

func getClient() momento.SimpleCacheClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
		DefaultTTL:         60 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return client
}

func setupCache(client momento.SimpleCacheClient, ctx context.Context) {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
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

func publishMessages(client momento.SimpleCacheClient, ctx context.Context) {
	for i := 0; i < 10; i++ {
		fmt.Printf("publishing message %d\n", i)
		_, err := client.TopicPublish(ctx, &momento.TopicPublishRequest{
			CacheName: cacheName,
			TopicName: topicName,
			Value: &momento.TopicValueString{
				Text: fmt.Sprintf("hello %d", i),
			},
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
