package incubating

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBasicHappyPathSubscriber(t *testing.T) {
	ctx := context.Background()
	testPortToUse := 3000
	client, err := NewLocalPubSubClient(testPortToUse) // TODO should we be returning error here?
	if err != nil {
		panic(err)
	}
	sub, err := client.SubscribeTopic(ctx, &TopicSubscribeRequest{
		TopicName: "test-topic",
	})
	if err != nil {
		panic(err)
	}

	err = sub.Recv(context.Background(), func(ctx context.Context, m *TopicMessageReceiveResponse) {
		layout := "2006-01-02T15:04:05.000Z07:00"
		trimmedValue := strings.ReplaceAll(m.StringValue(), "text:", "")
		trimmedValue = strings.ReplaceAll(trimmedValue, "\"", "")
		receivedTime, err := time.Parse(layout, trimmedValue)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Received  time=%v", receivedTime))
		currentTime := time.Now().Format(layout)
		parsedCurrentTime, err := time.Parse(layout, currentTime)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Current  time=%v", parsedCurrentTime))
		// latency is in nanoseconds so dividing it by a million
		latency := parsedCurrentTime.Sub(receivedTime)
		fmt.Println(fmt.Sprintf("Received a message! latency=%dms", latency/1000000))
		fmt.Println()
	})
	if err != nil {
		panic(err)
	}
}
