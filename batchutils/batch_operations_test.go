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

var _ = Describe("batch-utils", Label("cache-service"), func() {

	var (
		ctx       context.Context
		client    CacheClient
		cacheName string
		keys      []Value
	)

	BeforeEach(func() {
		ctx = context.Background()
		cacheName = fmt.Sprintf("golang-%s", uuid.NewString())
		//lint:ignore SA1019 // Still supporting FromEnvironmentVariable for backwards compatibility
		credentialProvider, err := auth.FromEnvironmentVariable("MOMENTO_API_KEY")
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

		for i := 0; i < 50; i++ {
			key := String(fmt.Sprintf("key%d", i))
			keys = append(keys, key)
			_, err := client.Set(ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String(fmt.Sprintf("val%d", i)),
			})
			if err != nil {
				panic(err)
			}
		}
	})

	AfterEach(func() {
		_, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: cacheName})
		if err != nil {
			panic(err)
		}
	})

	Describe("batch delete", func() {
		It("batch delete happy path", func() {
			errors := batchutils.BatchDelete(ctx, &batchutils.BatchDeleteRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      keys[5:21],
			})
			Expect(errors).To(BeNil())
			for i := 0; i < 50; i++ {
				resp, err := client.Get(ctx, &GetRequest{
					CacheName: cacheName,
					Key:       keys[i],
				})
				Expect(err).To(BeNil())
				switch resp.(type) {
				case *responses.GetHit:
					if i >= 5 && i <= 20 {
						Fail("got a hit for #%d that should be a miss", i)
					}
				case *responses.GetMiss:
					if !(i >= 5 && i <= 20) {
						Fail("got a miss for #%d that should be a hit", i)
					}
				}
			}
		})

		It("doesn't error trying to batch delete nonexistent keys", func() {
			keys := []Key{String("i"), String("don't"), String("exist")}
			errors := batchutils.BatchDelete(ctx, &batchutils.BatchDeleteRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      keys,
			})
			Expect(errors).To(BeNil())
		})

		It("super small request timeout test", func() {
			timeout := 1 * time.Nanosecond
			errors := batchutils.BatchDelete(ctx, &batchutils.BatchDeleteRequest{
				Client:         client,
				CacheName:      cacheName,
				Keys:           keys[5:21],
				RequestTimeout: &timeout,
			})

			Expect(len(errors.Errors())).To(Equal(16))
			for _, err := range errors.Errors() {
				Expect(err.Error()).To(ContainSubstring("TimeoutError: context deadline exceeded"))
			}
		})
	})

	Describe("batch get", func() {
		It("batch get happy path", func() {
			getBatch, errors := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      keys[5:21],
			})
			Expect(errors).To(BeNil())
			Expect(len(getBatch.Responses())).To(Equal(16))
			getResponses := getBatch.Responses()
			for i := 5; i < 21; i++ {
				switch r := getResponses[keys[i]].(type) {
				case *responses.GetHit:
					Expect(r.ValueString()).To(Equal(fmt.Sprintf("val%d", i)))
				case *responses.GetMiss:
					Fail("expected a hit but got a MISS")
				default:
					Fail(fmt.Sprintf("failed on %d", i))
				}
			}
		})

		It("returns misses for nonexistent keys", func() {
			keys := []Key{String("i"), String("don't"), String("exist")}
			getBatch, errors := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      keys,
			})
			Expect(errors).To(BeNil())
			getResponses := getBatch.Responses()
			for _, resp := range getResponses {
				Expect(resp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
			}
		})

		It("returns a mix of hits and misses", func() {
			keys := []Key{String("i"), String("key1"), String("don't"), String("exist"), String("key12")}
			getBatch, errors := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:    client,
				CacheName: cacheName,
				Keys:      keys,
			})
			Expect(errors).To(BeNil())
			getResponses := getBatch.Responses()
			for k, resp := range getResponses {
				if k == String("key1") || k == String("key12") {
					Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
				} else {
					Expect(resp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
				}
			}
		})

		It("super small request timeout test", func() {
			timeout := 1 * time.Nanosecond
			getBatch, errors := batchutils.BatchGet(ctx, &batchutils.BatchGetRequest{
				Client:         client,
				CacheName:      cacheName,
				Keys:           keys[5:21],
				RequestTimeout: &timeout,
			})

			Expect(getBatch).To(BeNil())
			Expect(len(errors.Errors())).To(Equal(16))
			for _, err := range errors.Errors() {
				Expect(err.Error()).To(ContainSubstring("TimeoutError: context deadline exceeded"))
			}
		})
	})
})
