package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/momentohq/client-sdk-go/config/logger"
)

type LoggingMiddleware struct {
	Log logger.MomentoLogger
}

func NewLoggingMiddleware(log logger.MomentoLogger) *LoggingMiddleware {
	return &LoggingMiddleware{
		Log: log,
	}
}

func (mw *LoggingMiddleware) OnRequest(requestId uint64, theRequest interface{}, metadata context.Context) {
	// Log request
	jsonStr, _ := json.MarshalIndent(theRequest, "", "  ")
	mw.Log.Info(fmt.Sprintf("\n(%d) Issuing %T:\n%s\nwith metadada: %+v\n", requestId, theRequest, string(jsonStr), metadata))
}

func (mw *LoggingMiddleware) OnResponse(requestId uint64, theResponse map[string]string) {
	if len(theResponse) == 0 {
		mw.Log.Debug("Got empty response")
		return
	}
	// Log response
	jsonStr, _ := json.MarshalIndent(theResponse, "", "  ")
	mw.Log.Info(fmt.Sprintf("\n(%d) Got response: %s\n", requestId, string(jsonStr)))
}
