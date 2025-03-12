package momento_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/momentohq/client-sdk-go/internal/retry"
	"github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"time"
)

type ResponseWithRetryTimestamps struct {
	Err error
	Reply interface{}
	Timestamps  []int64
}

func (e ResponseWithRetryTimestamps) Error() string {
	return e.Err.Error()
}

type ReplyWithRetryTimestamps struct {
	Reply interface{}
	Timestamps []int64
}

type retryMetrics struct {
	data map[string]map[string][]int64
}

type RetryMetricsCollector interface {
	addTimestamp(cacheName string, requestName string, timestamp int64)
	getTotalRetryCount(cacheName string, requestName string) int
	getAverageTimeBetweenRetries(cacheName string, requestName string) int64
	getAllMetrics() map[string]map[string][]int64
}

func NewRetryMetricsCollector() RetryMetricsCollector {
	return &retryMetrics{data: make(map[string]map[string][]int64)}
}

func (r *retryMetrics) addTimestamp(cacheName string, requestName string, timestamp int64) {
	if _, ok := r.data[cacheName]; !ok {
		r.data[cacheName] = make(map[string][]int64)
	}
	r.data[cacheName][requestName] = append(r.data[cacheName][requestName], timestamp)
}

func (r *retryMetrics) getTotalRetryCount(cacheName string, requestName string) int {
	// The first timestamp is the original request, so we subtract 1
	return len(r.data[cacheName][requestName]) - 1
}

// TODO: what resolution are we looking for here? I'm using Unix epoch time, so am currently
//  limited to seconds, but I can obviously change that.
func (r *retryMetrics) getAverageTimeBetweenRetries(cacheName string, requestName string) int64 {
	if timestamps, ok := r.data[cacheName][requestName]; !ok {
		return 0
	} else {
		if len(timestamps) < 2 {
			return 0
		}
		var sum int64
		for i := 1; i < len(timestamps); i++ {
			sum += timestamps[i] - timestamps[i-1]
		}
		return sum / int64(len(timestamps) - 1)
	}
}

func (r *retryMetrics) getAllMetrics() map[string]map[string][]int64 {
	return r.data
}

type retryMetricsMiddlewareProps struct {
	metricsCollector RetryMetricsCollector
	returnError *string
	errorRpcList *[]string
	errorCount *int
	delayRpcList *[]string
	delayMillis *int
	delayCount *int
}

type retryMetricsMiddleware struct {
	middleware.Middleware
	props    retryMetricsMiddlewareProps
}

func NewRetryMetricsMiddleware(props retryMetricsMiddlewareProps) middleware.Middleware {
	mw := middleware.NewMiddleware(middleware.Props{
		Logger: momento_default_logger.NewDefaultMomentoLoggerFactory(
			momento_default_logger.INFO).GetLogger("retry-metrics"),
	})
	return &retryMetricsMiddleware{Middleware: mw, props: props}
}

type retryMetricsMiddlewareRequestHandler struct {
	middleware.RequestHandler
	cacheName string
	props   retryMetricsMiddlewareProps
}

func (r *retryMetricsMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewRetryMetricsMiddlewareRequestHandler(
		baseHandler,
		r.props,
	), nil
}

func (r *retryMetricsMiddleware) AddUnaryRetryInterceptor(s retry.Strategy) func(
	ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		attempt := 1
		timestamps := make([]int64, 0)
		for {
			fmt.Printf("=====\nretry interceptor attempt %d for %s\n", attempt, method)
			// Execute api call
			timestamps = append(timestamps, time.Now().Unix())
			lastErr := invoker(ctx, method, req, reply, cc, opts...)
			fmt.Printf("----> lastErr: %v\n", lastErr)
			if lastErr == nil {
				// Success no error returned stop interceptor
				// TODO: if we have retries, we need to return them in an error response
				//  that has a nil error and a metadata field with the retry timestamps
				fmt.Printf("----> timestamps: %v\n", timestamps)
				if len(timestamps) > 1 {
					fmt.Printf("retry interceptor got success after %d retries; creating ReplyWithRetryTimestamps\n", len(timestamps) - 1)
					return ResponseWithRetryTimestamps{
						Reply:      reply,
						Timestamps: timestamps,
					}
				}
				return nil
			}

			// Check retry eligibility based off last error received
			retryBackoffTime := s.DetermineWhenToRetry(retry.StrategyProps{
				GrpcStatusCode: status.Code(lastErr),
				GrpcMethod:     method,
				AttemptNumber:  attempt,
			})

			if retryBackoffTime == nil {
				// Stop retrying and return the error with the timestamps
				return ResponseWithRetryTimestamps{
					Err: lastErr,
					Timestamps: timestamps,
				}
			}

			// Sleep for recommended time interval and increment attempts before trying again
			if *retryBackoffTime > 0 {
				time.Sleep(time.Duration(*retryBackoffTime) * time.Second)
			}
			attempt++
		}
	}
}

func NewRetryMetricsMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	props retryMetricsMiddlewareProps,
) middleware.RequestHandler {
	return &retryMetricsMiddlewareRequestHandler{RequestHandler: rh, cacheName: "", props: props}
}

