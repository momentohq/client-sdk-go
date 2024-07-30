package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/momentohq/client-sdk-go/config/logger"
)

type Middleware interface {
	OnRequest(theRequest interface{}, metadata context.Context)
	OnResponse(theResponse interface{})
}

type LoggingMiddleware struct {
	Log logger.MomentoLogger
}

func (mw *LoggingMiddleware) OnRequest(theRequest interface{}, metadata context.Context) {
	// Log request
	jsonStr, _ := json.MarshalIndent(theRequest, "", "  ")
	mw.Log.Info(fmt.Sprintf("\nIssuing %T:\n%s\nwith metadada: %+v\n", theRequest, string(jsonStr), metadata))
}

func (mw *LoggingMiddleware) OnResponse(theResponse interface{}) {
	// Log response
	jsonStr, _ := json.MarshalIndent(theResponse, "", "  ")
	mw.Log.Info(fmt.Sprintf("\nGot %T: %+v\n", theResponse, string(jsonStr)))
}
