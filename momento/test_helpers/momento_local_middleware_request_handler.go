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
	metricsChan   chan *timestampPayload
	metadataProps MomentoLocalMiddlewareMetadataProps
}

func NewMomentoLocalMiddlewareRequestHandler(
	rh middleware.RequestHandler,
	id uuid.UUID,
	metricsChan chan *timestampPayload,
	props MomentoLocalMiddlewareMetadataProps,
) middleware.RequestHandler {
	return &momentoLocalMiddlewareRequestHandler{
		RequestHandler: rh, middlewareId: id, metricsChan: metricsChan, metadataProps: props}
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
		requestMetadata["delay-millis"] = fmt.Sprintf("%d", *rh.metadataProps.DelayMillis)
	}

	if rh.metadataProps.DelayRpcList != nil {
		requestMetadata["delay-rpcs"] = strings.Join(*rh.metadataProps.DelayRpcList, " ")
	}

	return requestMetadata
}
