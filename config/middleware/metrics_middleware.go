package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type metricsMiddleware struct {
	Log         logger.MomentoLogger
	requestChan chan string
}

func NewMetricsMiddleware(log logger.MomentoLogger) *metricsMiddleware {
	mw := &metricsMiddleware{
		Log:         log,
		requestChan: make(chan string),
	}
	go func() {
		metricsSink(mw.requestChan, log)
	}()
	return mw
}

func metricsSink(requestChan chan string, log logger.MomentoLogger) {
	requestCount := make(map[string]uint64)
	for {
		select {
		case requestMsg := <-requestChan:
			requestCount[requestMsg]++
			jsonStr, _ := json.MarshalIndent(requestCount, "", "  ")
			log.Info(fmt.Sprintf("Request count: %s", string(jsonStr)))
		}
	}
}

func (mw *metricsMiddleware) OnRequest(_ uint64, theRequest interface{}, _ context.Context) {
	mw.requestChan <- fmt.Sprintf("%T", theRequest)
}

func (mw *metricsMiddleware) OnResponse(requestId uint64, _ map[string]string) {}
