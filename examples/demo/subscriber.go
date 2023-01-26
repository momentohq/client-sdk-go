package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/prozz/aws-embedded-metrics-golang/emf"
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
	client, err := incubating.NewScsClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
	})
	if err != nil {
		panic(err)
	}
	sub, err := client.SubscribeTopic(ctx, &incubating.TopicSubscribeRequest{
		CacheName: "default",
		TopicName: subscriberTopicName,
	})
	if err != nil {
		panic(err)
	}

	err = sub.Recv(context.Background(), func(ctx context.Context, m *incubating.TopicMessageReceiveResponse) {
		layout := "2006-01-02T15:04:05.000Z07:00"
		trimmedValue := strings.ReplaceAll(m.StringValue(), "text:", "")
		trimmedValue = strings.ReplaceAll(trimmedValue, "\"", "")
		receivedTime, err := time.Parse(layout, trimmedValue)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Received  time=%v", receivedTime))
		currentTime := time.Now().Format(layout)
		parsedCurrentTime, err := time.Parse(layout, currentTime)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Current  time=%v", parsedCurrentTime))
		// latency is in nanoseconds so dividing it by a million
		latency := parsedCurrentTime.Sub(receivedTime) / 1000000
		emf.New(emf.WithLogGroup("pubsub")).MetricAs("ReceivingMessageLatency", int(latency), emf.Milliseconds).Dimension("subscriber", "receiving").Log()
		fmt.Println(fmt.Sprintf("Received a message! latency=%dms", latency))
		fmt.Println()
	})
	if err != nil {
		panic(err)
	}
}
