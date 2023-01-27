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
		currentTime := int(time.Now().UnixMilli())
		receivedValue := m.StringValue()
		// for this demo, check if the value is not empty. Handle Discontinuity later
		if len(receivedValue) == 0 {
			fmt.Println("discontinuity is detected.")
		} else {
			receivedTime, err := strconv.Atoi(receivedValue)
			if err != nil {
				panic(err)
			}
			latency := currentTime - receivedTime
			// send metrics to CloudWatch
			emf.New(emf.WithLogGroup("pubsub")).MetricAs("ReceivingMessageLatency", latency, emf.Milliseconds).DimensionSet(emf.NewDimension("subscriber", "receiving"), emf.NewDimension("taskId", time.RFC3339)).Log()
		}
	})
	if err != nil {
		panic(err)
	}
}
