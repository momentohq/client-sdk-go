package main

import (
	"context"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"time"

	"github.com/loov/hrtime"
)

type timingMiddleware struct {
	middleware.Middleware
}

func (mw *timingMiddleware) GetRequestHandler() middleware.RequestHandler {
	return NewTimingMiddlewareRequestHandler(middleware.HandlerProps{Logger: mw.GetLogger(), IncludeTypes: mw.GetIncludeTypes()})
}

func NewTimingMiddleware(props middleware.Props) middleware.Middleware {
	mw := middleware.NewMiddleware(props)
	return &timingMiddleware{mw}
}

type timingMiddlewareRequestHandler struct {
	middleware.RequestHandler
	startTime time.Duration
}

func NewTimingMiddlewareRequestHandler(props middleware.HandlerProps) middleware.RequestHandler {
	rh := middleware.NewRequestHandler(props)
	return &timingMiddlewareRequestHandler{rh, 0}
}

func (rh *timingMiddlewareRequestHandler) OnRequest(theRequest interface{}, metadata context.Context)  error {
	err := rh.RequestHandler.OnRequest(theRequest, metadata)
	if err != nil {
		return err
	}
	rh.startTime = hrtime.Now()
	return nil
}

func (rh *timingMiddlewareRequestHandler) OnResponse(theResponse interface{}) {
	elapsed := hrtime.Since(rh.startTime)
	rh.GetLogger().Info("%T took %s", rh.GetRequest(), elapsed)
}

