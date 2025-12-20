package momento_test

import (
	"fmt"
	"time"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cache-client v2 api key tests", Label(CACHE_SERVICE_LABEL), func() {

	Describe("control plane", func() {
		It("creates, lists, and deletes caches", func() {
			cacheNames := []string{NewRandomString(), NewRandomString()}
			defer func() {
				for _, cacheName := range cacheNames {
					_, err := sharedContext.CacheClientApiKeyV2.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: cacheName})
					if err != nil {
						panic(err)
					}
				}
			}()

			for _, cacheName := range cacheNames {
				Expect(
					sharedContext.CacheClientApiKeyV2.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&CreateCacheSuccess{}))

				Expect(
					sharedContext.CacheClientApiKeyV2.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&CreateCacheAlreadyExists{}))
			}

			resp, err := sharedContext.CacheClientApiKeyV2.ListCaches(sharedContext.Ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())

			var listedCaches []string
			switch r := resp.(type) {
			case *ListCachesSuccess:
				for _, info := range r.Caches() {
					listedCaches = append(listedCaches, info.Name())
				}
				Expect(listedCaches).To(ContainElements(cacheNames))
			default:
				Fail("Unexpected response type")
			}

			for _, cacheName := range cacheNames {
				Expect(
					sharedContext.CacheClientApiKeyV2.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
			}
			resp, err = sharedContext.CacheClientApiKeyV2.ListCaches(sharedContext.Ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())
			Expect(resp).To(BeAssignableToTypeOf(&ListCachesSuccess{}))
			switch r := resp.(type) {
			case *ListCachesSuccess:
				Expect(r.Caches()).To(Not(ContainElements(cacheNames)))
			default:
				Fail("Unexpected response type")
			}
		})
	})

	Describe("get set delete", func() {
		It("Gets, Sets, and Deletes", func() {
			cacheName := sharedContext.DefaultCacheName
			client := sharedContext.CacheClientApiKeyV2
			key := NewRandomMomentoString()
			expectedString := NewRandomString()
			expectedBytes := []byte(expectedString)
			value := String(expectedString)

			Expect(
				client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: cacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				fmt.Println(`Expected GetHit, got:`, result)
				Fail("Unexpected type from Get")
			}

			Expect(
				client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&DeleteSuccess{}))

			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
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

			setBatchResp, setBatchErr := sharedContext.CacheClientApiKeyV2.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       10 * time.Second,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(SetBatchSuccess{}))
			setResponses := setBatchResp.(SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&SetSuccess{}))
			}

			getBatchResp, getBatchErr := sharedContext.CacheClientApiKeyV2.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(GetBatchSuccess{}))
			getResponses := getBatchResp.(GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for i, getResp := range getResponses {
				if i < 5 {
					Expect(getResp).To(BeAssignableToTypeOf(&GetHit{}))
					Expect(getResp.(*GetHit).ValueString()).To(Equal(fmt.Sprintf("Some-hits-%d", i)))
				} else {
					Expect(getResp).To(BeAssignableToTypeOf(&GetMiss{}))
				}
			}

			getValuesMap := getBatchResp.(GetBatchSuccess).ValueMap()
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

	// TODO: add remaining tests

	// Describe("compare and set", func() {})

	// Describe("other cache methods", func() {
	// 	It("KeysExist", func() {})

	// 	It("UpdateTtl", func() {})

	// 	It("Increment", func() {})

	// 	It("ItemGetType", func() {})

	// 	It("ItemGetTtl", func() {})
	// })

	// Describe("dictionary", func() {})

	// Describe("list", func() {})

	// Describe("set", func() {})

	// Describe("sorted set", func() {})

	// Describe("topics", func() {})

	// Describe("leaderboard", func() {})

})
