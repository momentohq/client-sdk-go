package main

import (
	"encoding/json"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/middleware"
)

type loggingMiddleware struct {
	middleware.Middleware
}

func (mw *loggingMiddleware) GetRequestHandler(baseHandler middleware.RequestHandler) (middleware.RequestHandler, error) {
	return NewLoggingMiddlewareRequestHandler(baseHandler), nil
}

func NewLoggingMiddleware(props middleware.Props) middleware.Middleware {
	mw := middleware.NewMiddleware(props)
	return &loggingMiddleware{mw}
}

type loggingMiddlewareRequestHandler struct {
	middleware.RequestHandler
}

func NewLoggingMiddlewareRequestHandler(rh middleware.RequestHandler) middleware.RequestHandler {
	return &loggingMiddlewareRequestHandler{rh}
}

func (rh *loggingMiddlewareRequestHandler) OnRequest() {
	jsonStr, _ := json.MarshalIndent(rh.GetRequest(), "", "  ")
	rh.GetLogger().Info(
		fmt.Sprintf(
			"\n(%s) Issuing %T:\n%s\nwith metadada: %+v\n",
			rh.GetId(), rh.GetRequest(), string(jsonStr), rh.GetMetadata(),
		),
	)
}

func (rh *loggingMiddlewareRequestHandler) OnResponse(theResponse interface{}) {
	jsonStr, _ := json.MarshalIndent(theResponse, "", "  ")
	rh.GetLogger().Info(fmt.Sprintf("\n(%s) Got response: %s\n", rh.GetId(), string(jsonStr)))
}
