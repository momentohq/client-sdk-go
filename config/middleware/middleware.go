package middleware

import (
	"context"
)

type Middleware interface {
	OnRequest(requestId uint64, theRequest interface{}, metadata context.Context)
	OnResponse(requestId uint64, theResponse map[string]string)
}
