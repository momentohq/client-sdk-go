package momento_test

import (
	"context"
	"time"

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

	It(`Publishes and receives`, func() {
		publishedValues := []TopicValue{
			&TopicValueString{Text: "aaa"},
			&TopicValueBytes{Bytes: []byte{1, 2, 3}},
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
})
