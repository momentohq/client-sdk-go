package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type metricsMiddleware struct {
	Middleware
	logEveryNRequests uint64
	requestChan chan string
}

func (mw *metricsMiddleware) GetRequestHandler() RequestHandler {
	return NewMetricsMiddlewareRequestHandler(HandlerProps{mw.GetLogger(), mw.GetIncludeTypes()}, mw.requestChan)
}

func NewMetricsMiddleware(props Props, logEveryNRequests uint64) Middleware {
	mw := NewMiddleware(props)
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
	RequestHandler
	requestChan chan string
}

func (rh *metricsMiddlewareRequestHandler) OnRequest(theRequest interface{}, _ context.Context) error {
	err := rh.RequestHandler.OnRequest(theRequest, nil)
	if err != nil {
		return err
	}
	rh.requestChan <- fmt.Sprintf("%T", theRequest)
	return nil
}

func (rh *metricsMiddlewareRequestHandler) OnResponse(_ interface{}) {}

func NewMetricsMiddlewareRequestHandler(props HandlerProps, requestChan chan string) RequestHandler {
	rh := NewRequestHandler(HandlerProps{Logger: props.Logger, IncludeTypes: props.IncludeTypes})
	return &metricsMiddlewareRequestHandler{rh, requestChan}
}
