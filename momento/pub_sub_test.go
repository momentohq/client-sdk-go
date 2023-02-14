package momento

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
)

var client ScsClient
var cacheName = os.Getenv("TEST_CACHE_NAME")

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
	client, err := NewSimpleCacheClient(&SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credProvider,
		DefaultTTL:         60 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return *client
}

func setup() {
	ctx := context.Background()
	client = getClient()
	err := client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != AlreadyExistsError {
				panic(err)
			}
		}
	}
}

func teardown() {
	client.Close()
}

func publishTopic(ctx context.Context, pubClient ScsClient, i int) {
	var topicVal TopicValue

	if i%2 == 0 {
		topicVal = &TopicValueString{Text: "hello txt"}
	} else {
		topicVal = &TopicValueBytes{Bytes: []byte("hello bytes")}
	}

	_, err := pubClient.TopicPublish(ctx, &TopicPublishRequest{
		CacheName: cacheName,
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
	cancelContext, cancelFunction := context.WithCancel(ctx)

	sub, err := client.TopicSubscribe(ctx, &TopicSubscribeRequest{
		CacheName: cacheName,
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	numMessagesToSend := 10
	numMessagesReceived := 0
	go func() {
		for {
			select {
			case <-cancelContext.Done():
				return
			default:
				_, err := sub.Item()
				if err != nil {
					panic(err)
				}
				numMessagesReceived++
			}
		}
	}()
	time.Sleep(time.Second)

	for i := 0; i < numMessagesToSend; i++ {
		publishTopic(ctx, client, i)
		time.Sleep(time.Second)
	}
	cancelFunction()

	if numMessagesReceived != numMessagesToSend {
		t.Errorf("expected %d messages but received %d", numMessagesToSend, numMessagesReceived)
	}
}
