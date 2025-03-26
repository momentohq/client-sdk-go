package momento_test

import (
	"context"
	"fmt"
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

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
)

func setupCacheClientTest(config config.Configuration) CacheClient {
	momentoLocalPort := os.Getenv("MOMENTO_PORT")
	if momentoLocalPort == "" {
		momentoLocalPort = "8080"
	}
	thePort, err := strconv.ParseUint(momentoLocalPort, 10, 32)
	Expect(err).To(BeNil())
	credentialProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{Port: uint(thePort)})
	Expect(err).To(BeNil())
	cacheClient, err := NewCacheClient(config, credentialProvider, 30*time.Second)
	Expect(err).To(BeNil())
	createResponse, err := cacheClient.CreateCache(context.Background(), &CreateCacheRequest{
		CacheName: "cache",
	})
	Expect(err).To(BeNil())
	Expect(createResponse).To(Not(BeNil()))
	return cacheClient
}

func setupTopicClientTest(config config.TopicsConfiguration) TopicClient {
	momentoLocalPort := os.Getenv("MOMENTO_PORT")
	if momentoLocalPort == "" {
		momentoLocalPort = "8080"
	}
	thePort, err := strconv.ParseUint(momentoLocalPort, 10, 32)
	Expect(err).To(BeNil())
	credentialProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{Port: uint(thePort)})
	Expect(err).To(BeNil())
	topicClient, err := NewTopicClient(config, credentialProvider)
	Expect(err).To(BeNil())
	return topicClient
}

var _ = Describe("retry eligibility-strategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
	Describe("Unary Interceptor Retries", func() {
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
				retryMiddleware := helpers.NewMomentoLocalMiddleware(
					helpers.MomentoLocalMiddlewareProps{
						MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
							ReturnError:  &status,
							ErrorRpcList: &[]string{"set"},
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
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
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
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
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
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))
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
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
						ErrorCount:   &errorCount,
						DelayRpcList: nil,
						DelayMillis:  nil,
						DelayCount:   nil,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(strategy).WithClientTimeout(10 * time.Second)
				cacheClient := setupCacheClientTest(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
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
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
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
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Set")).To(Equal(0))
			})

			It("should not try to retry if the rpc is not retryable", func() {
				status := "unavailable"
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"dictionary-increment"},
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
				incrResponse, err := cacheClient.DictionaryIncrement(context.Background(), &DictionaryIncrementRequest{
					CacheName:      "cache",
					DictionaryName: "dictionary",
					Field:          String("field"),
					Amount:         1,
				})
				Expect(incrResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
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
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
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
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: "cache",
					Key:       String("key"),
				})
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(getResponse).To(BeNil())

				Expect(metricsCollector.GetTotalRetryCount("cache", "Get")).To(Equal(3))
				Expect(metricsCollector.GetAverageTimeBetweenRetries("cache", "Get")).To(Equal(int64(0)))
			})

			It("should not retry if the status code is not retryable", func() {
				status := "unknown"
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
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

				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Set")).To(Equal(0))
			})

			It("should not retry if the api is not retryable", func() {
				status := "unavailable"
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"increment", "dictionary-increment"},
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

				incrementResponse, err := cacheClient.Increment(context.Background(), &IncrementRequest{
					CacheName: "cache",
					Field:     String("key"),
				})
				Expect(incrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))

				dictCreateResponse, err := cacheClient.DictionarySetField(context.Background(), &DictionarySetFieldRequest{
					CacheName:      "cache",
					DictionaryName: "myDict",
					Field:          String("key"),
					Value:          String("value"),
					Ttl:            &utils.CollectionTtl{Ttl: 600 * time.Second},
				})
				Expect(dictCreateResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())

				dictIncrementResponse, err := cacheClient.DictionaryIncrement(context.Background(), &DictionaryIncrementRequest{
					CacheName:      "cache",
					DictionaryName: "myDict",
					Field:          String("field"),
					Amount:         int64(1),
				})
				Expect(dictIncrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Increment")).To(Equal(0))
			})

			It("should return a value on success after a retry", func() {
				status := "unavailable"
				errCount := 1
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
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
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: "cache",
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: "cache",
					Key:       String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResponse).To(Not(BeNil()))
				Expect(getResponse.(*responses.GetHit).ValueString()).To(Equal("value"))
				Expect(metricsCollector.GetTotalRetryCount("cache", "Get")).To(Equal(1))
			})
		})
	})

	Describe("Stream Interceptor Retries", func() {
		It("publishes and subscribes", func() {
			ctx := context.Background()
			cacheName := "cache"
			topicName := "topic"
			status := "unavailable"
			streamErrorMessageLimit := 10
			strategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
				LoggerFactory: momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
				MaxAttempts:   10,
			})
			retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
				MomentoLocalMiddlewareRequestHandlerProps: helpers.MomentoLocalMiddlewareRequestHandlerProps{
					StreamError: &status,
					StreamErrorRpcList: &[]string{"topic-subscribe"},
					StreamErrorMessageLimit: &streamErrorMessageLimit,
				},
			})
			//metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
			//clientConfig := config.TopicsDefaultWithLogger(
			//	momento_default_logger.NewDefaultMomentoLoggerFactory(
			//		momento_default_logger.DEBUG)).AddMiddleware(retryMiddleware).WithRetryStrategy(strategy)
			clientConfig := config.TopicsDefaultWithLogger(
				momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE),
			).AddMiddleware(retryMiddleware).WithRetryStrategy(strategy)

			topicClient := setupTopicClientTest(clientConfig)
			_ = setupCacheClientTest(config.LaptopLatest())

			publishedValues := []TopicValue{
				String("aaa00"),
				String("aaa01"),
				String("aaa02"),
				String("aaa03"),
				String("aaa04"),
				String("aaa05"),
				String("aaa06"),
				String("aaa07"),
				String("aaa08"),
				String("aaa09"),
				String("aaa10"),
			}

			sub, err := topicClient.Subscribe(ctx, &TopicSubscribeRequest{
				CacheName: cacheName,
				TopicName: topicName,
			})
			fmt.Printf("client is back from subscribe with %v\n", sub)
			if err != nil {
				panic(err)
			}

			cancelContext, cancelFunction := context.WithCancel(ctx)
			var numReceivedValues int
			ready := make(chan int, 1)
			go func() {
				ready <- 1
				for {
					select {
					case <-cancelContext.Done():
						fmt.Println("subscribe goroutine done")
						return
					default:
						fmt.Println("*** waiting for item")
						item, err := sub.Item(ctx)
						if err != nil {
							panic(err)
						}
						fmt.Printf("received item: %v\n", item)
						numReceivedValues++
						fmt.Printf("received %d values\n", numReceivedValues)
					}
				}
			}()
			<-ready

			time.Sleep(time.Millisecond * 1000)
			for _, value := range publishedValues {
				fmt.Printf("Publishing %s\n", value)
				_, err := topicClient.Publish(ctx, &TopicPublishRequest{
					CacheName: cacheName,
					TopicName: topicName,
					Value:     value,
				})
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("sleeping before cancelling")
			time.Sleep(time.Millisecond * 5000)
			fmt.Println("CANCELLING")
			cancelFunction()

			//fmt.Printf("metrics: %#v", metricsCollector.GetAllMetrics())

			Expect(numReceivedValues).To(Equal(len(publishedValues)))

		})
	})
})
