package middleware

import (
	"fmt"
	"sync/atomic"
)

type inFlightRequestCountMiddleware struct {
	Middleware
	numActiveRequests atomic.Int64
}

func (mw *inFlightRequestCountMiddleware) GetRequestHandler(
	baseHandler RequestHandler,
) (RequestHandler, error) {
	return NewInFlightRequestCountMiddlewareRequestHandler(
		baseHandler,
		mw.add,
		mw.remove,
	), nil
}

func NewInFlightRequestCountMiddleware(props Props) Middleware {
	mw := NewMiddleware(props)
	return &inFlightRequestCountMiddleware{mw, atomic.Int64{}}
}

func (mw *inFlightRequestCountMiddleware) add() int64 {
	return mw.numActiveRequests.Add(1)
}

func (mw *inFlightRequestCountMiddleware) remove() int64 {
	return mw.numActiveRequests.Add(-1)
}

type inFlightRequestCountMiddlewareRequestHandler struct {
	RequestHandler
	requestsAtStart int64
	remover         func() int64
}

func NewInFlightRequestCountMiddlewareRequestHandler(
	rh RequestHandler, adder func() int64, remover func() int64,
) RequestHandler {
	countAfterAdd := adder()
	return &inFlightRequestCountMiddlewareRequestHandler{rh, countAfterAdd, remover}
}

func (rh *inFlightRequestCountMiddlewareRequestHandler) OnResponse(_ interface{}, _ error) error {
	countAfterRemove := rh.remover()
	rh.GetLogger().Info(
		fmt.Sprintf(
			"Request %s completed; in-flight requests at start: %d, in-flight requests at completion: %d",
			rh.GetRequestName(),
			rh.requestsAtStart,
			countAfterRemove,
		),
	)
	return nil
}
