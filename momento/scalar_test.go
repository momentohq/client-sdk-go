package momento_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
)

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
		).To(Succeed())
		DeferCleanup(func() {
			client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: testCacheName})
		})
	})

	DescribeTable(`Gets, Sets, and Deletes`,
		func(key Value, value Value, expectedString string, expectedBytes []byte) {
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
