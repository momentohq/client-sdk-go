package impl_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"

	"github.com/momentohq/client-sdk-go/config/compression"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/middleware/impl"
	. "github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("gzip-compression-middleware", Label("cache-service"), func() {
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
	Describe("when IncludeTypes is not specified", func() {
		It("should successfully set and get a value", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
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
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
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
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
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
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
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
	Describe("when IncludeTypes is specified", func() {
		It("should successfully set and get a value without compression when not included", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
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

		It("should not decompress when response was not compressed", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
					IncludeTypes: []interface{}{
						GetRequest{}, // try to decompress without any compression
					},
				}),
			}))

			value := "some-value"
			_, err := cacheClient.Set(testCtx, &SetRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(value),
			})
			Expect(err).To(BeNil())

			// should still be able to fetch, attempted decompression should be a no-op,
			// trace logs should indicate that decompression was not performed
			resp, err := cacheClient.Get(testCtx, &GetRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			Expect(resp.(*responses.GetHit).ValueString()).To(Equal(value))
		})
	})

	Describe("when using json data", func() {
		It("should successfully set and get and compress a json object", func() {
			createCacheClient(config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				impl.NewGzipCompressionMiddleware(impl.GzipCompressionMiddlewareProps{
					CompressionStrategyProps: compression.CompressionStrategyProps{
						CompressionLevel: compression.CompressionLevelDefault,
						Logger:           momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE).GetLogger("gzip-test"),
					},
				}),
			}))

			// User represents a sample JSON object
			type User struct {
				ID    int      `json:"id"`
				Name  string   `json:"name"`
				Email string   `json:"email"`
				Tags  []string `json:"tags"`
			}

			sampleUser := User{
				ID:    1,
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Tags:  []string{"tag1", "tag2"},
			}

			sampleUserJSON, err := json.Marshal(sampleUser)
			Expect(err).To(BeNil())

			_, err = cacheClient.Set(testCtx, &SetRequest{
				CacheName: cacheName,
				Key:       String("key"),
				Value:     String(sampleUserJSON),
			})
			Expect(err).To(BeNil())

			resp, err := cacheClient.Get(testCtx, &GetRequest{
				CacheName: cacheName,
				Key:       String("key"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			var retrievedUser User
			jsonErr := json.Unmarshal(resp.(*responses.GetHit).ValueByte(), &retrievedUser)
			Expect(jsonErr).To(BeNil())
			Expect(retrievedUser).To(Equal(sampleUser))
		})
	})
})
