package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/momento_rpc_names"
	"google.golang.org/grpc/metadata"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
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

// GetAverageTimeBetweenRetries returns the average time between retries in milliseconds.
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

type MomentoLocalMiddlewareProps struct {
	middleware.Props
	MomentoLocalMiddlewareRequestHandlerProps
}

type momentoLocalMiddleware struct {
	middleware.Middleware
	metricsCollector    RetryMetricsCollector
	metricsChan         chan *timestampPayload
	requestHandlerProps MomentoLocalMiddlewareRequestHandlerProps
}

type MomentoLocalMiddleware interface {
	middleware.Middleware
	GetMetricsCollector() *RetryMetricsCollector
}

func NewMomentoLocalMiddleware(props MomentoLocalMiddlewareProps) middleware.Middleware {
	var myLogger logger.MomentoLogger
	if props.Logger == nil {
		myLogger = momento_default_logger.NewDefaultMomentoLoggerFactory(
			momento_default_logger.INFO).GetLogger("retry-metrics")
	} else {
		myLogger = props.Logger
	}
	baseMw := middleware.NewMiddleware(middleware.Props{
		Logger:       myLogger,
		IncludeTypes: props.IncludeTypes,
	})
	metricsCollector := NewRetryMetricsCollector()
	metricsChan := make(chan *timestampPayload, 1000)
	mw := &momentoLocalMiddleware{
		Middleware:          baseMw,
		metricsCollector:    metricsCollector,
		metricsChan:         metricsChan,
		requestHandlerProps: props.MomentoLocalMiddlewareRequestHandlerProps,
	}
	go mw.listenForMetrics(metricsChan)
	return mw
}

func (r *momentoLocalMiddleware) GetMetricsCollector() *RetryMetricsCollector {
	return &r.metricsCollector
}

func (r *momentoLocalMiddleware) listenForMetrics(metricsChan chan *timestampPayload) {
	for {
		msg := <-metricsChan
		// this shouldn't happen under normal circumstances, but I thought it would
		// be good to provide a way to stop the goroutine.
		if msg == nil {
			return
		}
		shortRequestName := ConvertRpcNameToMomentoLocalRpcName(momento_rpc_names.MomentoRPCMethod(msg.requestName))
		r.metricsCollector.AddTimestamp(msg.cacheName, shortRequestName, msg.timestamp)
	}
}

func (r *momentoLocalMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewRetryMetricsMiddlewareRequestHandler(
		baseHandler,
		r.metricsChan,
		r.requestHandlerProps,
	), nil
}

func (r *momentoLocalMiddleware) OnInterceptorRequest(ctx context.Context, method string) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		r.metricsChan <- &timestampPayload{
			cacheName:   md.Get("cache")[0],
			requestName: method,
			timestamp:   time.Now().UnixMilli(),
		}
	} else {
		// Because this middleware is for test use only, we can panic here.
		panic(fmt.Sprintf("no metadata found in context: %#v", ctx))
	}
}

type momentoLocalMiddlewareRequestHandler struct {
	middleware.RequestHandler
	metricsChan chan *timestampPayload
	props       MomentoLocalMiddlewareRequestHandlerProps
}

type MomentoLocalMiddlewareRequestHandlerProps struct {
	ReturnError  *string
	ErrorRpcList *[]string
	ErrorCount   *int
	DelayRpcList *[]string
	DelayMillis  *int
	DelayCount   *int
}

func NewRetryMetricsMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	metricsChan chan *timestampPayload,
	props MomentoLocalMiddlewareRequestHandlerProps,
) middleware.RequestHandler {
	return &momentoLocalMiddlewareRequestHandler{RequestHandler: rh, metricsChan: metricsChan, props: props}
}

func (rh *momentoLocalMiddlewareRequestHandler) OnMetadata(requestMetadata map[string]string) map[string]string {
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
		requestMetadata["delay-ms"] = fmt.Sprintf("%d", *rh.props.DelayMillis)
	}

	if rh.props.DelayRpcList != nil {
		requestMetadata["delay-rpcs"] = strings.Join(*rh.props.DelayRpcList, " ")
	}

	return requestMetadata
}
