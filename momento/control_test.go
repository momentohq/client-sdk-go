package momento_test

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	helpers "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
)

var _ = Describe("control-ops", func() {
	Describe("cache-client happy-path", Label(CACHE_SERVICE_LABEL), func() {
		It("creates, lists, and deletes caches", func() {
			cacheNames := []string{helpers.NewRandomString(), helpers.NewRandomString()}
			defer func() {
				for _, cacheName := range cacheNames {
					_, err := sharedContext.Client.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: cacheName})
					if err != nil {
						panic(err)
					}
				}
			}()

			for _, cacheName := range cacheNames {
				Expect(
					sharedContext.Client.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&CreateCacheSuccess{}))

				Expect(
					sharedContext.Client.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&CreateCacheAlreadyExists{}))
			}

			resp, err := sharedContext.Client.ListCaches(sharedContext.Ctx, &ListCachesRequest{})
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
					sharedContext.Client.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
			}
			resp, err = sharedContext.Client.ListCaches(sharedContext.Ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())
			Expect(resp).To(BeAssignableToTypeOf(&ListCachesSuccess{}))
			switch r := resp.(type) {
			case *ListCachesSuccess:
				Expect(r.Caches()).To(Not(ContainElements(cacheNames)))
			default:
				Fail("Unexpected response type")
			}
		})

		It("creates and deletes using a default cache", func() {
			// Create a separate client with a default cache name to be used only in this test
			// to avoid affecting the shared context when all tests run
			defaultCacheName := fmt.Sprintf("golang-default-%s", helpers.NewRandomString())
			clientWithDefaultCacheName, err := NewCacheClientWithDefaultCache(
				sharedContext.Configuration, sharedContext.CredentialProvider, sharedContext.DefaultTtl, defaultCacheName,
			)
			if err != nil {
				panic(err)
			}
			DeferCleanup(func() { clientWithDefaultCacheName.Close() })

			Expect(
				clientWithDefaultCacheName.CreateCache(sharedContext.Ctx, &CreateCacheRequest{}),
			).Error().NotTo(HaveOccurred())
			Expect(
				clientWithDefaultCacheName.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{}),
			).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
		})

	})

	Describe("cache-client default-cache-name", Label(CACHE_SERVICE_LABEL), func() {
		It("overrides default cache name", func() {
			// Create a separate client with a default cache name to be used only in this test
			// to avoid affecting the shared context when all tests run
			defaultCacheName := fmt.Sprintf("golang-default-%s", uuid.NewString())
			clientWithDefaultCacheName, err := NewCacheClientWithDefaultCache(
				sharedContext.Configuration, sharedContext.CredentialProvider, sharedContext.DefaultTtl, defaultCacheName,
			)
			if err != nil {
				panic(err)
			}
			DeferCleanup(func() { clientWithDefaultCacheName.Close() })

			newCacheName := uuid.NewString()
			Expect(
				clientWithDefaultCacheName.CreateCache(
					sharedContext.Ctx, &CreateCacheRequest{CacheName: newCacheName},
				),
			).Error().NotTo(HaveOccurred())
			Expect(
				clientWithDefaultCacheName.Get(
					sharedContext.Ctx, &GetRequest{Key: helpers.NewRandomMomentoString()},
				),
			).Error().To(HaveMomentoErrorCode(CacheNotFoundError))
			Expect(
				clientWithDefaultCacheName.Get(
					sharedContext.Ctx, &GetRequest{
						CacheName: newCacheName,
						Key:       helpers.NewRandomMomentoString(),
					},
				),
			).To(BeAssignableToTypeOf(&GetMiss{}))
			Expect(
				clientWithDefaultCacheName.DeleteCache(
					sharedContext.Ctx, &DeleteCacheRequest{CacheName: newCacheName},
				),
			).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
		})
	})

	Describe("cache-client validate-cache-name", Label(CACHE_SERVICE_LABEL), func() {
		It("CreateCache and DeleteCache errors on bad cache names", func() {
			badCacheNames := []string{"", "   "}
			for _, badCacheName := range badCacheNames {
				createResp, err := sharedContext.Client.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: badCacheName})
				Expect(createResp).To(BeNil())
				var momentoErr MomentoError
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
				}

				deleteResp, err := sharedContext.Client.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: badCacheName})
				Expect(deleteResp).To(BeNil())
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
				}
			}
		})
	})

	Describe("cache-client delete-cache", Label(CACHE_SERVICE_LABEL), func() {
		It("succeeds even if the cache does not exist", func() {
			Expect(
				sharedContext.Client.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: uuid.NewString()}),
			).To(BeAssignableToTypeOf(&DeleteCacheSuccess{}))
		})
	})
})
