package momento_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
)

func HaveMomentoErrorCode(code string) types.GomegaMatcher {
	return WithTransform(
		func(err error) (string, error) {
			switch mErr := err.(type) {
			case MomentoError:
				return mErr.Code(), nil
			default:
				return "", fmt.Errorf("Expected MomentoError, but got %T", err)
			}
		}, Equal(code),
	)
}

var _ = Describe("Scalar methods", func() {
	var clientProps SimpleCacheClientProps
	var credentialProvider auth.CredentialProvider
	var configuration config.Configuration
	var client SimpleCacheClient
	var defaultTTL time.Duration
	var testCacheName string
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		credentialProvider, _ = auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
		configuration = config.LatestLaptopConfig()
		defaultTTL = 1 * time.Second

		clientProps = SimpleCacheClientProps{
			CredentialProvider: credentialProvider,
			Configuration:      configuration,
			DefaultTTL:         defaultTTL,
		}

		var err error
		client, err = NewSimpleCacheClient(&clientProps)
		if err != nil {
			panic(err)
		}
		DeferCleanup(func() { client.Close() })

		testCacheName = uuid.NewString()
		Expect(
			client.CreateCache(ctx, &CreateCacheRequest{CacheName: testCacheName}),
		).To(BeAssignableToTypeOf(&CreateCacheSuccess{}))
		DeferCleanup(func() {
			client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: testCacheName})
		})
	})

	DescribeTable(`Gets, Sets, and Deletes`,
		func(key Key, value Value, expectedString string, expectedBytes []byte) {
			Expect(
				client.Set(ctx, &SetRequest{
					CacheName: testCacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			getResp, err := client.Get(ctx, &GetRequest{
				CacheName: testCacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail("Unexpected type from Get")
			}

			Expect(
				client.Delete(ctx, &DeleteRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&DeleteSuccess{}))

			Expect(
				client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("when the key and value are strings", String("key"), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", Bytes([]byte{1, 2, 3}), Bytes([]byte("string")), "string", []byte("string")),
		Entry("when the value is empty", String("key"), String(""), "", []byte("")),
		Entry("when the value is blank", String("key"), String("  "), "  ", []byte("  ")),
	)

	It(`errors when the cache is missing`, func() {
		cacheName := uuid.NewString()
		key := String("key")
		value := String("value")

		getResp, err := client.Get(ctx, &GetRequest{
			CacheName: cacheName,
			Key:       key,
		})
		Expect(getResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		setResp, err := client.Set(ctx, &SetRequest{
			CacheName: cacheName,
			Key:       key,
			Value:     value,
		})
		Expect(setResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		deleteResp, err := client.Delete(ctx, &DeleteRequest{
			CacheName: cacheName,
			Key:       key,
		})
		Expect(deleteResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))
	})

	DescribeTable(`invalid cache names and keys`,
		func(cacheName string, key Key, value Key) {
			getResp, err := client.Get(ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(getResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			setResp, err := client.Set(ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(setResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			deleteResp, err := client.Delete(ctx, &DeleteRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(deleteResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry(`With an empty cache name`, "", String("key"), String("value")),
		Entry(`With an bank cache name`, "   ", String("key"), String("value")),
		Entry(`With an empty key`, cacheName(), String(""), String("value")),
	)

	Describe(`Set`, func() {
		It(`Uses the default TTL`, func() {
			key := String("key")
			value := String("value")

			Expect(
				client.Set(ctx, &SetRequest{
					CacheName: testCacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(defaultTTL / 2)

			Expect(
				client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(defaultTTL)

			Expect(
				client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It(`Overrides the default TTL`, func() {
			key := String("key")
			value := String("value")

			Expect(
				client.Set(ctx, &SetRequest{
					CacheName: testCacheName,
					Key:       key,
					Value:     value,
					TTL:       defaultTTL * 2,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(defaultTTL / 2)

			Expect(
				client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(defaultTTL)

			Expect(
				client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(defaultTTL)

			Expect(
				client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})
	})
})
