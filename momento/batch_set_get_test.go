package momento_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cache-client get-batch set-batch", Label(CACHE_SERVICE_LABEL), func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCaches()
		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	DescribeTable("validates cache name", func(cacheName string, expectedErrorCode string) {
		Expect(
			sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: cacheName,
				Items:     nil,
			}),
		).Error().To(HaveMomentoErrorCode(expectedErrorCode))

		Expect(
			sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: cacheName,
				Keys:      nil,
			}),
		).Error().To(HaveMomentoErrorCode(expectedErrorCode))
	},
		Entry("nonexistent cache name", uuid.NewString(), CacheNotFoundError),
		Entry("empty cache name", "", InvalidArgumentError),
		Entry("nil cache name", nil, InvalidArgumentError),
	)

	Describe("SetBatch", func() {
		It("GetBatch happy path with all misses", func() {
			var batchSetKeys []Value
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("Batch-key-%d", i)
				batchSetKeys = append(batchSetKeys, String(key))
			}

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))

			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))

			for _, getResp := range getResponses {
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
			}
		})

		It("SetBatch happy path", func() {
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				item := BatchSetItem{
					Key:   String(fmt.Sprintf("Batch-key-%d", i)),
					Value: String(fmt.Sprintf("Batch-value-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}
		})

		It("SetBatch happy path with TTL", func() {
			var batchSetKeys []Value
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("Batch-key-%d", i)
				batchSetKeys = append(batchSetKeys, String(key))
				item := BatchSetItem{
					Key:   String(key),
					Value: String(fmt.Sprintf("Batch-value-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       500 * time.Millisecond,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			time.Sleep(2 * time.Second)

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))
			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for _, getResp := range getResponses {
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
			}
		})

		It("GetBatch happy path with all hits", func() {
			var batchSetKeys []Value
			var batchSetKeysString []string
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("All-hits-%d", i)
				batchSetKeysString = append(batchSetKeysString, key)
				batchSetKeys = append(batchSetKeys, String(key))
				item := BatchSetItem{
					Key:   String(key),
					Value: String(fmt.Sprintf("All-hits-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       10 * time.Second,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))
			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for _, getResp := range getResponses {
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			}
			getValueMap := getBatchResp.(responses.GetBatchSuccess).ValueMap()
			for i := 0; i < len(batchSetKeysString); i++ {
				fetchedValue := getValueMap[batchSetKeysString[i]]
				Expect(fetchedValue).To(Equal(fmt.Sprintf("All-hits-%d", i)))
			}
		})

		It("GetBatch happy path with some hits and some misses", func() {
			var batchSetKeys []Value
			var batchSetKeysString []string
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("Some-hits-%d", i)

				if i < 5 {
					batchSetKeys = append(batchSetKeys, String(key))
					batchSetKeysString = append(batchSetKeysString, key)
				} else {
					differentKey := String(fmt.Sprintf("Some-hits-%d-miss", i))
					batchSetKeys = append(batchSetKeys, differentKey)
				}

				item := BatchSetItem{
					Key:   String(key),
					Value: String(fmt.Sprintf("Some-hits-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       10 * time.Second,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))
			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for i, getResp := range getResponses {
				if i < 5 {
					Expect(getResp).To(BeAssignableToTypeOf(&responses.GetHit{}))
					Expect(getResp.(*responses.GetHit).ValueString()).To(Equal(fmt.Sprintf("Some-hits-%d", i)))
				} else {
					Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
				}
			}

			getValuesMap := getBatchResp.(responses.GetBatchSuccess).ValueMap()
			// for each key, check if the value is as expected
			for i, key := range batchSetKeysString {
				value, ok := getValuesMap[key]
				if i < 5 {
					Expect(ok).To(BeTrue())
					Expect(value).To(Equal(fmt.Sprintf("Some-hits-%d", i)))
				} else {
					Expect(ok).To(BeFalse())
				}
			}
		})

	})
})
