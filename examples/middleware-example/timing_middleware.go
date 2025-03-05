package main

import (
	"github.com/momentohq/client-sdk-go/config/middleware"
	"time"

	"github.com/loov/hrtime"
)

type timingMiddleware struct {
	middleware.Middleware
}

func (mw *timingMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewTimingMiddlewareRequestHandler(baseHandler), nil
}

func NewTimingMiddleware(props middleware.Props) middleware.Middleware {
	mw := middleware.NewMiddleware(props)
	return &timingMiddleware{mw}
}

type timingMiddlewareRequestHandler struct {
	middleware.RequestHandler
	startTime time.Duration
}

func NewTimingMiddlewareRequestHandler(rh middleware.RequestHandler) middleware.RequestHandler {
	return &timingMiddlewareRequestHandler{rh, 0}
}

func (rh *timingMiddlewareRequestHandler) OnRequest() {
	rh.startTime = hrtime.Now()
}

func (rh *timingMiddlewareRequestHandler) OnResponse(_ interface{}) {
	elapsed := hrtime.Since(rh.startTime)
	rh.GetLogger().Info("%T took %s", rh.GetRequest(), elapsed)
}

