package momento

import (
	"context"
	"fmt"
	"os"
	"testing"
)

var (
	testPubSubAuthToken = os.Getenv("TEST_AUTH_TOKEN")
)

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathPubSub(t *testing.T) {
	go func() {
		newMockPubSubServer()
	}()
	client, err := NewPubSubClient(testPubSubAuthToken) // TODO should we be returning error here?
	if err != nil {
		panic(err)
	}
	sub, err := client.SubscribeTopic(&TopicSubscribeRequest{
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	err = sub.Recv(context.Background(), func(ctx context.Context, m *TopicMessageReceiveResponse) {
		fmt.Println(fmt.Sprintf("got a msg! val=%s", m.StringValue()))
	})
	if err != nil {
		panic(err)
	}

}
