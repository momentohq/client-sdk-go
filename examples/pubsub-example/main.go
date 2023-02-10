package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
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
	sub, err := client.SubscribeTopic(ctx, &incubating.TopicSubscribeRequest{
		CacheName: cacheName,
		TopicName: topicName,
	})
	if err != nil {
		panic(err)
	}

	// Receive and print messages in a goroutine
	go func() { pollForMessages(sub) }()
	time.Sleep(time.Second)

	// Publish messages for the subscriber
	publishMessages(client, ctx)
}

func pollForMessages(sub incubating.SubscriptionIFace) {
	for {
		item, err := sub.Item()
		if err != nil {
			panic(err)
		}
		switch msg := item.(type) {
		case *incubating.TopicValueString:
			fmt.Printf("received message: '%s'\n", msg.Text)
		case *incubating.TopicValueBytes:
			fmt.Printf("received message: '%s'\n", msg.Bytes)
		}
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

func publishMessages(client incubating.ScsClient, ctx context.Context) {
	for i := 0; i < 10; i++ {
		fmt.Printf("publishing message %d\n", i)
		err := client.PublishTopic(ctx, &incubating.TopicPublishRequest{
			CacheName: cacheName,
			TopicName: topicName,
			Value: &incubating.TopicValueString{
				Text: fmt.Sprintf("hello %d", i),
			},
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
