package momento_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/config/retry"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

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

const (
	CLIENT_TIMEOUT_MILLIS                 = 3 * time.Second
	RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS = 1000
	RETRY_DELAY_INTERVAL_MILLIS           = 100
)

var (
	testCtx     context.Context
	cacheName   string
	topicName   string
	cacheClient CacheClient
)

// I'm choosing not to use the usual shared context pattern here. That whole framework is bloated and
// overly prescriptive. I just want to set up a few things and run some tests.

func setCacheClient(config config.Configuration) {
	momentoLocalPort := os.Getenv("MOMENTO_PORT")
	if momentoLocalPort == "" {
		momentoLocalPort = "8080"
	}
	thePort, err := strconv.ParseUint(momentoLocalPort, 10, 32)
	Expect(err).To(BeNil())
	credentialProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{Port: uint(thePort)})
	Expect(err).To(BeNil())
	cacheClient, err = NewCacheClient(config, credentialProvider, 30*time.Second)
	Expect(err).To(BeNil())
}

func setupCacheClient(config config.Configuration) {
	setCacheClient(config)
	createResponse, err := cacheClient.CreateCache(context.Background(), &CreateCacheRequest{
		CacheName: cacheName,
	})
	Expect(err).To(BeNil())
	Expect(createResponse).To(Not(BeNil()))
}

func cleanup() {
	deleteResponse, err := cacheClient.DeleteCache(context.Background(), &DeleteCacheRequest{
		CacheName: cacheName,
	})
	Expect(err).To(BeNil())
	Expect(deleteResponse).To(Not(BeNil()))
}

func setupTopicClient(config config.TopicsConfiguration) TopicClient {
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

type clientConfigProps struct {
	status                  *string
	streamErrorMessageLimit *int
	streamErrorRpcList      *[]string
	delayRpcList            *[]string
	delayMillis             *int
}

func getClientConfig(props *clientConfigProps) (config.TopicsConfiguration, middleware.Middleware) {
	strategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
		LoggerFactory: momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
		MaxAttempts:   10,
	})
	retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
		MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
			StreamError:             props.status,
			StreamErrorRpcList:      props.streamErrorRpcList,
			StreamErrorMessageLimit: props.streamErrorMessageLimit,
			DelayRpcList:            props.delayRpcList,
			DelayMillis:             props.delayMillis,
		},
	})
	return config.TopicsDefaultWithLogger(
		momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.TRACE),
	).AddMiddleware(retryMiddleware.(middleware.TopicMiddleware)).WithRetryStrategy(strategy), retryMiddleware
}

