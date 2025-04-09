package impl

import (
	"fmt"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/middleware"
)

type inFlightRequestCountMiddleware struct {
	middleware.Middleware
	numActiveRequests atomic.Int64
}

func (mw *inFlightRequestCountMiddleware) GetRequestHandler(
	baseHandler middleware.RequestHandler,
) (middleware.RequestHandler, error) {
	return NewInFlightRequestCountMiddlewareRequestHandler(
		baseHandler,
		mw.add,
		mw.remove,
	), nil
}

func NewInFlightRequestCountMiddleware(props middleware.Props) middleware.Middleware {
	mw := middleware.NewMiddleware(props)
	return &inFlightRequestCountMiddleware{mw, atomic.Int64{}}
}

func (mw *inFlightRequestCountMiddleware) add() int64 {
	return mw.numActiveRequests.Add(1)
}

func (mw *inFlightRequestCountMiddleware) remove() int64 {
	return mw.numActiveRequests.Add(-1)
}

type inFlightRequestCountMiddlewareRequestHandler struct {
	middleware.RequestHandler
	requestsAtStart int64
	remover         func() int64
}

func NewInFlightRequestCountMiddlewareRequestHandler(
	rh middleware.RequestHandler, adder func() int64, remover func() int64,
) middleware.RequestHandler {
	countAfterAdd := adder()
	return &inFlightRequestCountMiddlewareRequestHandler{rh, countAfterAdd, remover}
}

func (rh *inFlightRequestCountMiddlewareRequestHandler) OnResponse(_ interface{}) (interface{}, error) {
	countAfterRemove := rh.remover()
	rh.GetLogger().Info(
		fmt.Sprintf(
			"Request %s completed; in-flight requests at start: %d, in-flight requests at completion: %d",
			rh.GetRequestName(),
			rh.requestsAtStart,
			countAfterRemove,
		),
	)
	return nil, nil
}
