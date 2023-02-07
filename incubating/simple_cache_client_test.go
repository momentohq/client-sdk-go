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

var client ScsClient

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func getClient() ScsClient {
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
	return client
}

func setup() {
	ctx := context.Background()
	client = getClient()
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

func teardown() {
	client.Close()
}

func publishTopic(pubClient ScsClient, i int, ctx context.Context) {
	var topicVal TopicValue

	if i%2 == 0 {
		topicVal = &TopicValueString{Text: "hello txt"}
	} else {
		topicVal = &TopicValueBytes{Bytes: []byte("hello bytes")}
	}

	err := pubClient.PublishTopic(ctx, &TopicPublishRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
		Value:     topicVal,
	})
	if err != nil {
		panic(err)
	}
}

// Basic happy path test using a context which we cancel
func TestHappyPathPubSub(t *testing.T) {
	ctx := context.Background()

	sub, err := client.SubscribeTopic(ctx, &TopicSubscribeRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	numMessagesToSend := 10
	// TODO: use a channel instead of a counter variable
	numMessagesReceived := 0
	go func() {
		for {
			_, err := sub.Item()
			if err != nil {
				panic(err)
			}
			numMessagesReceived++
			if numMessagesToSend == numMessagesReceived {
				return
			}
		}
	}()
	time.Sleep(time.Second)

	for i := 0; i < numMessagesToSend; i++ {
		publishTopic(client, i, ctx)
		time.Sleep(time.Second)
	}

	// if we have received more than numMessagesToSend, our cancel failed
	if numMessagesReceived != numMessagesToSend {
		t.Errorf("expected no more than %d messages but received %d", numMessagesToSend, numMessagesReceived)
	}
}

// Basic happy path test using local test server
// TODO: are we going to keep the local client and server around?
func TestBasicHappyPathLocalPubSub(t *testing.T) {
	ctx := context.Background()
	testPortToUse := 3000
	go func() {
		newMomentoLocalTestServer(testPortToUse)
	}()

	localClient, err := newLocalScsClient(testPortToUse)
	if err != nil {
		panic(err)
	}

	sub, err := localClient.SubscribeTopic(ctx, &TopicSubscribeRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	numMessagesReceived := 0
	numMessagesToSend := 10
	go func() {
		for {
			_, err := sub.Item()
			if err != nil {
				panic(err)
			}
			numMessagesReceived++
			if numMessagesToSend == numMessagesReceived {
				return
			}
		}
	}()

	for i := 0; i < numMessagesToSend; i++ {
		var topicVal TopicValue
		if i%2 == 0 {
			topicVal = &TopicValueString{
				Text: fmt.Sprintf("string hello %d", i),
			}
		} else {
			topicVal = &TopicValueBytes{
				Bytes: []byte(fmt.Sprintf("byte hello %d", i)),
			}
		}
		err := localClient.PublishTopic(ctx, &TopicPublishRequest{
			CacheName: "test-cache",
			TopicName: "test-topic",
			Value:     topicVal,
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	if numMessagesToSend != numMessagesReceived {
		t.Errorf("expected %d messages but got %d", numMessagesToSend, numMessagesReceived)
	}
}
