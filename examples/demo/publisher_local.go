package main

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
)

func PublisherLocal() {
	ctx := context.Background()
	credProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	// Create Momento client
	client, err := incubating.NewScsClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
	})
	if err != nil {
		panic(err)
	}

	// Create cache
	createCacheIfNotExist(ctx, client, "default")

	// Start publishing events to the topic
	fmt.Println("Publishing topic: local-test-topic")
	for {
		err = client.PublishTopic(ctx, &incubating.TopicPublishRequest{
			CacheName: "default",
			TopicName: "local-test-topic",
			Value:     "TacoBell!",
		})
		if err != nil {
			panic(err)
		}
		// Send an event every 2 seconds = total of 30 events per minute
		time.Sleep(2 * time.Second)
	}
}
