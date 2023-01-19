package incubating

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
)

var (
	authToken = os.Getenv("TEST_AUTH_TOKEN")
)

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathLocalPubSub(t *testing.T) {
	ctx := context.Background()
	testPortToUse := 3000
	go func() {
		newMockPubSubServer(testPortToUse)
	}()

	client, err := NewLocalPubSubClient(testPortToUse)
	if err != nil {
		panic(err)
	}

	sub, err := client.SubscribeTopic(ctx, &TopicSubscribeRequest{
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	go func() {
		err = sub.Recv(context.Background(), func(ctx context.Context, m *TopicMessageReceiveResponse) {
			fmt.Println(fmt.Sprintf("got a msg! val=%s", m.StringValue()))
		})
		if err != nil {
			panic(err)
		}
	}()

	for i := 0; i < 10; i++ {
		err = client.PublishTopic(ctx, &TopicPublishRequest{
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

	client, err := NewPubSubClient(authToken)
	if err != nil {
		panic(err)
	}
	err = client.CreateTopic(ctx, &CreateTopicRequest{
		TopicName: "test-topic",
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
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	//go func() {
	// Just block and make sure we get stubbed messages for now for quick test
	err = sub.Recv(context.Background(), func(ctx context.Context, m *TopicMessageReceiveResponse) {
		fmt.Println(fmt.Sprintf("got a msg! val=%s", m.StringValue()))
	})
	if err != nil {
		panic(err)
	}
	//}()

	// TODO remote api doesnt support publish yet
	//for i := 0; i < 10; i++ {
	//	err = client.PublishTopic(ctx, &TopicPublishRequest{
	//		TopicName: "test-topic",
	//		Value:     fmt.Sprintf("hello %d", i),
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//	time.Sleep(time.Second)
	//}
}
