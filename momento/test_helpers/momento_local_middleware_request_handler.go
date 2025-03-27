package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"strings"
)

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

func NewMomentoLocalMiddlewareRequestHandler(
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
