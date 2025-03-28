package helpers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"google.golang.org/grpc/metadata"
)

type MomentoLocalMiddlewareProps struct {
	middleware.Props
	// MomentoLocalMiddlewareMetadataProps holds properties from which the request handlers and topic middleware
	// will create metadata to be sent to the server to control the behavior of the request.
	MomentoLocalMiddlewareMetadataProps
}

type MomentoLocalMiddlewareMetadataProps struct {
	StreamErrorRpcList      *[]string
	StreamError             *string
	StreamErrorMessageLimit *int
	ReturnError             *string
	ErrorRpcList            *[]string
	ErrorCount              *int
	DelayRpcList            *[]string
	DelayMillis             *int
	DelayCount              *int
}

type momentoLocalMiddleware struct {
	middleware.Middleware
	id uuid.UUID
	metricsCollector    RetryMetricsCollector
	topicEventMetricsCollector TopicEventMetricsCollector
	metricsChan         chan *timestampPayload
	metadataProps MomentoLocalMiddlewareMetadataProps
	topicEventChan chan *topicEventPayload
}

// MomentoLocalMiddleware implements both the Middleware and TopicMiddleware interfaces. As such, instantiating one
// fires both goroutines to listen for both retry and topic event data.
type MomentoLocalMiddleware interface {
	middleware.Middleware
	middleware.TopicMiddleware
	GetMetricsCollector() *RetryMetricsCollector
	GetTopicEventCollector() *TopicEventMetricsCollector
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
		metadataProps: props.MomentoLocalMiddlewareMetadataProps,
		topicEventChan:      topicEventChan,
	}

	// launch goroutine to listen for metrics from the unary interceptor callback
	go mw.listenForMetrics(metricsChan)
	// launch goroutine to listen for topic events from the topic middleware callback
	go mw.listenForTopicEvents(topicEventChan)
	return mw
}

// MomentoLocalMiddleware interface methods

func (mw *momentoLocalMiddleware) GetMetricsCollector() *RetryMetricsCollector {
	return &mw.metricsCollector
}

func (mw *momentoLocalMiddleware) GetTopicEventCollector() *TopicEventMetricsCollector {
	return &mw.topicEventMetricsCollector
}

// middleware.Middleware interface methods

func (mw *momentoLocalMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewMomentoLocalMiddlewareRequestHandler(
		baseHandler,
		mw.id,
		mw.metricsChan,
		mw.metadataProps,
	), nil
}

// middleware.TopicMiddleware interface methods

func (mw *momentoLocalMiddleware) OnSubscribeMetadata(requestMetadata map[string]string) map[string]string {
	// request-id is actually a session id and must be the same throughout a given test.
	requestMetadata["request-id"] = mw.id.String()

	if mw.metadataProps.StreamErrorRpcList != nil {
		requestMetadata["stream-error-rpcs"] = strings.Join(*mw.metadataProps.StreamErrorRpcList, " ")
	}

	if mw.metadataProps.StreamError != nil {
		requestMetadata["stream-error"] = *mw.metadataProps.StreamError
	}

	if mw.metadataProps.StreamErrorMessageLimit != nil {
		requestMetadata["stream-error-message-limit"] = fmt.Sprintf("%d", *mw.metadataProps.StreamErrorMessageLimit)
	}

	if mw.metadataProps.DelayCount != nil {
		requestMetadata["delay-count"] = fmt.Sprintf("%d", *mw.metadataProps.DelayCount)
	}

	if mw.metadataProps.DelayMillis != nil {
		requestMetadata["delay-ms"] = fmt.Sprintf("%d", *mw.metadataProps.DelayMillis)
	}

	if mw.metadataProps.DelayRpcList != nil {
		requestMetadata["delay-rpcs"] = strings.Join(*mw.metadataProps.DelayRpcList, " ")
	}

	return requestMetadata
}

func (mw *momentoLocalMiddleware) OnPublishMetadata(requestMetadata map[string]string) map[string]string {
	// request-id is actually a session id and must be the same throughout a given test.
	requestMetadata["request-id"] = mw.id.String()

	if mw.metadataProps.ReturnError != nil {
		requestMetadata["return-error"] = *mw.metadataProps.ReturnError
	}

	if mw.metadataProps.ErrorRpcList != nil {
		requestMetadata["error-rpcs"] = strings.Join(*mw.metadataProps.ErrorRpcList, " ")
	}

	if mw.metadataProps.ErrorCount != nil {
		requestMetadata["error-count"] = fmt.Sprintf("%d", *mw.metadataProps.ErrorCount)
	}

	if mw.metadataProps.DelayCount != nil {
		requestMetadata["delay-count"] = fmt.Sprintf("%d", *mw.metadataProps.DelayCount)
	}

	if mw.metadataProps.DelayMillis != nil {
		requestMetadata["delay-ms"] = fmt.Sprintf("%d", *mw.metadataProps.DelayMillis)
	}

	if mw.metadataProps.DelayRpcList != nil {
		requestMetadata["delay-rpcs"] = strings.Join(*mw.metadataProps.DelayRpcList, " ")
	}

	return requestMetadata
}

// middleware.InterceptorCallbackMiddleware interface methods

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

// middleware.TopicEventCallbackMiddleware interface methods

func (mw *momentoLocalMiddleware) OnTopicEvent(cacheName string, requestName string, event middleware.TopicSubscriptionEventType) {
	mw.topicEventChan <- &topicEventPayload{
		cacheName:   cacheName,
		requestName: requestName,
		eventType:   event,
	}
}

// listeners for the goroutines that receive metrics and topic events on their respective channels

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

func MomentoErrorCodeToMomentoLocalMetadataValue(errorCode string) string {
	switch errorCode {
	case momentoerrors.InvalidArgumentError:
		return "invalid-argument"
	case momentoerrors.UnknownServiceError:
		return "unknown"
	case momentoerrors.AlreadyExistsError, momentoerrors.StoreAlreadyExistsError:
		return "already-exists"
	case momentoerrors.NotFoundError, momentoerrors.StoreNotFoundError:
		return "not-found"
	case momentoerrors.InternalServerError:
		return "internal"
	case momentoerrors.PermissionError:
		return "permission-denied"
	case momentoerrors.AuthenticationError:
		return "unauthenticated"
	case momentoerrors.CanceledError:
		return "cancelled"
	case momentoerrors.ConnectionError:
		return "unavailable"
	case momentoerrors.LimitExceededError:
		return "resource-exhausted"
	case momentoerrors.BadRequestError:
		return "invalid-argument"
	case momentoerrors.TimeoutError:
		return "deadline-exceeded"
	case momentoerrors.ServerUnavailableError:
		return "unavailable"
	case momentoerrors.FailedPreconditionError:
		return "failed-precondition"
	default:
		return "unknown"
	}
}
