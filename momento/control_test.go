package momento_test

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
)

var _ = Describe("Control ops", func() {
	var client SimpleCacheClient
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()

		client = getClient(&SimpleCacheClientProps{
			Configuration: config.LatestLaptopConfig(),
			DefaultTTL:    60 * time.Second,
		})

		DeferCleanup(func() { client.Close() })
	})

	Describe(`Happy Path`, func() {
		It(`creates, lists, and deletes caches`, func() {
			cacheNames := []string{uuid.NewString(), uuid.NewString()}
			defer func() {
				for _, cacheName := range cacheNames {
					_, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: cacheName})
					if err != nil {
						panic(err)
					}
				}
			}()

			for _, cacheName := range cacheNames {
				Expect(
					client.CreateCache(ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&CreateCacheSuccess{}))

				Expect(
					client.CreateCache(ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&CreateCacheAlreadyExists{}))
			}

			resp, err := client.ListCaches(ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())

			listedCaches := []string{}
			switch r := resp.(type) {
			case *ListCachesSuccess:
				for _, info := range r.Caches() {
					listedCaches = append(listedCaches, info.Name())
				}
				Expect(listedCaches).To(ContainElements(cacheNames))
			default:
				Fail("Unexpected repsonse type")
			}

			for _, cacheName := range cacheNames {
				Expect(
					client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
			}
			resp, err = client.ListCaches(ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())
			Expect(resp).To(BeAssignableToTypeOf(&ListCachesSuccess{}))
			switch r := resp.(type) {
			case *ListCachesSuccess:
				Expect(r.Caches()).To(Not(ContainElements(cacheNames)))
			default:
				Fail("Unexpected repsonse type")
			}
		})
	})

	Describe(`Validate cache name`, func() {
		It(`CreateCache and DeleteCache errors on bad cache names`, func() {
			badCacheNames := []string{``, `   `}
			for _, badCacheName := range badCacheNames {
				createResp, err := client.CreateCache(ctx, &CreateCacheRequest{CacheName: badCacheName})
				Expect(createResp).To(BeNil())
				var momentoErr MomentoError
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
				}

				deleteResp, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: badCacheName})
				Expect(deleteResp).To(BeNil())
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
				}
			}
		})
	})

	Describe(`DeleteCache`, func() {
		It(`succeeds even if the cache does not exist`, func() {
			Expect(
				client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: uuid.NewString()}),
			).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
		})
	})
})
