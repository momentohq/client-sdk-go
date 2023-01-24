package incubating

import (
	"context"
	"os"
	"testing"
	"time"
)

var (
	publisherTopicName = os.Getenv("TEST_TOPIC_NAME")
)

func TestBasicHappyPathPublisher(t *testing.T) {
	ctx := context.Background()
	testPortToUse := 3000
	go func() {
		newMockPubSubServer(testPortToUse)
	}()
	client, err := NewLocalPubSubClient(testPortToUse) // TODO should we be returning error here?
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
