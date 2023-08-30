package batchutils_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/batchutils"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	. "github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Batch set operations", func() {

	var (
		ctx       context.Context
		client    CacheClient
		cacheName string
	)

	BeforeEach(func() {
		ctx = context.Background()
		cacheName = fmt.Sprintf("golang-%s", uuid.NewString())
		credentialProvider, err := auth.FromString("eyJlbmRwb2ludCI6ImNlbGwtNC11cy13ZXN0LTItMS5wcm9kLmEubW9tZW50b2hxLmNvbSIsImFwaV9rZXkiOiJleUpoYkdjaU9pSklVekkxTmlKOS5leUp6ZFdJaU9pSndjbUYwYVd0QWJXOXRaVzUwYjJoeExtTnZiU0lzSW5abGNpSTZNU3dpY0NJNklrTkJRVDBpTENKbGVIQWlPakUyT1RVNE1qQTFNREI5Lk50cVdwYVBpNnNBVmd1WFNlVDg4OFFVa3JDOW5ldVlqbFk2TXp4ZW9DUVkifQ==")
		if err != nil {
			panic(err)
		}
		client, err = NewCacheClient(
			config.LaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory()),
			credentialProvider,
			time.Second*60,
		)
		if err != nil {
			panic(err)
		}

		_, err = client.CreateCache(ctx, &CreateCacheRequest{CacheName: cacheName})
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		_, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: cacheName})
		if err != nil {
			panic(err)
		}
	})

	Describe("batch set", func() {
		It("happy path", func() {
			var batchSetKeys []Value
			var items []batchutils.BatchSetItem

			for i := 0; i < 10; i++ {
				key := String(fmt.Sprintf("MSETk%d", i))
				batchSetKeys = append(batchSetKeys, key)
				item := batchutils.BatchSetItem{
					Key:   key,
					Value: String(fmt.Sprintf("MSETv%d", i)),
					Ttl:   1 * time.Second,
				}
				items = append(items, item)
			}

			setBatch, _ := batchutils.BatchSet(ctx, &batchutils.BatchSetRequest{
				Client:    client,
				CacheName: cacheName,
				Items:     items,
			})

			setResponses := setBatch.Responses()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, resp := range setResponses {
				Expect(resp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			getBatch, _ := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      batchSetKeys,
			})

			getBatchResponses := getBatch.Responses()
			for i := 0; i < len(batchSetKeys); i++ {
				switch r := getBatchResponses[batchSetKeys[i]].(type) {
				case *responses.GetHit:
					Expect(r.ValueString()).To(Equal(fmt.Sprintf("MSETv%d", i)))
				case *responses.GetMiss:
					Fail("expected a hit but got a MISS")
				default:
					Fail(fmt.Sprintf("failed on %d", i))
				}
			}
		})

		It("happy path with some items without ttl", func() {
			var batchSetKeys []Value
			var items []batchutils.BatchSetItem

			for i := 0; i < 5; i++ {
				key := String(fmt.Sprintf("MSETk%d", i))
				batchSetKeys = append(batchSetKeys, key)
				item := batchutils.BatchSetItem{
					Key:   key,
					Value: String(fmt.Sprintf("MSETv%d", i)),
					Ttl:   1 * time.Second,
				}
				items = append(items, item)
			}

			for i := 5; i < 10; i++ {
				key := String(fmt.Sprintf("MSETk%d", i))
				batchSetKeys = append(batchSetKeys, key)
				// without TTL
				item := batchutils.BatchSetItem{
					Key:   key,
					Value: String(fmt.Sprintf("MSETv%d", i)),
				}
				items = append(items, item)
			}

			setBatch, _ := batchutils.BatchSet(ctx, &batchutils.BatchSetRequest{
				Client:    client,
				CacheName: cacheName,
				Items:     items,
			})

			setResponses := setBatch.Responses()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, resp := range setResponses {
				Expect(resp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			getBatch, _ := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      batchSetKeys,
			})

			getBatchResponses := getBatch.Responses()
			for i := 0; i < len(batchSetKeys); i++ {
				switch r := getBatchResponses[batchSetKeys[i]].(type) {
				case *responses.GetHit:
					Expect(r.ValueString()).To(Equal(fmt.Sprintf("MSETv%d", i)))
				case *responses.GetMiss:
					Fail("expected a hit but got a MISS")
				default:
					Fail(fmt.Sprintf("failed on %d", i))
				}
			}
		})

		It("some items without a value gives invalid argument error", func() {
			var items []batchutils.BatchSetItem

			var batchSetSuccessfulKeys []Value
			for i := 0; i < 5; i++ {
				key := String(fmt.Sprintf("MSETk%d", i))
				batchSetSuccessfulKeys = append(batchSetSuccessfulKeys, key)
				item := batchutils.BatchSetItem{
					Key:   key,
					Value: String(fmt.Sprintf("MSETv%d", i)),
					Ttl:   1 * time.Second,
				}
				items = append(items, item)
			}

			var batchSetErrorKeys []Value
			for i := 5; i < 10; i++ {
				key := String(fmt.Sprintf("MSETk%d", i))
				batchSetErrorKeys = append(batchSetErrorKeys, key)
				// without TTL
				item := batchutils.BatchSetItem{
					Key: key,
				}
				items = append(items, item)
			}

			setBatch, errors := batchutils.BatchSet(ctx, &batchutils.BatchSetRequest{
				Client:    client,
				CacheName: cacheName,
				Items:     items,
			})

			// Assuming errors is an instance of *BatchSetError
			Expect(len(errors.Errors())).To(Equal(len(batchSetErrorKeys)))
			for v, e := range errors.Errors() {
				Expect(e.Error()).To(Equal("InvalidArgumentError: value cannot be nil"))
				isValidKey := false
				for _, erroredKey := range batchSetErrorKeys {
					if v == erroredKey {
						isValidKey = true
						break
					}
				}
				if !isValidKey {
					Fail("Found successful key must in list of successful keys sent to server ")
				}
			}

			setResponses := setBatch.Responses()
			Expect(len(setResponses)).To(Equal(len(batchSetSuccessfulKeys)))
			for v, resp := range setResponses {
				isValidKey := false
				for _, successfullKey := range batchSetSuccessfulKeys {
					if v == successfullKey {
						isValidKey = true
						break
					}
				}
				if !isValidKey {
					Fail("Found unsuccessful key must in list of unsuccessful keys sent to server ")
				}
				Expect(resp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			// assert hits for successful items
			getBatch, _ := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      batchSetSuccessfulKeys,
			})

			getBatchResponses := getBatch.Responses()
			for i := 0; i < len(batchSetSuccessfulKeys); i++ {
				switch r := getBatchResponses[batchSetSuccessfulKeys[i]].(type) {
				case *responses.GetHit:
					Expect(r.ValueString()).To(Equal(fmt.Sprintf("MSETv%d", i)))
				case *responses.GetMiss:
					Fail("expected a hit but got a MISS")
				default:
					Fail(fmt.Sprintf("failed on %d", i))
				}
			}

			// assert misses for successful items
			getBatch, _ = batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      batchSetErrorKeys,
			})

			for _, resp := range getBatch.Responses() {
				Expect(resp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
			}
		})

		It("super small request timeout test", func() {
			var items []batchutils.BatchSetItem

			for i := 0; i < 10; i++ {
				key := String(fmt.Sprintf("MSETk%d", i))
				item := batchutils.BatchSetItem{
					Key:   key,
					Value: String(fmt.Sprintf("MSETv%d", i)),
					Ttl:   1 * time.Second,
				}
				items = append(items, item)
			}

			timeout := 1 * time.Millisecond
			setBatchResp, setErrors := batchutils.BatchSet(ctx, &batchutils.BatchSetRequest{
				Client:         client,
				CacheName:      cacheName,
				Items:          items,
				RequestTimeout: &timeout,
			})

			Expect(setBatchResp).To(BeNil())
			Expect(len(setErrors.Errors())).To(Equal(len(items)))
			for _, err := range setErrors.Errors() {
				Expect(err.Error()).To(ContainSubstring("TimeoutError: context deadline exceeded"))
			}

		})

	})
})
