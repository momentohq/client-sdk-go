package helpers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"google.golang.org/grpc/metadata"
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
	fmt.Printf("adding timestamp for %s: %d\n", timestamp, cacheName)
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
//
//	Limited to second resolution, but I can obviously change that if desired.
//	This tracks with the JS implementation.
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

// TODO: make these thread safe
type TopicEventCounter struct {
	Heartbeats int
	Discontinuities int
	Items int
	Errors int
	Reconnects int
}

type topicEventMetrics struct {
	data map[string]map[string]TopicEventCounter
}

type TopicEventMetricsCollector interface {
	AddEvent(cacheName string, requestName string, event middleware.TopicSubscriptionEventType)
	GetEventCounter(cacheName string, requestName string) (*TopicEventCounter, error)
}

func NewTopicEventMetricsCollector() TopicEventMetricsCollector {
	return &topicEventMetrics{
		data: make(map[string]map[string]TopicEventCounter),
	}
}

func (t *topicEventMetrics) GetEventCounter(cacheName string, requestName string) (*TopicEventCounter, error) {
	if _, ok := t.data[cacheName]; !ok {
		return nil, fmt.Errorf("cache name '%s' is not valid", cacheName)
	}
	if _, ok := t.data[cacheName][requestName]; !ok {
		return nil, fmt.Errorf("request name '%s' is not valid", requestName)
	}
	counter := t.data[cacheName][requestName]
	return &counter, nil
}

func (t *topicEventMetrics) initializeEventMap(cacheName string, requestName string) {
	if _, ok := t.data[cacheName]; !ok {
		t.data[cacheName] = make(map[string]TopicEventCounter)
	}
	if _, ok := t.data[cacheName][requestName]; !ok {
		t.data[cacheName][requestName] = TopicEventCounter{
			Heartbeats: 0,
			Discontinuities: 0,
			Items: 0,
			Errors: 0,
			Reconnects: 0,
		}
	}
}

func (t *topicEventMetrics) AddEvent(cacheName string, requestName string, event middleware.TopicSubscriptionEventType) {
	t.initializeEventMap(cacheName, requestName)
	eventCounter := t.data[cacheName][requestName]

	switch event {
	case middleware.HEARTBEAT:
		eventCounter.Heartbeats++
	case middleware.DISCONTINUITY:
		eventCounter.Discontinuities++
	case middleware.ITEM:
		eventCounter.Items++
	case middleware.RECONNECT:
		eventCounter.Reconnects++
	case middleware.ERROR:
		eventCounter.Errors++
	}

	t.data[cacheName][requestName] = eventCounter
}

type MomentoLocalMiddlewareProps struct {
	middleware.Props
	MomentoLocalMiddlewareRequestHandlerProps
	MomentoLocalMiddlewareTopicProps
}

type momentoLocalMiddleware struct {
	middleware.Middleware
	id uuid.UUID
	metricsCollector    RetryMetricsCollector
	topicEventMetricsCollector TopicEventMetricsCollector
	metricsChan         chan *timestampPayload
	requestHandlerProps MomentoLocalMiddlewareRequestHandlerProps
	topicProps MomentoLocalMiddlewareTopicProps
}

type MomentoLocalMiddleware interface {
	middleware.Middleware
	middleware.TopicMiddleware
	GetMetricsCollector() *RetryMetricsCollector
	GetTopicEventCollector() TopicEventMetricsCollector
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
	topicEventMetricsCollector := NewTopicEventMetricsCollector()
	mw := &momentoLocalMiddleware{
		Middleware:          baseMw,
		id:                  uuid.New(),
		metricsCollector:    metricsCollector,
		topicEventMetricsCollector: topicEventMetricsCollector,
		metricsChan:         metricsChan,
		requestHandlerProps: props.MomentoLocalMiddlewareRequestHandlerProps,
		topicProps:          props.MomentoLocalMiddlewareTopicProps,
	}

	// launch goroutine to listen for metrics from the unary interceptor callback
	go mw.listenForMetrics(metricsChan)

	return mw
}

func (mw *momentoLocalMiddleware) GetMetricsCollector() *RetryMetricsCollector {
	return &mw.metricsCollector
}

