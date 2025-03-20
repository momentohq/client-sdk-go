package helpers

import (
	"context"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/logger"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/config/retry"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type timestampPayload struct {
	cacheName   string
	requestName string
	timestamp   int64
}

type retryMetrics struct {
	data map[string]map[string][]int64
}

type RetryMetricsCollector interface {
	AddTimestamp(cacheName string, requestName string, timestamp int64)
	GetTotalRetryCount(cacheName string, requestName string) (int, error)
	GetAverageTimeBetweenRetries(cacheName string, requestName string) (int64, error)
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

func (r *retryMetrics) GetTotalRetryCount(cacheName string, requestName string) (int, error) {
	if _, ok := r.data[cacheName]; !ok {
		return 0, fmt.Errorf("cache name '%s' is not valid", cacheName)
	}
	if timestamps, ok := r.data[cacheName][requestName]; ok {
		// The first timestamp is the original request, so we subtract 1
		return len(timestamps) - 1, nil
	}
	return 0, fmt.Errorf("request name '%s' is not valid", requestName)
}

// GetAverageTimeBetweenRetries returns the average time between retries in seconds.
// TODO: what resolution are we looking for here? I'm using Unix epoch time, so am currently
//
//	limited to seconds, but I can obviously change that.
func (r *retryMetrics) GetAverageTimeBetweenRetries(cacheName string, requestName string) (int64, error) {
	if _, ok := r.data[cacheName]; !ok {
		return int64(0), fmt.Errorf("cache name '%s' is not valid", cacheName)
	}
	if timestamps, ok := r.data[cacheName][requestName]; ok {
		if len(timestamps) < 2 {
			return 0, nil
		}
		var sum int64
		for i := 1; i < len(timestamps); i++ {
			sum += timestamps[i] - timestamps[i-1]
		}
		return sum / int64(len(timestamps)-1), nil
	}
	return 0, fmt.Errorf("request name '%s' is not valid", requestName)
}

func (r *retryMetrics) GetAllMetrics() map[string]map[string][]int64 {
	return r.data
}

type RetryMetricsMiddlewareProps struct {
	middleware.Props
	RetryMetricsMiddlewareRequestHandlerProps
}

type retryMetricsMiddleware struct {
	middleware.Middleware
	metricsCollector RetryMetricsCollector
	metricsChan chan *timestampPayload
	requestHandlerProps RetryMetricsMiddlewareRequestHandlerProps
}

type RetryMetricsMiddleware interface {
	middleware.Middleware
	GetMetricsCollector() *RetryMetricsCollector
}

func NewRetryMetricsMiddleware(props RetryMetricsMiddlewareProps) middleware.Middleware {
	var myLogger logger.MomentoLogger
	if props.Logger == nil {
		myLogger = momento_default_logger.NewDefaultMomentoLoggerFactory(
			momento_default_logger.INFO).GetLogger("retry-metrics")
	} else {
		myLogger = props.Logger
	}
	baseMw := middleware.NewMiddleware(middleware.Props{
		Logger: myLogger,
		IncludeTypes: props.IncludeTypes,
	})
	metricsCollector := NewRetryMetricsCollector()
	metricsChan := make(chan *timestampPayload, 1000)
	mw := &retryMetricsMiddleware{
		Middleware: baseMw,
		metricsCollector: metricsCollector,
		metricsChan: metricsChan,
		requestHandlerProps: props.RetryMetricsMiddlewareRequestHandlerProps,
	}
	go mw.listenForMetrics(metricsChan)
	return mw
}

func (r *retryMetricsMiddleware) GetMetricsCollector() *RetryMetricsCollector {
	return &r.metricsCollector
}

func (r *retryMetricsMiddleware) listenForMetrics(metricsChan chan *timestampPayload) {
	for {
		select {
		case msg := <-metricsChan:
			// this shouldn't happen under normal circumstances, but I thought it would
			// be good to provide a way to stop the goroutine.
			if msg == nil {
				return
			}
			// All requests are prefixed with "/cache_client.Scs/", so we cut that off
			parsedRequest, ok := strings.CutPrefix(msg.requestName, "/cache_client.Scs/")
			if !ok {
				// Because this middleware is for test use only, we can panic here.
				panic(fmt.Sprintf("Could not parse request name %s", msg.requestName))
			}
			r.metricsCollector.AddTimestamp(msg.cacheName, parsedRequest, msg.timestamp)
		}
	}
}

func (r *retryMetricsMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewRetryMetricsMiddlewareRequestHandler(
		baseHandler,
		r.metricsChan,
		r.requestHandlerProps,
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
		for {
			// Execute api call
			md, ok := metadata.FromOutgoingContext(ctx); if ok {
				r.metricsChan <- &timestampPayload{
					cacheName:   md.Get("cache")[0],
					requestName: method,
					timestamp:   time.Now().Unix(),
				}
			} else {
				// Because this middleware is for test use only, we can panic here.
				panic(fmt.Sprintf("no metadata found in context: %#v", ctx))
			}
			lastErr := invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				// Success. No error was returned so we can return from the interceptor.
				return nil
			}

			// Check retry eligibility based off last error received
			retryBackoffTime := s.DetermineWhenToRetry(retry.StrategyProps{
				GrpcStatusCode: status.Code(lastErr),
				GrpcMethod:     method,
				AttemptNumber:  attempt,
			})

			if retryBackoffTime == nil {
				// Request is not retryable. Return the error.
				return lastErr
			}

			// Sleep for recommended time interval and increment attempts before trying again
			if *retryBackoffTime > 0 {
				time.Sleep(time.Duration(*retryBackoffTime) * time.Millisecond)
			}
			attempt++
		}
	}
}

type retryMetricsMiddlewareRequestHandler struct {
	middleware.RequestHandler
	metricsChan chan *timestampPayload
	props RetryMetricsMiddlewareRequestHandlerProps
}

type RetryMetricsMiddlewareRequestHandlerProps struct {
	ReturnError      *string
	ErrorRpcList     *[]string
	ErrorCount       *int
	DelayRpcList     *[]string
	DelayMillis      *int
	DelayCount       *int
}

func NewRetryMetricsMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	metricsChan chan *timestampPayload,
	props RetryMetricsMiddlewareRequestHandlerProps,
) middleware.RequestHandler {
	return &retryMetricsMiddlewareRequestHandler{RequestHandler: rh, metricsChan: metricsChan, props: props}
}

func (rh *retryMetricsMiddlewareRequestHandler) OnMetadata(requestMetadata map[string]string) map[string]string {
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

	if rh.props.DelayCount != nil {
		requestMetadata["delay-count"] = fmt.Sprintf("%d", *rh.props.DelayCount)
	}

	if rh.props.DelayMillis != nil {
		requestMetadata["delay-millis"] = fmt.Sprintf("%d", *rh.props.DelayMillis)
	}

	if rh.props.DelayRpcList != nil {
		requestMetadata["delay-rpcs"] = strings.Join(*rh.props.DelayRpcList, " ")
	}

	return requestMetadata
}
