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
	topicEventChan chan *topicEventPayload
}

// MomentoLocalMiddleware implements both the Middleware and TopicMiddleware interfaces. As such, instantiating one
// fires both goroutines to listen for both retry and resubscribe data.
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
	topicEventChan := make(chan *topicEventPayload, 1000)
	mw := &momentoLocalMiddleware{
		Middleware:          baseMw,
		id:                  uuid.New(),
		metricsCollector:    metricsCollector,
		topicEventMetricsCollector: topicEventMetricsCollector,
		metricsChan:         metricsChan,
		requestHandlerProps: props.MomentoLocalMiddlewareRequestHandlerProps,
		topicProps:          props.MomentoLocalMiddlewareTopicProps,
		topicEventChan:      topicEventChan,
	}

	// launch goroutine to listen for metrics from the unary interceptor callback
	go mw.listenForMetrics(metricsChan)
	go mw.listenForTopicEvents(topicEventChan)
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

func (mw *momentoLocalMiddleware) listenForTopicEvents(topicEventChan chan *topicEventPayload) {
	for {
		msg := <-topicEventChan
		if msg == nil {
			return
		}
		mw.topicEventMetricsCollector.AddEvent(msg.cacheName, msg.requestName, msg.eventType)
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
	mw.topicEventChan <- &topicEventPayload{
		cacheName:   cacheName,
		requestName: requestName,
		eventType:   event,
	}
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
