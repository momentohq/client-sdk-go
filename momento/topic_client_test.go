package momento_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

var _ = Describe("topic-client", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCaches()

		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	DescribeTable("Validates the names",
		func(cacheName string, collectionName string, expectedError string) {
			ctx := sharedContext.Ctx
			client := sharedContext.TopicClient
			value := String("foo")

			Expect(
				client.Subscribe(ctx, &TopicSubscribeRequest{
					CacheName: cacheName, TopicName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.Publish(ctx, &TopicPublishRequest{
					CacheName: cacheName, TopicName: collectionName, Value: value,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))
		},
		Entry("Empty cache name", "", sharedContext.CollectionName, InvalidArgumentError),
		Entry("Blank cache name", "  ", sharedContext.CollectionName, InvalidArgumentError),
		Entry("Empty collection name", sharedContext.CacheName, "", InvalidArgumentError),
		Entry("Blank collection name", sharedContext.CacheName, "  ", InvalidArgumentError),
		Entry("Non-existent cache", uuid.NewString(), uuid.NewString(), CacheNotFoundError),
	)

	It(`Publishes and receives`, func() {
		publishedValues := []TopicValue{
			String("aaa"),
			Bytes([]byte{1, 2, 3}),
		}

		sub, err := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: sharedContext.CollectionName,
		})
		if err != nil {
			panic(err)
		}

		cancelContext, cancelFunction := context.WithCancel(sharedContext.Ctx)
		var receivedValues []TopicValue
		ready := make(chan int, 1)
		go func() {
			ready <- 1
			for {
				select {
				case <-cancelContext.Done():
					return
				default:
					value, err := sub.Item(sharedContext.Ctx)
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
			_, err := sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
				CacheName: sharedContext.CacheName,
				TopicName: sharedContext.CollectionName,
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

	It("Cancels the context immediataly after subscribing and asserts as such", func() {

		sub, _ := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: sharedContext.CollectionName,
		})

		// Create a new context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		done := make(chan bool)

		// Run Item function in a goroutine
		go func() {
			_, err := sub.Item(ctx)
			if err == nil {
				fmt.Println("Expected an error due to context cancellation, got nil")
			}
			close(done)
		}()

		// immediately cancel the context
		cancel()

		// Wait for either the Item function to return or the test to timeout
		select {
		case <-done:
			// Test passed
		case <-time.After(time.Second * 2):
			Fail("Test timed out, likely due to infinite loop in Item function")
		}

	})

	It("returns an error when trying to publish a nil topic value", func() {
		Expect(
			sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
				CacheName: sharedContext.CacheName,
				TopicName: sharedContext.CollectionName,
				Value:     nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	Describe(`Subscribe`, func() {
		It(`Does not error on a non-existent topic`, func() {
			Expect(
				sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: sharedContext.CacheName,
					TopicName: sharedContext.CollectionName,
				}),
			).Error().NotTo(HaveOccurred())
		})
	})

	It("Can close individual topics subscriptions without closing the grpc channel", func() {
		topic1 := fmt.Sprintf("golang-topics-test-%s", uuid.NewString())
		topic2 := fmt.Sprintf("golang-topics-test-%s", uuid.NewString())

		// subscribe to one topic
		sub1, err := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topic1,
		})
		if err != nil {
			panic(err)
		}

		// subscribe to another topic
		sub2, err := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topic2,
		})
		if err != nil {
			panic(err)
		}

		// publish messages to both
		_, err = sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topic1,
			Value:     String("hello-1"),
		})
		if err != nil {
			panic(err)
		}
		_, err = sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topic2,
			Value:     String("hello-2"),
		})
		if err != nil {
			panic(err)
		}

		// expect two Item() successes
		item, err := sub1.Item(sharedContext.Ctx)
		if err != nil {
			panic(err)
		}
		switch msg := item.(type) {
		case String:
			Expect(msg).To(Equal(String("hello-1")))
		case Bytes:
			Fail("Expected topic item to be a string")
		}

		item, err = sub2.Item(sharedContext.Ctx)
		if err != nil {
			panic(err)
		}
		switch msg := item.(type) {
		case String:
			Expect(msg).To(Equal(String("hello-2")))
		case Bytes:
			Fail("Expected topic item to be a string")
		}

		// close one subscription
		sub1.Close()

		// publish messages to both
		_, err = sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topic1,
			Value:     String("hello-again-1"),
		})
		if err != nil {
			panic(err)
		}

		_, err = sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topic2,
			Value:     String("hello-again-2"),
		})
		if err != nil {
			panic(err)
		}

		// expect one Item() success and one failure
		item, err = sub1.Item(sharedContext.Ctx)
		Expect(item).To(BeNil())
		Expect(err.Error()).To(Equal("context canceled"))

		item, err = sub2.Item(sharedContext.Ctx)
		if err != nil {
			panic(err)
		}
		switch msg := item.(type) {
		case String:
			Expect(msg).To(Equal(String("hello-again-2")))
		case Bytes:
			Fail("Expected topic item to be a string")
		}
	})
})
