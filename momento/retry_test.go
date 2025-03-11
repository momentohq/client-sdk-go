package momento_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
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

type ErrorWithRetryTimestamps struct {
	Err error
	Timestamps  []int64
}

func (e ErrorWithRetryTimestamps) Error() string {
	return e.Err.Error()
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
	fmt.Printf("adding timestamp for cache: %s, request: %s, timestamp: %d\n", cacheName, requestName, timestamp)
	if _, ok := r.data[cacheName]; !ok {
		r.data[cacheName] = make(map[string][]int64)
	}
	r.data[cacheName][requestName] = append(r.data[cacheName][requestName], timestamp)
}

func (r *retryMetrics) getTotalRetryCount(cacheName string, requestName string) int {
	return len(r.data[cacheName][requestName])
}

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
		return sum / int64(len(timestamps)) - 1
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
			// Execute api call
			requestTimestamp := time.Now().Unix()
			lastErr := invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				// Success no error returned stop interceptor
				// TODO: if we have retries, we need to return them in an error response
				//  that has a nil error and a metadata field with the retry timestamps
				return nil
			}
			timestamps = append(timestamps, requestTimestamp)

			// Check retry eligibility based off last error received
			retryBackoffTime := s.DetermineWhenToRetry(retry.StrategyProps{
				GrpcStatusCode: status.Code(lastErr),
				GrpcMethod:     method,
				AttemptNumber:  attempt,
			})

			if retryBackoffTime == nil {
				// Stop retrying and return the error with the timestamps
				return ErrorWithRetryTimestamps{
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

func (rh *retryMetricsMiddlewareRequestHandler) OnRequest() {
	fmt.Printf("retry metrics middleware request handler on request with %s %s\n", rh.GetResourceType(), rh.GetResourceName())
}

func (rh *retryMetricsMiddlewareRequestHandler) OnResponse(_ interface{}, err error) error {
	if err != nil {
		switch e := err.(type) {
		case ErrorWithRetryTimestamps:
			for _, ts := range e.Timestamps {
				rh.props.metricsCollector.addTimestamp(rh.GetResourceName(), rh.GetRequestName(), ts)
			}
			if e.Err != nil {
				return e.Err
			}
		default:
			return err
		}
	}
	return nil
}

func (rh *retryMetricsMiddlewareRequestHandler) GetMetadata() map[string]string {
	requestMetadata := rh.RequestHandler.GetMetadata()
	fmt.Printf("======> retry metrics middleware request handler get metadata %+v\n", requestMetadata)
	requestMetadata["request-id"] = rh.GetId().String()
	if rh.props.returnError != nil {
		requestMetadata["return-error"] = *rh.props.returnError
	}
	if rh.props.errorRpcList != nil {
		requestMetadata["error-rpcs"] = strings.Join(*rh.props.errorRpcList, " ")
	}

	// add other stuff here :-)/

	return requestMetadata
}

var _ = Describe("cache-client retry eligibility-strategy", func() {
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

	Describe("DefaultEligibilityStrategy -- test retires with fixed count strategy", func() {
		It("should not retry if the status code is not retryable", func() {
			metricsCollector := NewRetryMetricsCollector()
			cancelledStatus := "unavailable"
			clientConfig := config.LaptopLatest().WithMiddleware([]middleware.Middleware{
				middleware.NewInFlightRequestCountMiddleware(middleware.Props{
					Logger: momento_default_logger.NewDefaultMomentoLoggerFactory(
						momento_default_logger.INFO).GetLogger("in-flight-request-count"),
				}),
				NewRetryMetricsMiddleware(retryMetricsMiddlewareProps{
					metricsCollector: metricsCollector,
					returnError:      &cancelledStatus,
					errorRpcList:     &[]string{"get"},
					errorCount:       nil,
					delayRpcList:     nil,
					delayMillis:      nil,
					delayCount:       nil,
				}),
			})
			credentialProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{})
			Expect(err).To(BeNil())
			cacheClient, err := momento.NewCacheClient(clientConfig, credentialProvider, 30*time.Second)
			Expect(err).To(BeNil())
			createResponse, err := cacheClient.CreateCache(context.Background(), &momento.CreateCacheRequest{
				CacheName: "cache",
			})
			Expect(err).To(BeNil())
			Expect(createResponse).To(Not(BeNil()))

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

			allMetrics, err := json.MarshalIndent(metricsCollector.getAllMetrics(), "", "  ")
			Expect(err).To(BeNil())
			fmt.Printf("%s\n", allMetrics)
		})
	})
})
