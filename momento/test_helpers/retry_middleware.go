package helpers

import (
	"context"
	"errors"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/retry"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ResponseWithRetryTimestamps struct {
	Err        error
	Reply      interface{}
	Timestamps []int64
}

func (e ResponseWithRetryTimestamps) Error() string {
	return e.Err.Error()
}

type ReplyWithRetryTimestamps struct {
	Reply      interface{}
	Timestamps []int64
}

type retryMetrics struct {
	data map[string]map[string][]int64
}

type RetryMetricsCollector interface {
	AddTimestamp(cacheName string, requestName string, timestamp int64)
	GetTotalRetryCount(cacheName string, requestName string) int
	GetAverageTimeBetweenRetries(cacheName string, requestName string) int64
	GetAllMetrics() map[string]map[string][]int64
}

func NewRetryMetricsCollector() RetryMetricsCollector {
	return &retryMetrics{data: make(map[string]map[string][]int64)}
}

func (r *retryMetrics) AddTimestamp(cacheName string, requestName string, timestamp int64) {
	if _, ok := r.data[cacheName]; !ok {
		r.data[cacheName] = make(map[string][]int64)
	}
	r.data[cacheName][requestName] = append(r.data[cacheName][requestName], timestamp)
}

func (r *retryMetrics) GetTotalRetryCount(cacheName string, requestName string) int {
	// The first timestamp is the original request, so we subtract 1
	return len(r.data[cacheName][requestName]) - 1
}

// GetAverageTimeBetweenRetries returns the average time between retries in seconds.
// TODO: what resolution are we looking for here? I'm using Unix epoch time, so am currently
//
//	limited to seconds, but I can obviously change that.
func (r *retryMetrics) GetAverageTimeBetweenRetries(cacheName string, requestName string) int64 {
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
		return sum / int64(len(timestamps)-1)
	}
}

func (r *retryMetrics) GetAllMetrics() map[string]map[string][]int64 {
	return r.data
}

type RetryMetricsMiddlewareProps struct {
	MetricsCollector RetryMetricsCollector
	ReturnError      *string
	ErrorRpcList     *[]string
	ErrorCount       *int
	DelayRpcList     *[]string
	DelayMillis      *int
	DelayCount       *int
}

type retryMetricsMiddleware struct {
	middleware.Middleware
	props RetryMetricsMiddlewareProps
}

func NewRetryMetricsMiddleware(props RetryMetricsMiddlewareProps) middleware.Middleware {
	mw := middleware.NewMiddleware(middleware.Props{
		Logger: momento_default_logger.NewDefaultMomentoLoggerFactory(
			momento_default_logger.INFO).GetLogger("retry-metrics"),
	})
	return &retryMetricsMiddleware{Middleware: mw, props: props}
}

type retryMetricsMiddlewareRequestHandler struct {
	middleware.RequestHandler
	props     RetryMetricsMiddlewareProps
}

func (r *retryMetricsMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewRetryMetricsMiddlewareRequestHandler(
		baseHandler,
		r.props,
	), nil
}

// AddUnaryRetryInterceptor returns a unary interceptor that will retry the request based on the retry strategy. It is
// essentially identical to the SDK retry interceptor, but it also collects metrics on the retry attempts and returns
// them in a custom error type. If modifications are made to the retry interceptor in the SDK, they should be mirrored
// here. Injecting the retry interceptor is a bit of a hack, but it keeps the retry metrics gathering out of the production
// code.
func (r *retryMetricsMiddleware) AddUnaryRetryInterceptor(s retry.Strategy) func(
	ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		attempt := 1
		timestamps := make([]int64, 0)
		for {
			// Execute api call
			timestamps = append(timestamps, time.Now().Unix())
			lastErr := invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				// Success. No error was returned so we can return from the interceptor.
				// If we have retries, we need to return them in an error response
				// that has a nil error and a metadata field with the retry timestamps.
				if len(timestamps) > 1 {
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
				// Stop retrying and return the error with the timestamps.
				return ResponseWithRetryTimestamps{
					Err:        lastErr,
					Timestamps: timestamps,
				}
			}

			// Sleep for recommended time interval and increment attempts before trying again
			if *retryBackoffTime > 0 {
				time.Sleep(time.Duration(*retryBackoffTime) * time.Millisecond)
			}
			attempt++
		}
	}
}

func NewRetryMetricsMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	props RetryMetricsMiddlewareProps,
) middleware.RequestHandler {
	return &retryMetricsMiddlewareRequestHandler{RequestHandler: rh, props: props}
}

func (rh *retryMetricsMiddlewareRequestHandler) OnResponse(theResponse interface{}, err error) (interface{}, error) {
	if err != nil {
		var e ResponseWithRetryTimestamps
		switch {
		case errors.As(err, &e):
			for _, ts := range e.Timestamps {
				rh.props.MetricsCollector.AddTimestamp(rh.GetResourceName(), rh.GetRequestName(), ts)
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
	if rh.props.ReturnError != nil {
		requestMetadata["return-error"] = *rh.props.ReturnError
	}
	if rh.props.ErrorRpcList != nil {
		requestMetadata["error-rpcs"] = strings.Join(*rh.props.ErrorRpcList, " ")
	}
	if rh.props.ErrorCount != nil {
		requestMetadata["error-count"] = fmt.Sprintf("%d", *rh.props.ErrorCount)
	}

	// add other stuff here :-)/

	return requestMetadata
}
