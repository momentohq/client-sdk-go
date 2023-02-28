package momento_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

var _ = Describe("Pubsub", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCache()

		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	DescribeTable(`Validates the names`,
		func(cacheName string, collectionName string, expectedError string) {
			ctx := sharedContext.Ctx
			client := sharedContext.Client
			value := TopicValueString("foo")

			Expect(
				client.TopicSubscribe(ctx, &TopicSubscribeRequest{
					CacheName: cacheName, TopicName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.TopicPublish(ctx, &TopicPublishRequest{
					CacheName: cacheName, TopicName: collectionName, Value: value,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))
		},
		Entry("Empty cache name", "", sharedContext.CollectionName, InvalidArgumentError),
		Entry("Blank cache name", "  ", sharedContext.CollectionName, InvalidArgumentError),
		Entry("Empty collection name", sharedContext.CacheName, "", InvalidArgumentError),
		Entry("Blank collection name", sharedContext.CacheName, "  ", InvalidArgumentError),
		Entry("Non-existent cache", uuid.NewString(), uuid.NewString(), NotFoundError),
	)

	It(`Publishes and receives`, func() {
		publishedValues := []TopicValue{
			TopicValueString("aaa"),
			TopicValueBytes([]byte{1, 2, 3}),
		}

		sub, err := sharedContext.Client.TopicSubscribe(sharedContext.Ctx, &TopicSubscribeRequest{
			CacheName: sharedContext.CacheName,
			TopicName: sharedContext.CollectionName,
		})
		if err != nil {
			panic(err)
		}

		cancelContext, cancelFunction := context.WithCancel(sharedContext.Ctx)
		receivedValues := []TopicValue{}
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
			_, err := sharedContext.Client.TopicPublish(sharedContext.Ctx, &TopicPublishRequest{
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

	Describe(`TopicSubscribe`, func() {
		It(`Does not error on a non-existent topic`, func() {
			Expect(
				sharedContext.Client.TopicSubscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: sharedContext.CacheName,
					TopicName: sharedContext.CollectionName,
				}),
			).Error().NotTo(HaveOccurred())
		})
	})
})
