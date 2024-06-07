package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"math/rand"
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
	logLevel := momento_default_logger.INFO
	// for more verbose logging, comment out the previous line and uncomment the line below
	//logLevel := momento_default_logger.DEBUG
	loggerFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(logLevel)
	log := loggerFactory.GetLogger("sharded-topic-example")

	rand.Seed(time.Now().UnixNano())

	topicClient := getTopicClient(loggerFactory)
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
	go func() { pollForMessages(ctx, sub, log) }()
	time.Sleep(time.Second)

	// Publish messages for the subscriber
	publishMessages(topicClient, ctx, log)

	sub.Close()
}

func pollForMessages(ctx context.Context, sub momento.TopicSubscription, log logger.MomentoLogger) {
	for {
		item, err := sub.Item(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Debug("context canceled, shutting down subscription poll loop")
				return
			}
			panic(err)
		}
		switch msg := item.(type) {
		case momento.String:
			log.Info(fmt.Sprintf("received message as string: '%v'", msg))
		case momento.Bytes:
			log.Info(fmt.Sprintf("received message as bytes: '%v'", msg))
		}
	}
}

func getTopicClient(loggerFactory logger.MomentoLoggerFactory) momento.TopicClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	topicClient, err := NewShardedTopicClient(
		config.TopicsDefaultWithLogger(loggerFactory),
		credProvider,
		16,
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

func publishMessages(client momento.TopicClient, ctx context.Context, log logger.MomentoLogger) {
	for i := 0; i < 10; i++ {
		log.Info(fmt.Sprintf("publishing message %d", i))
		_, err := client.Publish(ctx, &momento.TopicPublishRequest{
			CacheName: cacheName,
			TopicName: topicName,
			Value:     momento.String(fmt.Sprintf("hello %d", i)),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
