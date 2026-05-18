package momento_test

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TopicSubscription", func() {
	Describe("pool slot release", func() {
		It("Subscribe + Close cycles do not exhaust the pool", func() {
			// Pin to a single stream channel so pool capacity is exactly 100 slots,
			// then run more Subscribe+Close cycles than the capacity. Each cycle
			// must return its slot, or the loop fails with ClientResourceExhaustedError
			// once the leaked slots fill the pool.
			cfg := config.TopicsDefault().WithNumStreamGrpcChannels(1)
			topicClient, err := NewTopicClient(cfg, sharedContext.CredentialProvider)
			Expect(err).NotTo(HaveOccurred())
			defer topicClient.Close()

			for i := 0; i < 150; i++ {
				sub, err := topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: sharedContext.CacheName,
					TopicName: fmt.Sprintf("leak-regression-%d", i),
				})
				Expect(err).NotTo(HaveOccurred(), "Subscribe failed at iteration %d", i)
				Expect(sub).NotTo(BeNil())
				sub.Close()
			}
		})

		It("Subscribe + Event + Close cycles release exactly one slot per subscription", func() {
			// Same setup as above, but with an Event() call after Close. Event
			// observes the canceled context and would normally decrement the slot;
			// Close already did. The two paths must coexist without double-counting,
			// otherwise the pool counter drifts and either exhausts early or goes
			// negative. With idempotent release, the pool returns to zero each
			// iteration and all cycles succeed.
			cfg := config.TopicsDefault().WithNumStreamGrpcChannels(1)
			topicClient, err := NewTopicClient(cfg, sharedContext.CredentialProvider)
			Expect(err).NotTo(HaveOccurred())
			defer topicClient.Close()

			for i := 0; i < 150; i++ {
				sub, err := topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: sharedContext.CacheName,
					TopicName: fmt.Sprintf("idempotency-%d", i),
				})
				Expect(err).NotTo(HaveOccurred(), "Subscribe failed at iteration %d", i)
				Expect(sub).NotTo(BeNil())
				sub.Close()
				_, eventErr := sub.Event(sharedContext.Ctx)
				Expect(eventErr).To(HaveOccurred(),
					"Event after Close should return the canceled-context error at iteration %d", i)
			}
		})
	})
})
