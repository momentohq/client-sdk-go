package momento_test

import (
	"context"
	"os"
	"strconv"

	"github.com/momentohq/client-sdk-go/config/retry"

	helpers "github.com/momentohq/client-sdk-go/momento/test_helpers"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"

	"time"

	"github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
)

func setupCacheClientTest(config config.Configuration) momento.CacheClient {
	momentoLocalPort := os.Getenv("MOMENTO_PORT")
	if momentoLocalPort == "" {
		momentoLocalPort = "8080"
	}
	thePort, err := strconv.ParseUint(momentoLocalPort, 10, 32)
	Expect(err).To(BeNil())
	credentialProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{Port: uint(thePort)})
	Expect(err).To(BeNil())
	cacheClient, err := momento.NewCacheClient(config, credentialProvider, 30*time.Second)
	Expect(err).To(BeNil())
	createResponse, err := cacheClient.CreateCache(context.Background(), &momento.CreateCacheRequest{
		CacheName: "cache",
	})
	Expect(err).To(BeNil())
	Expect(createResponse).To(Not(BeNil()))
	return cacheClient
}

var _ = Describe(
	"cache-client retry eligibility-strategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
		DescribeTable(
			"DefaultEligibilityStrategy -- determine retry eligibility given grpc status code and request method",
			func(grpcStatus codes.Code, requestMethod string, expected bool) {
				strategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
					MaxAttempts:   3,
				})
				retryResult := strategy.DetermineWhenToRetry(
					retry.StrategyProps{GrpcStatusCode: grpcStatus, GrpcMethod: requestMethod, AttemptNumber: 1},
				)

				if expected == false {
					Expect(retryResult).To(BeNil())
				} else {
					Expect(retryResult).To(Not(BeNil()))
					Expect(*retryResult).To(Equal(0))
				}
			},
			Entry("name", codes.Internal, "/cache_client.Scs/Get", true),
			Entry("name", codes.Internal, "/cache_client.Scs/Set", true),
			Entry("name", codes.Internal, "/cache_client.Scs/DictionaryIncrement", false),
			Entry("name", codes.Unknown, "/cache_client.Scs/Get", false),
			Entry("name", codes.Unknown, "/cache_client.Scs/Set", false),
			Entry("name", codes.Unknown, "/cache_client.Scs/DictionaryIncrement", false),
			Entry("name", codes.Unavailable, "/cache_client.Scs/Get", true),
			Entry("name", codes.Unavailable, "/cache_client.Scs/Set", true),
			Entry("name", codes.Unavailable, "/cache_client.Scs/DictionaryIncrement", false),
			Entry("name", codes.Canceled, "/cache_client.Scs/Get", false),
			Entry("name", codes.Canceled, "/cache_client.Scs/Set", false),
			Entry("name", codes.Canceled, "/cache_client.Scs/DictionaryIncrement", false),
			Entry("name", codes.DeadlineExceeded, "/cache_client.Scs/Get", false),
			Entry("name", codes.DeadlineExceeded, "/cache_client.Scs/Set", false),
			Entry("name", codes.DeadlineExceeded, "/cache_client.Scs/DictionaryIncrement", false),
		)

		Describe("cache-client retry neverRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("shouldn't retry", func() {
				status := "unavailable"
				strategy := retry.NewNeverRetryStrategy()
				retryMiddleware := helpers.NewRetryMetricsMiddleware(
					helpers.RetryMetricsMiddlewareProps{
						RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
							ReturnError:  &status,
							ErrorRpcList: &[]string{"set"},
							ErrorCount:   nil,
							DelayRpcList: nil,
							DelayMillis:  nil,
							DelayCount:   nil,
						},
					},
				)
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(strategy)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Set")).To(Equal(0))
			})
		})

		Describe("cache-client retry exponentialBackoffRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should receive a timeout error after multiple retries", func() {
				status := "unavailable"
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(strategy).WithClientTimeout(1 * time.Second)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.TimeoutError))
				retries, err := metricsCollector.GetTotalRetryCount("cache", "Set")
				Expect(err).To(BeNil())
				Expect(retries > 1).To(BeTrue())
			})

			It("should succeed after multiple retries", func() {
				status := "unavailable"
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				errorCount := 5
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
						ErrorCount:   &errorCount,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(strategy).WithClientTimeout(10 * time.Second)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(setResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())
				Expect(setResponse).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
				retries, err := metricsCollector.GetTotalRetryCount("cache", "Set")
				Expect(err).To(BeNil())
				Expect(retries > 1).To(BeTrue())
			})

			It("should not try to retry if the status code is not retryable", func() {
				status := "unknown"
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(strategy)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Set")).To(Equal(0))
			})

			It("should not try to retry if the rpc is not retryable", func() {
				status := "unavailable"
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"dictionary-increment"},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(strategy)
				cacheClient := setupCacheClientTest(clientConfig)
				incrResponse, err := cacheClient.DictionaryIncrement(context.Background(), &momento.DictionaryIncrementRequest{
					CacheName:      "cache",
					DictionaryName: "dictionary",
					Field:          momento.String("field"),
					Amount:         1,
				})
				Expect(incrResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "DictionaryIncrement")).To(Equal(0))
			})
		})

		Describe("cache-client retry fixedCountRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should retry 3 times if the status code is retryable", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
					MaxAttempts:   3,
				})
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))
				Expect(getResponse).To(BeNil())

				Expect(metricsCollector.GetTotalRetryCount("cache", "Get")).To(Equal(3))
				Expect(metricsCollector.GetAverageTimeBetweenRetries("cache", "Get")).To(Equal(int64(0)))
			})

			It("should not retry if the status code is not retryable", func() {
				status := "unknown"
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				})
				cacheClient := setupCacheClientTest(clientConfig)

				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Set")).To(Equal(0))
			})

			It("should not retry if the api is not retryable", func() {
				status := "unavailable"
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"increment", "dictionary-increment"},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				})
				cacheClient := setupCacheClientTest(clientConfig)

				incrementResponse, err := cacheClient.Increment(context.Background(), &momento.IncrementRequest{
					CacheName: "cache",
					Field:     momento.String("key"),
				})
				Expect(incrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))

				dictCreateResponse, err := cacheClient.DictionarySetField(context.Background(), &momento.DictionarySetFieldRequest{
					CacheName:      "cache",
					DictionaryName: "myDict",
					Field:          momento.String("key"),
					Value:          momento.String("value"),
					Ttl:            &utils.CollectionTtl{Ttl: 600 * time.Second},
				})
				Expect(dictCreateResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())

				dictIncrementResponse, err := cacheClient.DictionaryIncrement(context.Background(), &momento.DictionaryIncrementRequest{
					CacheName:      "cache",
					DictionaryName: "myDict",
					Field:          momento.String("field"),
					Amount:         int64(1),
				})
				Expect(dictIncrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Increment")).To(Equal(0))
			})

			It("should return a value on success after a retry", func() {
				status := "unavailable"
				errCount := 1
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						ErrorCount:   &errCount,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				})
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResponse).To(Not(BeNil()))
				Expect(getResponse.(*responses.GetHit).ValueString()).To(Equal("value"))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Get")).To(Equal(1))
			})

			It("shouldn't gather metrics if the request is not included", func() {
				status := "unavailable"
				errCount := 3
				retryMiddleware := helpers.NewRetryMetricsMiddleware(helpers.RetryMetricsMiddlewareProps{
					RetryMetricsMiddlewareRequestHandlerProps: helpers.RetryMetricsMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						ErrorCount:   &errCount,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.RetryMetricsMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				})
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResponse).To(Not(BeNil()))
				Expect(getResponse.(*responses.GetHit).ValueString()).To(Equal("value"))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Get")).To(Equal(0))
			})

		})
	},
)
