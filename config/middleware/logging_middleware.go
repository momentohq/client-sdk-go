package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type loggingMiddleware struct {
	Middleware
}

func (mw *loggingMiddleware) GetRequestHandler() RequestHandler {
	return NewLoggingMiddlewareRequestHandler(mw.GetLogger())
}

func NewLoggingMiddleware(props Props) Middleware {
	mw := NewMiddleware(props)
	return &loggingMiddleware{mw}
}

type loggingMiddlewareRequestHandler struct {
	RequestHandler
}

func NewLoggingMiddlewareRequestHandler(log logger.MomentoLogger) RequestHandler {
	rh := NewRequestHandler(HandlerProps{Logger: log})
	return &loggingMiddlewareRequestHandler{rh}
}

func (rh *loggingMiddlewareRequestHandler) OnRequest(theRequest interface{}, metadata context.Context) error {
	err := rh.RequestHandler.OnRequest(theRequest, metadata)
	if err != nil {
		return err
	}
	// Logger request
	jsonStr, _ := json.MarshalIndent(theRequest, "", "  ")
	rh.GetLogger().Info(
		fmt.Sprintf(
			"\n(%d) Issuing %T:\n%s\nwith metadada: %+v\n", rh.GetId(), rh.GetRequest(), string(jsonStr), metadata,
		),
	)
	return nil;
}

func (rh *loggingMiddlewareRequestHandler) OnResponse(theResponse interface{}) {
	// Logger response
	jsonStr, _ := json.MarshalIndent(theResponse, "", "  ")
	rh.GetLogger().Info(fmt.Sprintf("\n(%d) Got response: %s\n", rh.GetRequest(), string(jsonStr)))
}
