package incubating

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathLocalPubSub(t *testing.T) {
	ctx := context.Background()
	testPortToUse := 3000
	go func() {
		newMomentoLocalTestServer(testPortToUse)
	}()

	client, err := newLocalScsClient(testPortToUse)
	if err != nil {
		panic(err)
	}

	sub, err := client.SubscribeTopic(ctx, &TopicSubscribeRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	go func() {
		err := sub.Recv(context.Background(), func(ctx context.Context, m *TopicMessageReceiveResponse) {
			fmt.Printf("got a msg! val=%s\n", m.StringValue())
		})
		if err != nil {
			panic(err)
		}
	}()

	for i := 0; i < 10; i++ {
		err := client.PublishTopic(ctx, &TopicPublishRequest{
			CacheName: "test-cache",
			TopicName: "test-topic",
			Value:     fmt.Sprintf("hello %d", i),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathPubSubIntegrationTest(t *testing.T) {
	ctx := context.Background()
	credProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := NewScsClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
	})
	if err != nil {
		panic(err)
	}
	err = client.CreateCache(ctx, &momento.CreateCacheRequest{
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

	sub, err := client.SubscribeTopic(ctx, &TopicSubscribeRequest{
		CacheName: "default",
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	go func() {
		// Just block and make sure we get stubbed messages for now for quick test
		err := sub.Recv(context.Background(), func(ctx context.Context, m *TopicMessageReceiveResponse) {
			fmt.Printf("got a msg! val=%s\n", m.StringValue())
		})
		if err != nil {
			panic(err)
		}
	}()

	for i := 0; i < 10; i++ {
		err := client.PublishTopic(ctx, &TopicPublishRequest{
			CacheName: "default",
			TopicName: "test-topic",
			Value:     fmt.Sprintf("hello %d", i),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
