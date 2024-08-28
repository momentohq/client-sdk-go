package main

import (
	"context"
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

// You wouldn't need the cacheClient object, setupCache function, or associated code if you were connecting to an existing cache in Momento Serverless Cache.
// For demonstration purposes, this example creates a cache for the new Topic to show the publish and subscribe functionality.
func main() {
	// Initialization
	topicClient := getTopicClient()
	cacheClient := getCacheClient()
	ctx := context.Background()
	setupCache(cacheClient, ctx)

	// Instantiate subscriber
	sub, err := topicClient.Subscribe(ctx, &momento.TopicSubscribeRequest{
		CacheName: cacheName,
		TopicName: topicName,
	})
	if err != nil {
		panic(err)
	}

	// Receive and print messages in a goroutine
	go func() { pollForMessages(ctx, sub) }()

	// Receive and print all events in a goroutine
	// go func() { pollForEvents(ctx, sub) }()

	time.Sleep(time.Second)

	// Publish messages for the subscriber
	publishMessages(topicClient, ctx)

	sub.Close()
}

func pollForMessages(ctx context.Context, sub momento.TopicSubscription) {
	for {
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
}

func pollForEvents(ctx context.Context, sub momento.TopicSubscription) {
	for {
		event, err := sub.Event(ctx)
		if err != nil {
			panic(err)
		}
		switch e := event.(type) {
		case momento.TopicHeartbeat:
			fmt.Printf("received heartbeat\n")
		case momento.TopicDiscontinuity:
			fmt.Printf("received discontinuity\n")
		case momento.TopicItem:
			fmt.Printf(
				"received message with sequence number %d: %v \n",
				e.GetTopicSequenceNumber(),
				e.GetValue(),
			)
		}
	}
}

func getTopicClient() momento.TopicClient {
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
	return topicClient
}

func getCacheClient() momento.CacheClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	cacheClient, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credProvider,
		time.Second*60,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}
	return cacheClient
}

func setupCache(client momento.CacheClient, ctx context.Context) {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "test-cache",
	})
	if err != nil {
		panic(err)
	}
}

func publishMessages(client momento.TopicClient, ctx context.Context) {
	for i := 0; i < 5; i++ {
		fmt.Printf("publishing message %d\n", i)
		_, err := client.Publish(ctx, &momento.TopicPublishRequest{
			CacheName: cacheName,
			TopicName: topicName,
			Value:     momento.String(fmt.Sprintf("hello %d", i)),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("publishing message %d\n", i)
		_, err := client.Publish(ctx, &momento.TopicPublishRequest{
			CacheName: cacheName,
			TopicName: topicName,
			Value:     momento.Bytes(fmt.Sprintf("hello %d", i)),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
