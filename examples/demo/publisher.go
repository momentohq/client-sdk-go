package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
)

var (
	publisherTopicName = os.Getenv("TEST_TOPIC_NAME")
)

func Publisher() {
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
	fmt.Println(fmt.Sprintf("Publishing topic: %s", publisherTopicName))
	for {
		err = client.PublishTopic(ctx, &incubating.TopicPublishRequest{
			CacheName: "default",
			TopicName: publisherTopicName,
			Value:     strconv.Itoa(int(time.Now().UnixMilli())),
		})
		if err != nil {
			panic(err)
		}
		// Send an event every 2 seconds = total of 30 events per minute
		time.Sleep(2 * time.Second)
	}
}
