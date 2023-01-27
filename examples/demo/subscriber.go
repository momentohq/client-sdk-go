package main

import (
	"context"
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
		receivedTime, err := strconv.Atoi(m.StringValue())
		if err != nil {
			panic(err)
		}
		latency := currentTime - receivedTime
		// send metrics to CloudWatch
		emf.New(emf.WithLogGroup("pubsub")).MetricAs("ReceivingMessageLatency", latency, emf.Milliseconds).Dimension("subscriber", "receiving").Log()
	})
	if err != nil {
		panic(err)
	}
}
