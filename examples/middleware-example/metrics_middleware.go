package main

import (
	"encoding/json"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/middleware"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type metricsMiddleware struct {
	middleware.Middleware
	logEveryNRequests uint64
	requestChan chan string
}

func (mw *metricsMiddleware) GetRequestHandler(baseHandler middleware.RequestHandler) (middleware.RequestHandler, error) {
	return NewMetricsMiddlewareRequestHandler(baseHandler, mw.requestChan), nil
}

func NewMetricsMiddleware(props middleware.Props, logEveryNRequests uint64) middleware.Middleware {
	mw := middleware.NewMiddleware(props)
	requestChan := make(chan string, 100)
	go func() {
		metricsSink(requestChan, props.Logger, logEveryNRequests)
	}()
	return &metricsMiddleware{mw, logEveryNRequests, requestChan}
}

func metricsSink(requestChan chan string, log logger.MomentoLogger, logEveryNRequests uint64) {
	totalRequests := 0
	requestCount := make(map[string]uint64)
	for {
		requestMsg := <-requestChan
		totalRequests++
		requestCount[requestMsg]++
		if (totalRequests % int(logEveryNRequests)) == 0 {
			jsonStr, _ := json.MarshalIndent(requestCount, "", "  ")
			log.Info(fmt.Sprintf("Request count: %s", string(jsonStr)))
		}
	}
}

type metricsMiddlewareRequestHandler struct {
	middleware.RequestHandler
	requestChan chan string
}

func (rh *metricsMiddlewareRequestHandler) OnRequest() {
	rh.requestChan <- fmt.Sprintf("%T", rh.GetRequest())
}

func (rh *metricsMiddlewareRequestHandler) OnResponse(_ interface{}) {}

func NewMetricsMiddlewareRequestHandler(rh middleware.RequestHandler, requestChan chan string) middleware.RequestHandler {
	return &metricsMiddlewareRequestHandler{rh, requestChan}
}
