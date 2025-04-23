package momento_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
)

var _ = Describe("topic-client", Label(TOPICS_SERVICE_LABEL), func() {
	var topicName string

	BeforeEach(func() {
		topicName = uuid.NewString()
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
		Entry("Empty cache name", "", topicName, InvalidArgumentError),
		Entry("Blank cache name", "  ", topicName, InvalidArgumentError),
		Entry("Empty collection name", sharedContext.CacheName, "", InvalidArgumentError),
		Entry("Blank collection name", sharedContext.CacheName, "  ", InvalidArgumentError),
		Entry("Non-existent cache", uuid.NewString(), uuid.NewString(), CacheNotFoundError),
	)

	It("Publishes and receives", func() {
		publishedValues := []TopicValue{
			String("aaa"),
			Bytes([]byte{1, 2, 3}),
		}

		sub, err := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topicName,
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

	It("should error on deadline exceeded", func() {
		newGrpcConfig := sharedContext.TopicConfiguration.GetTransportStrategy().GetGrpcConfig()
		newCfg := sharedContext.TopicConfiguration.WithTransportStrategy(
			sharedContext.TopicConfiguration.GetTransportStrategy().WithGrpcConfig(newGrpcConfig.WithClientTimeout(1)))
		newTopicClient, err := NewTopicClient(newCfg, sharedContext.CredentialProvider)
		if err != nil {
			panic(err)
		}

		_, err = newTopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
			CacheName: sharedContext.CacheName,
			TopicName: "topic",
			Value:     String("hi"),
		})
		Expect(err).To(HaveMomentoErrorCode(TimeoutError))
	})

	It("Publishes and receives detailed subscription items", func() {
		publishedValues := []TopicValue{
			String("aaa"),
			Bytes([]byte{1, 2, 3}),
		}

		// Value does not implement TopicValue
		expectedValues := []Value{
			String("aaa"),
			Bytes([]byte{1, 2, 3}),
		}

		sub, err := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topicName,
		})
		if err != nil {
			panic(err)
		}

		cancelContext, cancelFunction := context.WithCancel(sharedContext.Ctx)
		var receivedItems []TopicEvent
		ready := make(chan int, 1)
		go func() {
			ready <- 1
			for {
				select {
				case <-cancelContext.Done():
					return
				default:
					item, err := sub.Event(sharedContext.Ctx)
					if err != nil {
						panic(err)
					}
					receivedItems = append(receivedItems, item)
				}
			}
		}()
		<-ready

		time.Sleep(time.Millisecond * 100)
		for _, value := range publishedValues {
			_, err := sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
				CacheName: sharedContext.CacheName,
				TopicName: topicName,
				Value:     value,
			})
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second * 10)
		cancelFunction()

		Expect(len(receivedItems)).To(BeNumerically(">=", len(publishedValues)+1)) // +1 for the heartbeat(s)

		numberOfHeartbeats := 0
		numberOfDiscontinuities := 0

		// Collect only TopicItem events for comparison
		receivedTopicItems := []TopicItem{}
		for _, receivedItem := range receivedItems {
			switch r := receivedItem.(type) {
			case TopicItem:
				receivedTopicItems = append(receivedTopicItems, r)
			case TopicHeartbeat:
				numberOfHeartbeats++
			case TopicDiscontinuity:
				numberOfDiscontinuities++
			}
		}

		// Ensure we have received the expected number of TopicItems
		Expect(len(receivedTopicItems)).To(Equal(len(expectedValues)))

		// Now, compare the received TopicItems to the expected values
		for i, r := range receivedTopicItems {
			Expect(r.GetValue()).To(Equal(expectedValues[i]))
			Expect(r.GetTopicSequenceNumber()).To(BeNumerically(">", 0))
			Expect(r.GetTopicSequenceNumber()).To(Equal(uint64(i + 1)))
		}

		Expect(numberOfHeartbeats).To(BeNumerically(">=", 1))
		Expect(numberOfDiscontinuities).To(Equal(0))
	})

	It("Cancels the context immediately after subscribing and asserts as such", func() {

		sub, _ := sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: topicName,
		})

		// immediately cancel the context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		done := make(chan bool)

		// Run Item function in a goroutine
		go func() {
			_, err := sub.Item(ctx)
			if err == nil {
				fmt.Println("Expected an error due to context cancellation, got nil")
			}
			close(done)
		}()

		// Wait for either the Item function to return or the test to timeout
		select {
		case <-done:
			// Test passed
			Succeed()
		case <-time.After(time.Second * 5):
			Fail("Test timed out, likely due to infinite loop in Item function")
		}

	})

	It("returns an error when trying to publish a nil topic value", func() {
		Expect(
			sharedContext.TopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
				CacheName: sharedContext.CacheName,
				TopicName: topicName,
				Value:     nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	Describe(`Subscribe`, func() {
		It(`Does not error on a non-existent topic`, func() {
			Expect(
				sharedContext.TopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: sharedContext.CacheName,
					TopicName: topicName,
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

	Describe("Dynamic transport strategy", func() {
		It("Publishes and receives when using dynamic transport strategy", func() {
			publishedValues := []TopicValue{
				String("aaa"),
				Bytes([]byte{1, 2, 3}),
			}

			sub, err := sharedContext.DynamicTopicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
				CacheName: sharedContext.CacheName,
				TopicName: topicName,
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
				_, err := sharedContext.DynamicTopicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: sharedContext.CacheName,
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
})
