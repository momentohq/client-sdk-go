package impl_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/middleware/impl"
	. "github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const MIDDLEWARE_IMPL_LABEL = "middleware-impl"

var (
	testCtx     context.Context
	cacheName   string
	cacheClient CacheClient
)

func cleanup() {
	deleteResponse, err := cacheClient.DeleteCache(context.Background(), &DeleteCacheRequest{
		CacheName: cacheName,
	})
	Expect(err).To(BeNil())
	Expect(deleteResponse).To(Not(BeNil()))
}

// Each test may use a different compression middleware in the config
func createCacheClient(config config.Configuration) {
	credentialProvider, err := auth.FromEnvironmentVariable("MOMENTO_API_KEY")
	Expect(err).To(BeNil())

	cacheClient, err = NewCacheClient(
		config,
		credentialProvider,
		time.Second*60,
	)
	Expect(err).To(BeNil())

	_, err = cacheClient.CreateCache(testCtx, &CreateCacheRequest{CacheName: cacheName})
	Expect(err).To(BeNil())
}

func getCompressableString() string {
	longString := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum"
	return fmt.Sprintf("%s %s", longString, uuid.NewString())
}

var _ = Describe("middleware-impl zstd-compression", Label(MIDDLEWARE_IMPL_LABEL), func() {
	BeforeEach(func() {
		testCtx = context.Background()
		cacheName = fmt.Sprintf("golang-%s", uuid.NewString())
		cacheClient = nil
	})

	AfterEach(func() {
		cleanup()
	})

	// Some happy-path tests to verify set/get methods still work even when compression middleware is enabled
	// and no IncludeTypes are specified to narrow down the types of requests that should be compressed.
	Describe("should compress and decompress when IncludeTypes is not specified", func() {
		It("should successfully set and get a value", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewZstdCompressionMiddleware(impl.ZstdCompressionMiddlewareProps{
					Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("zstd-test"),
					CompressionLevel: zstd.SpeedBetterCompression,
				}),
			}))

			value := getCompressableString()
			_, err := cacheClient.Set(testCtx, &SetRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(value),
			})
			Expect(err).To(BeNil())

			resp, err := cacheClient.Get(testCtx, &GetRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			Expect(resp.(*responses.GetHit).ValueString()).To(Equal(value))
		})

		It("should successfully setIf and get a value", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewZstdCompressionMiddleware(impl.ZstdCompressionMiddlewareProps{
					Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("zstd-test"),
					CompressionLevel: zstd.SpeedBetterCompression,
				}),
			}))

			setIfAbsentValue := getCompressableString()
			_, err := cacheClient.SetIfAbsent(testCtx, &SetIfAbsentRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(setIfAbsentValue),
			})
			Expect(err).To(BeNil())

			resp, err := cacheClient.Get(testCtx, &GetRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			Expect(resp.(*responses.GetHit).ValueString()).To(Equal(setIfAbsentValue))

			setIfPresentValue := getCompressableString()
			_, err = cacheClient.SetIfPresentAndNotEqual(testCtx, &SetIfPresentAndNotEqualRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(setIfPresentValue),
				NotEqual:  String("some other string"),
			})
			Expect(err).To(BeNil())

			resp, err = cacheClient.Get(testCtx, &GetRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			Expect(resp.(*responses.GetHit).ValueString()).To(Equal(setIfPresentValue))
		})

		It("should successfully setWithHash and getWithHash", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewZstdCompressionMiddleware(impl.ZstdCompressionMiddlewareProps{
					Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("zstd-test"),
					CompressionLevel: zstd.SpeedBetterCompression,
				}),
			}))

			value := getCompressableString()
			setResp, err := cacheClient.SetWithHash(testCtx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(value),
			})
			Expect(err).To(BeNil())
			Expect(setResp).To(BeAssignableToTypeOf(&responses.SetWithHashStored{}))
			Expect(setResp.(*responses.SetWithHashStored).HashByte()).To(Not(BeEmpty()))
			hash := setResp.(*responses.SetWithHashStored).HashByte()

			resp, err := cacheClient.GetWithHash(testCtx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetWithHashHit{}))
			Expect(resp.(*responses.GetWithHashHit).ValueString()).To(Equal(value))
			Expect(resp.(*responses.GetWithHashHit).HashByte()).To(Equal(hash))
		})

		It("should successfully setIfHash and getWithHash", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewZstdCompressionMiddleware(impl.ZstdCompressionMiddlewareProps{
					Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("zstd-test"),
					CompressionLevel: zstd.SpeedBetterCompression,
				}),
			}))

			setIfAbsentValue := getCompressableString()
			setAbsentResp, err := cacheClient.SetIfAbsentOrHashEqual(testCtx, &SetIfAbsentOrHashEqualRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(setIfAbsentValue),
				HashEqual: Bytes("some-hash-value"),
			})
			Expect(err).To(BeNil())
			Expect(setAbsentResp).To(BeAssignableToTypeOf(&responses.SetIfAbsentOrHashEqualStored{}))
			Expect(setAbsentResp.(*responses.SetIfAbsentOrHashEqualStored).HashByte()).To(Not(BeEmpty()))
			hash := setAbsentResp.(*responses.SetIfAbsentOrHashEqualStored).HashByte()

			resp, err := cacheClient.GetWithHash(testCtx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetWithHashHit{}))
			Expect(resp.(*responses.GetWithHashHit).ValueString()).To(Equal(setIfAbsentValue))
			Expect(resp.(*responses.GetWithHashHit).HashByte()).To(Equal(hash))

			setIfPresentValue := getCompressableString()
			setPresentResp, err := cacheClient.SetIfPresentAndHashEqual(testCtx, &SetIfPresentAndHashEqualRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(setIfPresentValue),
				HashEqual: Bytes(hash),
			})
			Expect(err).To(BeNil())
			Expect(setPresentResp).To(BeAssignableToTypeOf(&responses.SetIfPresentAndHashEqualStored{}))
			Expect(setPresentResp.(*responses.SetIfPresentAndHashEqualStored).HashByte()).To(Not(BeEmpty()))
			hash = setPresentResp.(*responses.SetIfPresentAndHashEqualStored).HashByte()

			resp, err = cacheClient.GetWithHash(testCtx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetWithHashHit{}))
			Expect(resp.(*responses.GetWithHashHit).ValueString()).To(Equal(setIfPresentValue))
			Expect(resp.(*responses.GetWithHashHit).HashByte()).To(Equal(hash))
		})

	})

	// Some tests to verify set/get methods still work even when compression middleware is enabled and
	// IncludeTypes is specified to narrow down the types of requests that should be compressed.
	Describe("should compress and decompress when IncludeTypes is specified", func() {
		It("should successfully set and get a value without compression when not included", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewZstdCompressionMiddleware(impl.ZstdCompressionMiddlewareProps{
					Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG).GetLogger("zstd-test"),
					CompressionLevel: zstd.SpeedBetterCompression,
					IncludeTypes: []interface{}{
						SetWithHashRequest{},
						GetWithHashRequest{},
					},
				}),
			}))

			// Should not see Get or Set mentioned in compression logs

			value := getCompressableString()
			nonHashKey := uuid.NewString()
			_, err := cacheClient.Set(testCtx, &SetRequest{
				CacheName: cacheName,
				Key:       String(nonHashKey),
				Value:     String(value),
			})
			Expect(err).To(BeNil())

			resp, err := cacheClient.Get(testCtx, &GetRequest{
				CacheName: cacheName,
				Key:       String(nonHashKey),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			Expect(resp.(*responses.GetHit).ValueString()).To(Equal(value))

			// Should see SetWithHash and GetWithHash mentioned in compression logs

			hashKey := getCompressableString()
			setResp, err := cacheClient.SetWithHash(testCtx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       String(hashKey),
				Value:     String(value),
			})
			Expect(err).To(BeNil())
			Expect(setResp).To(BeAssignableToTypeOf(&responses.SetWithHashStored{}))
			Expect(setResp.(*responses.SetWithHashStored).HashByte()).To(Not(BeEmpty()))
			hash := setResp.(*responses.SetWithHashStored).HashByte()

			getResp, err := cacheClient.GetWithHash(testCtx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       String(hashKey),
			})
			Expect(err).To(BeNil())
			Expect(getResp).To(BeAssignableToTypeOf(&responses.GetWithHashHit{}))
			Expect(getResp.(*responses.GetWithHashHit).ValueString()).To(Equal(value))
			Expect(getResp.(*responses.GetWithHashHit).HashByte()).To(Equal(hash))
		})
	})
})
