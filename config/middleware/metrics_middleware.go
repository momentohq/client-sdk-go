package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/logger"
)

type MetricsMiddleware struct {
	Log          logger.MomentoLogger
	requestCount map[string]uint64
}

func NewMetricsMiddleware(log logger.MomentoLogger) *MetricsMiddleware {
	return &MetricsMiddleware{
		Log: log,
	}
}

func (mw *MetricsMiddleware) OnRequest(theRequest interface{}, _ context.Context) {
	// Log request
	if mw.requestCount == nil {
		mw.requestCount = make(map[string]uint64)
	}
	mw.requestCount[fmt.Sprintf("%T", theRequest)]++
	jsonStr, _ := json.MarshalIndent(mw.requestCount, "", "  ")
	mw.Log.Info(fmt.Sprintf("%+v\n", string(jsonStr)))
}

func (mw *MetricsMiddleware) OnResponse(_ map[string]string) {
	return
}
