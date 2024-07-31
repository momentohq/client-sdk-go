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

func (mw *LoggingMiddleware) OnRequest(theRequest interface{}, metadata context.Context) {
	// Log request
	jsonStr, _ := json.MarshalIndent(theRequest, "", "  ")
	mw.Log.Info(fmt.Sprintf("\nIssuing %T:\n%s\nwith metadada: %+v\n", theRequest, string(jsonStr), metadata))
}

func (mw *LoggingMiddleware) OnResponse(theResponse map[string]string) {
	if len(theResponse) == 0 {
		mw.Log.Debug("Got empty response")
		return
	}
	// Log response
	jsonStr, _ := json.MarshalIndent(theResponse, "", "  ")
	mw.Log.Info(fmt.Sprintf("\nGot response: %s\n", string(jsonStr)))
}