func (mw *momentoLocalMiddleware) GetTopicEventCollector() TopicEventMetricsCollector {
	return mw.topicEventMetricsCollector
}

func (mw *momentoLocalMiddleware) listenForMetrics(metricsChan chan *timestampPayload) {
	for {
		msg := <-metricsChan
		// this shouldn't happen under normal circumstances, but I thought it would
		// be good to provide a way to stop the goroutine.
		if msg == nil {
			return
		}
		// All requests are prefixed with "/cache_client.Scs/", so we cut that off.
		parts := strings.Split(msg.requestName, "/")
		if len(parts) < 2 {
			// Because this middleware is for test use only, we can panic here.
			panic(fmt.Sprintf("Could not parse request name %s", msg.requestName))
		}
		shortRequestName := parts[2]
		mw.metricsCollector.AddTimestamp(msg.cacheName, shortRequestName, msg.timestamp)
	}
}

func (mw *momentoLocalMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewRetryMetricsMiddlewareRequestHandler(
		baseHandler,
		mw.id,
		mw.metricsChan,
		mw.requestHandlerProps,
	), nil
}

func (mw *momentoLocalMiddleware) OnInterceptorRequest(ctx context.Context, method string) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		fmt.Printf("interceptor reqest got metadata: %#v\n", md)
		mw.metricsChan <- &timestampPayload{
			cacheName:   md.Get("cache")[0],
			requestName: method,
			timestamp:   time.Now().Unix(),
		}
	} else {
		// Because this middleware is for test use only, we can panic here.
		panic(fmt.Sprintf("no metadata found in context: %#v", ctx))
	}
}

func (mw *momentoLocalMiddleware) OnTopicEvent(cacheName string, requestName string, event middleware.TopicSubscriptionEventType) {
	mw.topicEventMetricsCollector.AddEvent(cacheName, requestName, event)
}

type momentoLocalMiddlewareRequestHandler struct {
	middleware.RequestHandler
	middlewareId uuid.UUID
	metricsChan chan *timestampPayload
	props       MomentoLocalMiddlewareRequestHandlerProps
}

type MomentoLocalMiddlewareRequestHandlerProps struct {
	ReturnError             *string
	ErrorRpcList            *[]string
	ErrorCount              *int
	DelayRpcList            *[]string
	DelayMillis             *int
	DelayCount              *int
}

type MomentoLocalMiddlewareTopicProps struct {
	StreamErrorRpcList      *[]string
	StreamError             *string
	StreamErrorMessageLimit *int
}

func NewRetryMetricsMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	id uuid.UUID,
	metricsChan chan *timestampPayload,
	props MomentoLocalMiddlewareRequestHandlerProps,
) middleware.RequestHandler {
	return &momentoLocalMiddlewareRequestHandler{
		RequestHandler: rh, middlewareId: id, metricsChan: metricsChan, props: props}
}

func (mw *momentoLocalMiddleware) OnSubscribeMetadata(requestMetadata map[string]string) map[string]string {
	requestMetadata["request-id"] = mw.id.String()

	if mw.topicProps.StreamErrorRpcList != nil {
		requestMetadata["stream-error-rpcs"] = strings.Join(*mw.topicProps.StreamErrorRpcList, " ")
	}

	if mw.topicProps.StreamError != nil {
		requestMetadata["stream-error"] = *mw.topicProps.StreamError
	}

	if mw.topicProps.StreamErrorMessageLimit != nil {
		requestMetadata["stream-error-message-limit"] = fmt.Sprintf("%d", *mw.topicProps.StreamErrorMessageLimit)
	}

	return requestMetadata
}

func (mw *momentoLocalMiddleware) OnPublishMetadata(requestMetadata map[string]string) map[string]string {
	// request-id is a little misleading-- this is actually more of a session id
	requestMetadata["request-id"] = mw.id.String()
	return requestMetadata
}

func (rh *momentoLocalMiddlewareRequestHandler) OnMetadata(requestMetadata map[string]string) map[string]string {
	// request-id is a little misleading-- this is actually more of a session id
	requestMetadata["request-id"] = rh.middlewareId.String()

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
