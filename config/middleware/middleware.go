package middleware

import (
	"context"
)

type Middleware interface {
	OnRequest(theRequest interface{}, metadata context.Context)
	OnResponse(theResponse map[string]string)
}
