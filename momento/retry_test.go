package momento_test

import (
	"context"
	"os"
	"strconv"

	"github.com/momentohq/client-sdk-go/config/retry"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	helpers "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/momento_rpc_names"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/responses"

	"time"

	"github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
)

const (
	CLIENT_TIMEOUT_MILLIS                 = 3 * time.Second
	RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS = 1000
	RETRY_DELAY_INTERVAL_MILLIS           = 100
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
			func(grpcStatus codes.Code, requestMethod momento_rpc_names.MomentoRPCMethod, expected bool) {
				strategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
					MaxAttempts:   3,
				})
				retryResult := strategy.DetermineWhenToRetry(
					retry.StrategyProps{GrpcStatusCode: grpcStatus, GrpcMethod: string(requestMethod), AttemptNumber: 1},
				)

				if expected == false {
					Expect(retryResult).To(BeNil())
				} else {
					Expect(retryResult).To(Not(BeNil()))
					Expect(*retryResult).To(Equal(0))
				}
			},
			// Entry("name", codes.Internal, "/cache_client.Scs/Get", true),
			Entry("name", codes.Internal, momento_rpc_names.Get, true),
			Entry("name", codes.Internal, momento_rpc_names.Set, true),
			Entry("name", codes.Internal, momento_rpc_names.DictionaryIncrement, false),
			Entry("name", codes.Unknown, momento_rpc_names.Get, false),
			Entry("name", codes.Unknown, momento_rpc_names.Set, false),
			Entry("name", codes.Unknown, momento_rpc_names.DictionaryIncrement, false),
			Entry("name", codes.Unavailable, momento_rpc_names.Get, true),
			Entry("name", codes.Unavailable, momento_rpc_names.Set, true),
			Entry("name", codes.Unavailable, momento_rpc_names.DictionaryIncrement, false),
			Entry("name", codes.Canceled, momento_rpc_names.Get, false),
			Entry("name", codes.Canceled, momento_rpc_names.Set, false),
			Entry("name", codes.Canceled, momento_rpc_names.DictionaryIncrement, false),
			Entry("name", codes.DeadlineExceeded, momento_rpc_names.Get, false),
			Entry("name", codes.DeadlineExceeded, momento_rpc_names.Set, false),
			Entry("name", codes.DeadlineExceeded, momento_rpc_names.DictionaryIncrement, false),
		)

		Describe("cache-client retry neverRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("shouldn't retry", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				strategy := retry.NewNeverRetryStrategy()
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Set)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(
					helpers.MomentoLocalMiddlewareProps{
						MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
							ReturnError:  &status,
							ErrorRpcList: &[]string{rpcName},
							ErrorCount:   nil,
							DelayRpcList: nil,
							DelayMillis:  nil,
							DelayCount:   nil,
						},
					},
				)
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})
		})

		Describe("cache-client retry exponentialBackoffRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should receive a timeout error after multiple retries", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Set)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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
				retries, err := metricsCollector.GetTotalRetryCount("cache", rpcName)
				Expect(err).To(BeNil())
				Expect(retries > 1).To(BeTrue())
			})

			It("should not try to retry if the status code is not retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.UnknownServiceError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Set)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})

			It("should not try to retry if the rpc is not retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.DictionaryIncrement)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})
		})

		Describe("cache-client retry fixedCountRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should retry 3 times if the status code is retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
					MaxAttempts:   3,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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

				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(3))
				Expect(metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)).To(BeNumerically("<=", 10))
			})

			It("should not retry if the status code is not retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.UnknownServiceError)
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Set)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})

			It("should not retry if the api is not retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.DictionaryIncrement)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   nil,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				})
				cacheClient := setupCacheClientTest(clientConfig)

				dictIncrementResponse, err := cacheClient.DictionaryIncrement(context.Background(), &momento.DictionaryIncrementRequest{
					CacheName:      "cache",
					DictionaryName: "myDict",
					Field:          momento.String("field"),
					Amount:         int64(1),
				})
				Expect(dictIncrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})

			It("should return a value on success after a retry", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				errCount := 1
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   &errCount,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
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
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(1))
			})
		})

		Describe("cache-client retry fixedTimeoutRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should not retry if the status code is not retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.UnknownServiceError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Set)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &momento.SetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
					Value:     momento.String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})

			It("should not retry if the rpc is not retryable", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.DictionaryIncrement)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
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
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(0))
			})

			It("should use default timeout values when not specified", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				cacheClient := setupCacheClientTest(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.TimeoutError))
				Expect(getResponse).To(BeNil())

				// Should immediately receive errors and retry every DefaultRetryDelayIntervalMillis
				// until the client timeout is reached.
				maxAttempts := CLIENT_TIMEOUT_MILLIS / retry.DefaultRetryDelayIntervalMillis
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				Expect(metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)).To(BeNumerically("<=", retry.DefaultRetryDelayIntervalMillis*1.1))
				Expect(metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)).To(BeNumerically(">=", retry.DefaultRetryDelayIntervalMillis*0.9))
			})

			It("should retry until client timeout when responses have no delays during full outage", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				cacheClient := setupCacheClientTest(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.TimeoutError))

				// Should immediately receive errors and retry every DefaultRetryDelayIntervalMillis
				// until the client timeout is reached.
				maxAttempts := CLIENT_TIMEOUT_MILLIS / RETRY_DELAY_INTERVAL_MILLIS
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 1.1
				minDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)
				Expect(err).To(BeNil())
				Expect(average).To(BeNumerically("<=", int64(maxDelay)))
				Expect(average).To(BeNumerically(">=", int64(minDelay)))
			})

			It("should retry until client timeout when responses have short delays during full outage", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				shortDelay := RETRY_DELAY_INTERVAL_MILLIS + 100
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						DelayRpcList: &[]string{rpcName},
						DelayMillis:  &shortDelay,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				cacheClient := setupCacheClientTest(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.TimeoutError))

				// Should receive errors after shortDelay ms and retry every RETRY_DELAY_INTERVAL_MILLIS
				// until the client timeout is reached.
				delayBetweenAttempts := RETRY_DELAY_INTERVAL_MILLIS + shortDelay
				maxAttempts := int(CLIENT_TIMEOUT_MILLIS.Milliseconds()) / delayBetweenAttempts
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(delayBetweenAttempts) * 1.1
				minDelay := float64(delayBetweenAttempts) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)
				Expect(err).To(BeNil())
				Expect(float64(average)).To(BeNumerically("<=", maxDelay))
				Expect(float64(average)).To(BeNumerically(">=", minDelay))
			})

			It("should retry until client timeout when responses have long delays during full outage", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				longDelay := RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS + 100
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						DelayRpcList: &[]string{rpcName},
						DelayMillis:  &longDelay,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				cacheClient := setupCacheClientTest(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(momento.TimeoutError))

				// Should receive errors after longDelay ms and retry every RETRY_DELAY_INTERVAL_MILLIS
				// until the client timeout is reached.
				delayBetweenAttempts := RETRY_DELAY_INTERVAL_MILLIS + longDelay
				maxAttempts := int(CLIENT_TIMEOUT_MILLIS.Milliseconds()) / delayBetweenAttempts
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(delayBetweenAttempts) * 1.1
				minDelay := float64(delayBetweenAttempts) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)
				Expect(err).To(BeNil())
				Expect(float64(average)).To(BeNumerically("<=", maxDelay))
				Expect(float64(average)).To(BeNumerically(">=", minDelay))
			})

			It("should retry until partial outage is resolved", func() {
				status := helpers.ConvertErrorCodeToMomentoLocalErrorCode(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				errCount := 3
				rpcName := helpers.ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.Get)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{rpcName},
						ErrorCount:   &errCount,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				cacheClient := setupCacheClientTest(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &momento.GetRequest{
					CacheName: "cache",
					Key:       momento.String("key"),
				})
				Expect(getResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())

				// Should retry until the server stops returning errors
				Expect(metricsCollector.GetTotalRetryCount("cache", rpcName)).To(Equal(errCount))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 1.1
				minDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries("cache", rpcName)
				Expect(err).To(BeNil())
				Expect(average).To(BeNumerically("<=", int64(maxDelay)))
				Expect(average).To(BeNumerically(">=", int64(minDelay)))
			})
		})
	},
)
