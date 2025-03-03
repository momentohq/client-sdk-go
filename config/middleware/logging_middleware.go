package middleware

import (
	"context"
	"encoding/json"
	"fmt"
)

type loggingMiddleware struct {
	Middleware
}

func (mw *loggingMiddleware) GetRequestHandler() RequestHandler {
	return NewLoggingMiddlewareRequestHandler(HandlerProps{Logger: mw.GetLogger(), IncludeTypes: mw.GetIncludeTypes()})
}

func NewLoggingMiddleware(props Props) Middleware {
	mw := NewMiddleware(props)
	return &loggingMiddleware{mw}
}

type loggingMiddlewareRequestHandler struct {
	RequestHandler
}

func NewLoggingMiddlewareRequestHandler(props HandlerProps) RequestHandler {
	rh := NewRequestHandler(props)
	return &loggingMiddlewareRequestHandler{rh}
}

func (rh *loggingMiddlewareRequestHandler) OnRequest(theRequest interface{}, metadata context.Context) error {
	err := rh.RequestHandler.OnRequest(theRequest, metadata)
	if err != nil {
		return err
	}
	jsonStr, _ := json.MarshalIndent(theRequest, "", "  ")
	rh.GetLogger().Info(
		fmt.Sprintf(
			"\n(%d) Issuing %T:\n%s\nwith metadada: %+v\n", rh.GetId(), rh.GetRequest(), string(jsonStr), metadata,
		),
	)
	return nil;
}

func (rh *loggingMiddlewareRequestHandler) OnResponse(theResponse interface{}) {
	jsonStr, _ := json.MarshalIndent(theResponse, "", "  ")
	rh.GetLogger().Info(fmt.Sprintf("\n(%d) Got response: %s\n", rh.GetRequest(), string(jsonStr)))
}
