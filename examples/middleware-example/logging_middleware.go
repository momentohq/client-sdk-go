package main

import (
	"encoding/json"
	"fmt"

	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/responses"
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

func (rh *loggingMiddlewareRequestHandler) OnRequest(_ interface{}) (interface{}, error) {
	jsonStr, _ := json.Marshal(rh.GetRequest())
	rh.GetLogger().Info(
		fmt.Sprintf(
			"(%s) Issuing %T: %s",
			rh.GetId(), rh.GetRequest(), string(jsonStr),
		),
	)
	return nil, nil
}

func (rh *loggingMiddlewareRequestHandler) OnResponse(theResponse interface{}) (interface{}, error) {
	switch r := theResponse.(type) {
	case *responses.GetHit:
		rh.GetLogger().Info(fmt.Sprintf("(%s) Got response: %T, %s", rh.GetId(), r, r.ValueString()))
	}
	return nil, nil
}
