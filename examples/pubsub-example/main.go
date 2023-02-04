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
	//cancelContext, cancelFunction := context.WithCancel(ctx)
	// Or you can set a timeout after which the goroutine will be cancelled
	//cancelContext, cancelFunction := context.WithTimeout(ctx, time.Second*10)
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
	//go func() { pollForMessages(sub, cancelContext) }()
	go func() {
		for i := 0; i < 10; i++ {
			msg := popMessage(sub)
			fmt.Println(msg)
			time.Sleep(time.Second * 5)
		}
	}()

	// Publish messages and then shut down the subscriber goroutine
	publishMessages(client, ctx)
	//cancelFunction()
	// Prove that the goroutine is stopped by publishing more messages that
	// won't be output to the console
	publishMessages(client, ctx)
	time.Sleep(time.Second * 60)
}

func getClient() incubating.ScsClient {
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

func pollForMessages(sub incubating.SubscriptionIFace, cancelContext context.Context) {
	err := sub.Consume(cancelContext, func(ctx context.Context, m incubating.TopicValue) {
		switch msg := m.(type) {
		case *incubating.TopicValueString:
			fmt.Printf("received message: '%s'\n", msg.Text)
		case *incubating.TopicValueBytes:
			fmt.Printf("received message: '%s'\n", msg.Bytes)
		}
	})
	if err != nil {
		panic(err)
	}
}

func popMessage(sub incubating.SubscriptionIFace) string {
	ctx := context.Background()
	var msgOut string
	err := sub.Recv(ctx, func(ctx context.Context, m incubating.TopicValue) {
		switch msg := m.(type) {
		case *incubating.TopicValueString:
			msgOut = msg.Text
		case *incubating.TopicValueBytes:
			msgOut = string(msg.Bytes)
		}
	})
	if err != nil {
		panic(err)
	}
	return msgOut
}
