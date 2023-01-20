package incubating

import (
	"context"
	"testing"
	"time"
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
			TopicName: "test-topic",
			Value:     time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}
