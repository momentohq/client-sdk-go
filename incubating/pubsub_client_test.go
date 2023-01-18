package incubating

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	testPubSubAuthToken = os.Getenv("TEST_AUTH_TOKEN")
)

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathPubSub(t *testing.T) {
	ctx := context.Background()
	go func() {
		newMockPubSubServer()
	}()
	client, err := NewPubSubClient(testPubSubAuthToken) // TODO should we be returning error here?
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