func (rh *retryMetricsMiddlewareRequestHandler) OnResponse(theResponse interface{}, err error) (interface{}, error) {
	if err != nil {
		fmt.Printf("retryMetricsMiddlewareRequestHandler.OnResponse got error: %T\n", err)
		var e ResponseWithRetryTimestamps
		switch {
		case errors.As(err, &e):
			fmt.Printf("Reply -> %T - %v\n", e.Reply, e.Reply)
			fmt.Printf("Timestamps -> %v\n", e.Timestamps)
			for _, ts := range e.Timestamps {
				rh.props.metricsCollector.addTimestamp(rh.GetResourceName(), rh.GetRequestName(), ts)
			}
			if e.Err != nil {
				return nil, e.Err
			} else {
				return e.Reply, nil
			}
		default:
			return nil, err
		}
	}
	return theResponse, nil
}

func (rh *retryMetricsMiddlewareRequestHandler) GetMetadata() map[string]string {
	requestMetadata := rh.RequestHandler.GetMetadata()
	requestMetadata["request-id"] = rh.GetId().String()
	if rh.props.returnError != nil {
		requestMetadata["return-error"] = *rh.props.returnError
	}
	if rh.props.errorRpcList != nil {
		requestMetadata["error-rpcs"] = strings.Join(*rh.props.errorRpcList, " ")
	}
	if rh.props.errorCount != nil {
		requestMetadata["error-count"] = fmt.Sprintf("%d", *rh.props.errorCount)
	}

	// add other stuff here :-)/

	return requestMetadata
}

func setupCacheClientTest(config config.Configuration) momento.CacheClient {
	credentialProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{})
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

var _ = Describe("cache-client retry eligibility-strategy", Label(CACHE_SERVICE_LABEL), func() {
	DescribeTable(
		"DefaultEligibilityStrategy -- determine retry eligibility given grpc status code and request method",
		func(grpcStatus codes.Code, requestMethod string, expected bool) {
			strategy := retry.NewFixedCountRetryStrategy(momento_default_logger.NewDefaultMomentoLoggerFactory(
				momento_default_logger.INFO,
			))
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

	Describe("cache-client retry NeverRetryStrategy", Label(CACHE_SERVICE_LABEL), func() {
		It("shouldn't retry", func() {
			metricsCollector := NewRetryMetricsCollector()
			status := "unavailable"
			strategy := retry.NewNeverRetryStrategy()
			clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				NewRetryMetricsMiddleware(retryMetricsMiddlewareProps{
					metricsCollector: metricsCollector,
					returnError:      &status,
					errorRpcList:     &[]string{"set"},
					errorCount:       nil,
					delayRpcList:     nil,
					delayMillis:      nil,
					delayCount:       nil,
				}),
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
			Expect(metricsCollector.getTotalRetryCount("cache", "Set")).To(Equal(0))
		})
	})

	Describe("cache-client retry DefaultEligibilityStrategy", Label(CACHE_SERVICE_LABEL), func() {
		It("should retry 3 times if the status code is retryable", func() {
			metricsCollector := NewRetryMetricsCollector()
			status := "unavailable"
			clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				NewRetryMetricsMiddleware(retryMetricsMiddlewareProps{
					metricsCollector: metricsCollector,
					returnError:      &status,
					errorRpcList:     &[]string{"get"},
					errorCount:       nil,
					delayRpcList:     nil,
					delayMillis:      nil,
					delayCount:       nil,
				}),
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
			Expect(err).To(Not(BeNil()))
			Expect(err).To(HaveMomentoErrorCode(momento.ServerUnavailableError))
			Expect(getResponse).To(BeNil())

			Expect(metricsCollector.getTotalRetryCount("cache", "Get")).To(Equal(3))
			Expect(metricsCollector.getAverageTimeBetweenRetries("cache", "Get")).To(Equal(int64(0)))
		})

		It("should not retry if the status code is not retryable", func() {
			metricsCollector := NewRetryMetricsCollector()
			status := "unknown"
			clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				NewRetryMetricsMiddleware(retryMetricsMiddlewareProps{
					metricsCollector: metricsCollector,
					returnError:      &status,
					errorRpcList:     &[]string{"set"},
					errorCount:       nil,
					delayRpcList:     nil,
					delayMillis:      nil,
					delayCount:       nil,
				}),
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
			Expect(metricsCollector.getTotalRetryCount("cache", "Set")).To(Equal(0))
		})

		It("should return a value on success after a retry", func() {
			metricsCollector := NewRetryMetricsCollector()
			status := "unavailable"
			errCount := 1
			clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				NewRetryMetricsMiddleware(retryMetricsMiddlewareProps{
					metricsCollector: metricsCollector,
					returnError:      &status,
					errorRpcList:     &[]string{"get"},
					errorCount:       &errCount,
					delayRpcList:     nil,
					delayMillis:      nil,
					delayCount:       nil,
				}),
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
			fmt.Printf("%v", metricsCollector.getAllMetrics())
			Expect(err).To(BeNil())
			Expect(getResponse).To(Not(BeNil()))
			Expect(getResponse.(*responses.GetHit).ValueString()).To(Equal("value"))
			Expect(metricsCollector.getTotalRetryCount("cache", "Get")).To(Equal(1))
		})

	})
})
