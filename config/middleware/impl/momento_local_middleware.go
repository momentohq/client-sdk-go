package impl

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
)

// MomentoLocalMiddlewareProps holds properties from which the middleware will be instantiated.
type MomentoLocalMiddlewareProps struct {
	middleware.Props
	MomentoLocalMiddlewareMetadataProps
}

// MomentoLocalMiddlewareMetadataProps holds properties from which the request handlers
// will create metadata to be sent to the server to control the behavior of the request.
type MomentoLocalMiddlewareMetadataProps struct {
	// ReturnError is the error code to return for a request.
	//
	// Valid values are: "unavailable", "unknown", "already-exists", "not-found", "internal",
	// "permission-denied", "unauthenticated", "cancelled", "resource-exhausted", "invalid-argument",
	// "deadline-exceeded", "failed-precondition"
	ReturnError *string

	// ErrorRpcList is the list of RPCs for which to return an error.
	//
	// Valid values are: "get", "get-batch", "set", "set-batch", "set-if", "delete", "keys-exist", "increment",
	// "update-ttl", "item-get-ttl", "item-get-type", "dictionary-get", "dictionary-fetch", "dictionary-set",
	// "dictionary-increment", "dictionary-delete", "dictionary-length", "set-fetch", "set-sample", "set-union",
	// "set-difference", "set-contains", "set-length", "set-pop", "list-push-front", "list-push-back",
	// "list-pop-front", "list-pop-back", "list-remove", "list-fetch", "list-length", "list-concatenate-front",
	// "list-concatenate-back", "list-retain", "sorted-set-put", "sorted-set-fetch", "sorted-set-get-score",
	// "sorted-set-remove", "sorted-set-increment", "sorted-set-get-rank", "sorted-set-length",
	// "sorted-set-length-by-score", "topic-publish", "topic-subscribe"
	ErrorRpcList *[]string

	// ErrorCount is the number of times for which to return an error.
	ErrorCount *int

	// DelayRpcList is the list of RPCs to delay for.
	//
	// Valid values are: "get", "get-batch", "set", "set-batch", "set-if", "delete", "keys-exist", "increment",
	// "update-ttl", "item-get-ttl", "item-get-type", "dictionary-get", "dictionary-fetch", "dictionary-set",
	// "dictionary-increment", "dictionary-delete", "dictionary-length", "set-fetch", "set-sample", "set-union",
	// "set-difference", "set-contains", "set-length", "set-pop", "list-push-front", "list-push-back",
	// "list-pop-front", "list-pop-back", "list-remove", "list-fetch", "list-length", "list-concatenate-front",
	// "list-concatenate-back", "list-retain", "sorted-set-put", "sorted-set-fetch", "sorted-set-get-score",
	// "sorted-set-remove", "sorted-set-increment", "sorted-set-get-rank", "sorted-set-length",
	// "sorted-set-length-by-score", "topic-publish", "topic-subscribe"
	DelayRpcList *[]string

	// DelayMillis is the number of milliseconds for which to delay a response.
	DelayMillis *int

	// DelayCount is the number of times to delay a response.
	DelayCount *int

	// StreamErrorRpcList is the list of RPCs to return a stream error for.
	// Valid values are: "topic-subscribe"
	StreamErrorRpcList *[]string

	// StreamError is the error code to return for a stream error.
	StreamError *string

	// StreamErrorMessageLimit is the limit of messages to return for a stream error.
	StreamErrorMessageLimit *int
}

type momentoLocalMiddleware struct {
	middleware.Middleware
	id            uuid.UUID
	metadataProps MomentoLocalMiddlewareMetadataProps
}

// MomentoLocalMiddleware implements both the Middleware and TopicMiddleware interfaces.
type MomentoLocalMiddleware interface {
	middleware.Middleware
	middleware.TopicMiddleware
}

// NewMomentoLocalMiddleware creates a new MomentoLocalMiddleware instance.
func NewMomentoLocalMiddleware(props MomentoLocalMiddlewareProps) middleware.Middleware {
	var myLogger logger.MomentoLogger
	if props.Logger == nil {
		myLogger = momento_default_logger.NewDefaultMomentoLoggerFactory(
			momento_default_logger.INFO).GetLogger("momento-local-middleware")
	} else {
		myLogger = props.Logger
	}
	baseMw := middleware.NewMiddleware(middleware.Props{
		Logger:       myLogger,
		IncludeTypes: props.IncludeTypes,
	})
	mw := &momentoLocalMiddleware{
		Middleware:    baseMw,
		id:            uuid.New(),
		metadataProps: props.MomentoLocalMiddlewareMetadataProps,
	}
	return mw
}

// middleware.Middleware interface methods

func (mw *momentoLocalMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewMomentoLocalMiddlewareRequestHandler(
		baseHandler,
		mw.id,
		mw.metadataProps,
	), nil
}

type momentoLocalMiddlewareRequestHandler struct {
	middleware.RequestHandler
	middlewareId  uuid.UUID
	metadataProps MomentoLocalMiddlewareMetadataProps
}

func NewMomentoLocalMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	id uuid.UUID,
	props MomentoLocalMiddlewareMetadataProps,
) middleware.RequestHandler {
	return &momentoLocalMiddlewareRequestHandler{
		RequestHandler: rh,
		middlewareId:   id,
		metadataProps:  props,
	}
}

func (rh *momentoLocalMiddlewareRequestHandler) OnMetadata(requestMetadata map[string]string) map[string]string {
	// request-id is actually a session id and must be the same throughout a given test.
	requestMetadata["request-id"] = rh.middlewareId.String()

	if rh.metadataProps.ReturnError != nil {
		requestMetadata["return-error"] = *rh.metadataProps.ReturnError
	}

	if rh.metadataProps.ErrorRpcList != nil {
		requestMetadata["error-rpcs"] = strings.Join(*rh.metadataProps.ErrorRpcList, " ")
	}

	if rh.metadataProps.ErrorCount != nil {
		requestMetadata["error-count"] = fmt.Sprintf("%d", *rh.metadataProps.ErrorCount)
	}

	if rh.metadataProps.DelayCount != nil {
		requestMetadata["delay-count"] = fmt.Sprintf("%d", *rh.metadataProps.DelayCount)
	}

	if rh.metadataProps.DelayMillis != nil {
		requestMetadata["delay-ms"] = fmt.Sprintf("%d", *rh.metadataProps.DelayMillis)
	}

	if rh.metadataProps.DelayRpcList != nil {
		requestMetadata["delay-rpcs"] = strings.Join(*rh.metadataProps.DelayRpcList, " ")
	}

	return requestMetadata
}

// middleware.TopicMiddleware interface methods

func (mw *momentoLocalMiddleware) OnSubscribeMetadata(requestMetadata map[string]string) map[string]string {
	// Subscribe shares most metadata with publish but adds some streaming config.
	requestMetadata = mw.OnPublishMetadata(requestMetadata)

	if mw.metadataProps.StreamErrorRpcList != nil {
		requestMetadata["stream-error-rpcs"] = strings.Join(*mw.metadataProps.StreamErrorRpcList, " ")
	}

	if mw.metadataProps.StreamError != nil {
		requestMetadata["stream-error"] = *mw.metadataProps.StreamError
	}

	if mw.metadataProps.StreamErrorMessageLimit != nil {
		requestMetadata["stream-error-message-limit"] = fmt.Sprintf("%d", *mw.metadataProps.StreamErrorMessageLimit)
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
