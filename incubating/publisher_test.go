package incubating

import (
	"context"
	"errors"
	"github.com/momentohq/client-sdk-go/auth"
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
		newMomentoLocalTestServer(testPortToUse)
	}()
	client, err := newLocalScsClient(testPortToUse) // TODO should we be returning error here?
	if err != nil {
		panic(err)
	}
	for {
		err = client.PublishTopic(ctx, &TopicPublishRequest{
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
	client, err := NewScsClient(credProvider, 3600)
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
	for {
		err = client.PublishTopic(ctx, &TopicPublishRequest{
			TopicName: publisherTopicName,
			Value:     time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}
