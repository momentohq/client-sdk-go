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
	subscriberTopicName = os.Getenv("TEST_TOPIC_NAME")
)

func Subscriber() {
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

	// Subscribe to a topic
	sub, err := client.SubscribeTopic(ctx, &incubating.TopicSubscribeRequest{
		CacheName: "default",
		TopicName: subscriberTopicName,
	})
	if err != nil {
		panic(err)
	}

	// Kick off goroutine to send subscriber count to CloudWatch
	sendSubscriberCountToCw()

	// Start receiving events
	err = sub.Recv(context.Background(), func(ctx context.Context, m *incubating.TopicMessageReceiveResponse) {
		currentTime := int(time.Now().UnixMilli())
		publishedTime, err := strconv.Atoi(m.StringValue())
		if err != nil {
			fmt.Printf("Received non-time value: %s\n", m.StringValue())
		}
		latency := currentTime - publishedTime
		// Send metrics to CloudWatch
		sendLatencyToCw(latency)
	})
	if err != nil {
		panic(err)
	}
}
