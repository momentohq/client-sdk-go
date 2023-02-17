package momento_test

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
)

func getClient(client_props *SimpleCacheClientProps) SimpleCacheClient {
	credProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client_props.CredentialProvider = credProvider
	client, err := NewSimpleCacheClient(client_props)
	if err != nil {
		panic(err)
	}
	return client
}

func createCache(client SimpleCacheClient) string {
	ctx := context.Background()
	cacheName := "go-pubsub-" + uuid.NewString()

	_, err := client.CreateCache(ctx, &CreateCacheRequest{
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

	return cacheName
}

func deleteCache(client SimpleCacheClient, cacheName string) {
	ctx := context.Background()

	_, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

var _ = Describe("Pubsub", func() {
	var client SimpleCacheClient
	var cacheName string
	var ctx context.Context
	var topicName string

	BeforeEach(func() {
		ctx = context.Background()
		topicName = uuid.NewString()

		client = getClient(&SimpleCacheClientProps{
			Configuration: config.LatestLaptopConfig(),
			DefaultTTL:    60 * time.Second,
		})
		DeferCleanup(func() {
			client.Close()
		})

		cacheName = createCache(client)
		DeferCleanup(func() {
			deleteCache(client, cacheName)
		})
	})

	It(`Publishes and receives`, func() {
		publishedValues := []TopicValue{
			&TopicValueString{Text: "aaa"},
			&TopicValueBytes{Bytes: []byte{1, 2, 3}},
		}

		sub, err := client.TopicSubscribe(ctx, &TopicSubscribeRequest{
			CacheName: cacheName,
			TopicName: topicName,
		})
		if err != nil {
			panic(err)
		}

		cancelContext, cancelFunction := context.WithCancel(ctx)
		receivedValues := []TopicValue{}
		ready := make(chan int, 1)
		go func() {
			ready <- 1
			for {
				select {
				case <-cancelContext.Done():
					return
				default:
					value, err := sub.Item()
					if err != nil {
						panic(err)
					}
					receivedValues = append(receivedValues, value)
				}
			}
		}()
		<-ready

		time.Sleep(time.Millisecond * 100)
		for _, value := range publishedValues {
			_, err := client.TopicPublish(ctx, &TopicPublishRequest{
				CacheName: cacheName,
				TopicName: topicName,
				Value:     value,
			})
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Millisecond * 100)
		cancelFunction()

		Expect(receivedValues).To(Equal(publishedValues))
	})
})
