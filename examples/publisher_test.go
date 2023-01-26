package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/incubating"
	"github.com/momentohq/client-sdk-go/momento"
	"os"
	"testing"
	"time"
)

var (
	publisherTopicName = os.Getenv("TEST_TOPIC_NAME")
)

func TestLocalBasicHappyPathPublisher(t *testing.T) {
	ctx := context.Background()
	testPortToUse := 3000
	go func() {
		incubating.NewLocalScsClient(testPortToUse)
	}()
	client, err := incubating.NewLocalScsClient(testPortToUse) // TODO should we be returning error here?
	if err != nil {
		panic(err)
	}
	for {
		err = client.PublishTopic(ctx, &incubating.TopicPublishRequest{
			TopicName: publisherTopicName,
			Value:     time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func TestBasicHappyPathPublisher(t *testing.T) {
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
	err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "default",
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.AlreadyExistsError {
				panic(err)
			}
		}
	}
	fmt.Println(fmt.Sprintf("Publishing topic: %s", publisherTopicName))
	for {
		err = client.PublishTopic(ctx, &incubating.TopicPublishRequest{
			CacheName: "default",
			TopicName: publisherTopicName,
			Value:     time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}