func doPubSub(topicClient TopicClient, publishedValues []TopicValue) error {
	sub, err := topicClient.Subscribe(testCtx, &TopicSubscribeRequest{
		CacheName: cacheName,
		TopicName: topicName,
	})
	if err != nil {
		fmt.Printf("error from topic subscription: %v\n", err)
		return err
	}

	cancelContext, cancelFunction := context.WithCancel(testCtx)
	ready := make(chan int, 1)
	go func() {
		ready <- 1
		for {
			select {
			case <-cancelContext.Done():
				return
			default:
				_, err := sub.Item(cancelContext)
				if err != nil {
					// canceled errors are expected, so we can ignore them
					if errors.Is(err, context.Canceled) {
						return
					}
					var svcErr momentoerrors.MomentoSvcErr
					switch {
					case errors.As(err, &svcErr):
						if svcErr.Code() == momentoerrors.CanceledError {
							return
						}
					default:
						panic(err)
					}
				}
			}
		}
	}()
	<-ready

	time.Sleep(time.Millisecond * 1000)
	for _, value := range publishedValues {
		_, err := topicClient.Publish(testCtx, &TopicPublishRequest{
			CacheName: cacheName,
			TopicName: topicName,
			Value:     value,
		})
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(time.Millisecond * 5000)
	cancelFunction()
	return nil
}

var _ = Describe("retry eligibility-strategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
	Describe("Eligibility Strategy Testing", func() {
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
	})

	Describe("Retry Strategy Testing", func() {

		BeforeEach(func() {
			testCtx = context.Background()
			cacheName = uuid.NewString()
			topicName = uuid.NewString()
			cacheClient = nil
		})

		AfterEach(func() {
			cleanup()
		})

		Describe("cache-client retry neverRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("shouldn't retry", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				strategy := retry.NewNeverRetryStrategy()
				retryMiddleware := helpers.NewMomentoLocalMiddleware(
					helpers.MomentoLocalMiddlewareProps{
						MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Set")).To(Equal(0))
			})
		})

		Describe("cache-client retry exponentialBackoffRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should receive a timeout error after multiple retries", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))
				retries, err := metricsCollector.GetTotalRetryCount(cacheName, "Set")
				Expect(err).To(BeNil())
				Expect(retries > 1).To(BeTrue())
			})

			It("should succeed after multiple retries", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				errorCount := 5
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())
				Expect(setResponse).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
				retries, err := metricsCollector.GetTotalRetryCount(cacheName, "Set")
				Expect(err).To(BeNil())
				Expect(retries > 1).To(BeTrue())
			})

			It("should not try to retry if the status code is not retryable", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.UnknownServiceError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Set")).To(Equal(0))
			})

			It("should not try to retry if the rpc is not retryable", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				strategy := retry.NewExponentialBackoffRetryStrategy(retry.ExponentialBackoffRetryStrategyProps{
					LoggerFactory:      momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
					InitialDelayMillis: 100,
					MaxBackoffMillis:   2000,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				incrResponse, err := cacheClient.DictionaryIncrement(context.Background(), &DictionaryIncrementRequest{
					CacheName:      cacheName,
					DictionaryName: "dictionary",
					Field:          String("field"),
					Amount:         1,
				})
				Expect(incrResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "DictionaryIncrement")).To(Equal(0))
			})
		})

		Describe("cache-client retry fixedCountRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should retry 3 times if the status code is retryable", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				retryStrategy := retry.NewFixedCountRetryStrategy(retry.FixedCountRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
					MaxAttempts:   3,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(getResponse).To(BeNil())

				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(Equal(3))
				Expect(metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")).To(BeNumerically("<=", 10))
			})

			It("should not retry if the status code is not retryable", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.UnknownServiceError)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)

				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Set")).To(Equal(0))
			})

			It("should not retry if the api is not retryable", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)

				incrementResponse, err := cacheClient.Increment(context.Background(), &IncrementRequest{
					CacheName: cacheName,
					Field:     String("key"),
				})
				Expect(incrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))

				dictCreateResponse, err := cacheClient.DictionarySetField(context.Background(), &DictionarySetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: "myDict",
					Field:          String("key"),
					Value:          String("value"),
					Ttl:            &utils.CollectionTtl{Ttl: 600 * time.Second},
				})
				Expect(dictCreateResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())

				dictIncrementResponse, err := cacheClient.DictionaryIncrement(context.Background(), &DictionaryIncrementRequest{
					CacheName:      cacheName,
					DictionaryName: "myDict",
					Field:          String("field"),
					Amount:         int64(1),
				})
				Expect(dictIncrementResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Increment")).To(Equal(0))
			})

			It("should return a value on success after a retry", func() {
				status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
				errCount := 1
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
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
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResponse).To(Not(BeNil()))

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResponse).To(Not(BeNil()))
				Expect(getResponse.(*responses.GetHit).ValueString()).To(Equal("value"))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(Equal(1))
			})
		})

		Describe("cache-client retry fixedTimeoutRetryStrategy", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
			It("should not retry if the status code is not retryable", func() {
				status := "unknown"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"set"},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)
				setResponse, err := cacheClient.Set(context.Background(), &SetRequest{
					CacheName: cacheName,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(setResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(UnknownServiceError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Set")).To(Equal(0))
			})

			It("should not retry if the rpc is not retryable", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"dictionary-increment"},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)
				incrResponse, err := cacheClient.DictionaryIncrement(context.Background(), &DictionaryIncrementRequest{
					CacheName:      cacheName,
					DictionaryName: "dictionary",
					Field:          String("field"),
					Amount:         1,
				})
				Expect(incrResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(ServerUnavailableError))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "DictionaryIncrement")).To(Equal(0))
			})

			It("should use default timeout values when not specified", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory: momento_default_logger.DefaultMomentoLoggerFactory{},
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))
				Expect(getResponse).To(BeNil())

				// Should immediately receive errors and retry every DefaultRetryDelayIntervalMillis
				// until the client timeout is reached.
				maxAttempts := CLIENT_TIMEOUT_MILLIS / retry.DefaultRetryDelayIntervalMillis
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				Expect(metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")).To(BeNumerically("<=", retry.DefaultRetryDelayIntervalMillis*1.1))
				Expect(metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")).To(BeNumerically(">=", retry.DefaultRetryDelayIntervalMillis*0.9))
			})

			It("should retry until client timeout when responses have no delays during full outage", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))

				// Should immediately receive errors and retry every DefaultRetryDelayIntervalMillis
				// until the client timeout is reached.
				maxAttempts := CLIENT_TIMEOUT_MILLIS / RETRY_DELAY_INTERVAL_MILLIS
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 1.1
				minDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")
				Expect(err).To(BeNil())
				Expect(average).To(BeNumerically("<=", int64(maxDelay)))
				Expect(average).To(BeNumerically(">=", int64(minDelay)))
			})

			It("should retry until client timeout when responses have short delays during full outage", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				shortDelay := RETRY_DELAY_INTERVAL_MILLIS + 100
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						DelayRpcList: &[]string{"get"},
						DelayMillis:  &shortDelay,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))

				// Should receive errors after shortDelay ms and retry every RETRY_DELAY_INTERVAL_MILLIS
				// until the client timeout is reached.
				delayBetweenAttempts := RETRY_DELAY_INTERVAL_MILLIS + shortDelay
				maxAttempts := int(CLIENT_TIMEOUT_MILLIS.Milliseconds()) / delayBetweenAttempts
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(delayBetweenAttempts) * 1.1
				minDelay := float64(delayBetweenAttempts) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")
				Expect(err).To(BeNil())
				Expect(float64(average)).To(BeNumerically("<=", maxDelay))
				Expect(float64(average)).To(BeNumerically(">=", minDelay))
			})

			It("should retry until client timeout when responses have long delays during full outage", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				longDelay := RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS + 100
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						DelayRpcList: &[]string{"get"},
						DelayMillis:  &longDelay,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))

				// Should receive errors after longDelay ms and retry every RETRY_DELAY_INTERVAL_MILLIS
				// until the client timeout is reached.
				delayBetweenAttempts := RETRY_DELAY_INTERVAL_MILLIS + longDelay
				maxAttempts := int(CLIENT_TIMEOUT_MILLIS.Milliseconds()) / delayBetweenAttempts
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically("<=", maxAttempts))
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(BeNumerically(">", 0))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(delayBetweenAttempts) * 1.1
				minDelay := float64(delayBetweenAttempts) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")
				Expect(err).To(BeNil())
				Expect(float64(average)).To(BeNumerically("<=", maxDelay))
				Expect(float64(average)).To(BeNumerically(">=", minDelay))
			})

			It("should retry until partial outage is resolved", func() {
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: RESPONSE_DATA_RECEIVED_TIMEOUT_MILLIS,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				errCount := 3
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						ErrorCount:   &errCount,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(CLIENT_TIMEOUT_MILLIS)
				setupCacheClient(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(getResponse).To(Not(BeNil()))
				Expect(err).To(BeNil())

				// Should retry until the server stops returning errors
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(Equal(errCount))

				// Jitter will be +/- 10% of the retry delay interval
				maxDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 1.1
				minDelay := float64(RETRY_DELAY_INTERVAL_MILLIS) * 0.9
				average, err := metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")
				Expect(err).To(BeNil())
				Expect(average).To(BeNumerically("<=", int64(maxDelay)))
				Expect(average).To(BeNumerically(">=", int64(minDelay)))
			})

			It("should not exceed client timeout when retry timeout is greater than client timeout", func() {
				response_data_received_timeout_millis := 3000
				client_timeout_millis := 2000
				response_delay := 1000
				status := "unavailable"
				retryStrategy := retry.NewFixedTimeoutRetryStrategy(retry.FixedTimeoutRetryStrategyProps{
					LoggerFactory:                     momento_default_logger.DefaultMomentoLoggerFactory{},
					ResponseDataReceivedTimeoutMillis: response_data_received_timeout_millis,
					RetryDelayIntervalMillis:          RETRY_DELAY_INTERVAL_MILLIS,
				})
				retryMiddleware := helpers.NewMomentoLocalMiddleware(helpers.MomentoLocalMiddlewareProps{
					MomentoLocalMiddlewareMetadataProps: helpers.MomentoLocalMiddlewareMetadataProps{
						ReturnError:  &status,
						ErrorRpcList: &[]string{"get"},
						DelayRpcList: &[]string{"get"},
						DelayMillis:  &response_delay,
					},
				})
				metricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetMetricsCollector()
				clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
					retryMiddleware,
				}).WithRetryStrategy(retryStrategy).WithClientTimeout(time.Duration(client_timeout_millis) * time.Millisecond)
				setupCacheClient(clientConfig)

				getResponse, err := cacheClient.Get(context.Background(), &GetRequest{
					CacheName: cacheName,
					Key:       String("key"),
				})
				Expect(getResponse).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(err).To(HaveMomentoErrorCode(TimeoutError))

				// Should retry once and retry attempt should not exceedclient timeout
				Expect(metricsCollector.GetTotalRetryCount(cacheName, "Get")).To(Equal(1))
				Expect(metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")).To(BeNumerically("<=", int64(client_timeout_millis)))
				Expect(metricsCollector.GetAverageTimeBetweenRetries(cacheName, "Get")).To(BeNumerically(">", int64(0)))
			})
		})
	})

	Describe("Topic Subscription Reconnects", func() {

		BeforeEach(func() {
			testCtx = context.Background()
			cacheName = uuid.NewString()
			topicName = uuid.NewString()
			setupCacheClient(config.LaptopLatest())
		})

		AfterEach(func() {
			cleanup()
		})

		It("should reconnect on recoverable error", func() {
			msgLimit := 9
			status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.ServerUnavailableError)
			clientConfig, retryMiddleware := getClientConfig(&clientConfigProps{
				status:                  &status,
				streamErrorMessageLimit: &msgLimit,
				streamErrorRpcList:      &[]string{"topic-subscribe"},
			})
			topicClient := setupTopicClient(clientConfig)

			publishedValues := make([]TopicValue, 0)
			for i := 0; i < 10; i++ {
				publishedValues = append(publishedValues, String(fmt.Sprintf("aaa%02d", i)))
			}

			err := doPubSub(topicClient, publishedValues)
			Expect(err).To(BeNil())

			topicEventMetricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetTopicEventCollector()
			counter, err := topicEventMetricsCollector.GetEventCounter(cacheName, "Subscribe")

			Expect(err).To(BeNil())
			Expect(counter.Errors).To(Equal(1))
			Expect(counter.Reconnects).To(Equal(1))
			Expect(counter.Items > 0).To(BeTrue())
			Expect(counter.Heartbeats > 0).To(BeTrue())
			Expect(counter.Discontinuities).To(Equal(0))
		})

		It("should not reconnect on unrecoverable error", func() {
			msgLimit := 8
			status := helpers.MomentoErrorCodeToMomentoLocalMetadataValue(momentoerrors.CanceledError)
			clientConfig, retryMiddleware := getClientConfig(&clientConfigProps{
				status:                  &status,
				streamErrorMessageLimit: &msgLimit,
				streamErrorRpcList:      &[]string{"topic-subscribe"},
			})
			topicClient := setupTopicClient(clientConfig)
			publishedValues := make([]TopicValue, 0)
			for i := 0; i < 10; i++ {
				publishedValues = append(publishedValues, String(fmt.Sprintf("aaa%02d", i)))
			}

			err := doPubSub(topicClient, publishedValues)
			Expect(err).To(BeNil())

			topicEventMetricsCollector := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetTopicEventCollector()
			counter, err := topicEventMetricsCollector.GetEventCounter(cacheName, "Subscribe")

			Expect(err).To(BeNil())
			Expect(counter.Errors).To(Equal(1))
			Expect(counter.Reconnects).To(Equal(0))
			Expect(counter.Items > 0).To(BeTrue())
			Expect(counter.Heartbeats > 0).To(BeTrue())
			Expect(counter.Discontinuities).To(Equal(0))
		})

		It("should timeout if deadline exceeds client timeout on first message", func() {
			delayMillis := 10_000
			clientConfig, _ := getClientConfig(&clientConfigProps{
				delayRpcList: &[]string{"topic-subscribe"},
				delayMillis:  &delayMillis,
			})
			topicClient := setupTopicClient(clientConfig)

			publishedValues := make([]TopicValue, 0)
			for i := 0; i < 10; i++ {
				publishedValues = append(publishedValues, String(fmt.Sprintf("aaa%02d", i)))
			}

			err := doPubSub(topicClient, publishedValues)
			Expect(err).To(HaveMomentoErrorCode(TimeoutError))
		})
	})

	Describe("Network Outage", func() {
		BeforeEach(func() {
			testCtx = context.Background()
			cacheName = uuid.NewString()
			topicName = uuid.NewString()
			setupCacheClient(config.LaptopLatest())
		})

		AfterEach(func() {
			cleanup()
		})

		It("should pause subscription when admin port is blocked and resume subscription once admin port is unblocked", func() {
			clientConfig, retryMiddleware := getClientConfig(&clientConfigProps{})
			topicClient := setupTopicClient(clientConfig.WithTransportStrategy(
				clientConfig.GetTransportStrategy().WithClientTimeout(time.Duration(5) * time.Minute)))
			sub, err := topicClient.Subscribe(testCtx, &TopicSubscribeRequest{
				CacheName: cacheName,
				TopicName: topicName,
			})
			Expect(err).To(BeNil())

			publishedValues := make([]TopicValue, 0)
			for i := 0; i < 10; i++ {
				publishedValues = append(publishedValues, String(fmt.Sprintf("aaa%02d", i)))
			}

			cancelContext, cancelFunction := context.WithCancel(testCtx)
			ready := make(chan int, 1)
			go func() {
				ready <- 1
				for {
					select {
					case <-cancelContext.Done():
						return
					default:
						_, err := sub.Item(cancelContext)
						if err != nil {
							fmt.Printf("error receiving item: %v\n", err)
							return
						}
					}
				}
			}()
			<-ready

			topicEventCounter := *retryMiddleware.(helpers.MomentoLocalMiddleware).GetTopicEventCollector()
			time.Sleep(time.Millisecond * 1000)
			numItemsAtBlock := 0
			for idx, value := range publishedValues {
				_, err := topicClient.Publish(testCtx, &TopicPublishRequest{
					CacheName: cacheName,
					TopicName: topicName,
					Value:     value,
				})
				if err != nil {
					panic(err)
				}

				if idx == 5 {
					testAdminHost, ok := os.LookupEnv("TEST_ADMIN_ENDPOINT")
					if !ok {
						testAdminHost = "127.0.0.1"
					}
					testAdminPort, ok := os.LookupEnv("TEST_ADMIN_PORT")
					if !ok {
						testAdminPort = "9090"
					}
					testAdminUrl := fmt.Sprintf("http://%s:%s/", testAdminHost, testAdminPort)
					_, err := http.Get(fmt.Sprintf("%s/block", testAdminUrl))
					Expect(err).To(BeNil())
					counter, err := topicEventCounter.GetEventCounter(cacheName, "Subscribe")
					Expect(err).To(BeNil())
					numItemsAtBlock = counter.Items
					Expect(numItemsAtBlock >= 5).To(BeTrue())
					time.Sleep(time.Millisecond * 2000)
					Expect(numItemsAtBlock).To(Equal(counter.Items))
					_, err = http.Get(fmt.Sprintf("%s/unblock", testAdminUrl))
					Expect(err).To(BeNil())
				}
			}
			time.Sleep(time.Millisecond * 5000)
			cancelFunction()
			counter, err := topicEventCounter.GetEventCounter(cacheName, "Subscribe")
			Expect(err).To(BeNil())
			Expect(counter.Items > numItemsAtBlock).To(BeTrue())
		})
	})
})
