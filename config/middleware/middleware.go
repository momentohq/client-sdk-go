package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/config/logger"
	"reflect"
)

type middleware struct {
	logger         logger.MomentoLogger
	requestHandler RequestHandler
	includeTypes   map[string]bool
}

type Middleware interface {
	GetLogger() logger.MomentoLogger
	GetBaseRequestHandler(theRequest interface{}, metadata context.Context) (RequestHandler, error)
	GetRequestHandler(baseRequestHandler RequestHandler) (RequestHandler, error)
	GetIncludeTypes() map[string]bool
}

type Props struct {
	Logger       logger.MomentoLogger
	IncludeTypes []interface{}
}

func (mw *middleware) GetBaseRequestHandler(theRequest interface{}, metadata context.Context) (RequestHandler, error) {
	allowedTypes := mw.GetIncludeTypes()
	if allowedTypes != nil {
		requestType := reflect.TypeOf(theRequest)
		if _, ok := allowedTypes[requestType.String()]; !ok {
			return nil, fmt.Errorf("request type %T not in includeTypes", theRequest)
		}
	}

	// Return the "base" request handler. User request handlers will be composed on top of this.
	return NewRequestHandler(
		HandlerProps {
			Metadata: metadata,
			Request: theRequest,
			Logger: mw.GetLogger(),
			IncludeTypes: mw.GetIncludeTypes(),
		},
	), nil
}

func (mw *middleware) GetLogger() logger.MomentoLogger {
	return mw.logger
}

func (mw *middleware) GetRequestHandler(_ RequestHandler) (RequestHandler, error) {
	return nil, fmt.Errorf("GetRequestHandler not implemented in middleware")
}

func (mw *middleware) GetIncludeTypes() map[string]bool {
	return mw.includeTypes
}

// NewMiddleware creates a new middleware with a logger and optional list of request types it should handle.
// If the IncludeTypes are omitted or empty, all request types will be processed. For example, to limit processing
// to only requests of type *momento.SetRequest and *momento.GetRequest, pass the following slice as the IncludeTypes:
//   []interface{}{&momento.SetRequest{}, &momento.GetRequest{}}
func NewMiddleware(props Props) Middleware {
	// convert the slice of types to a map of type names for quick lookup in the data client
	var includeTypeMap map[string]bool
	if props.IncludeTypes != nil {
		includeTypeMap = make(map[string]bool)
		for _, t := range props.IncludeTypes {
			includeTypeMap[reflect.TypeOf(t).String()] = true
		}
	} else {
		includeTypeMap = nil
	}
	if props.Logger == nil {
		props.Logger = logger.NewNoopMomentoLoggerFactory().GetLogger("noop")
	}
	return &middleware{logger: props.Logger, includeTypes: includeTypeMap}
}

type requestHandler struct {
	id           uuid.UUID
	logger       logger.MomentoLogger
	request      interface{}
	metadata    context.Context
	includeTypes map[string]bool
}

type RequestHandler interface {
	GetId() uuid.UUID
	GetRequest() interface{}
	GetMetadata() context.Context
	GetLogger() logger.MomentoLogger
	GetIncludeTypes() map[string]bool
	OnRequest()
	OnResponse(theResponse interface{})
}

type HandlerProps struct {
	Request		 interface{}
	Metadata     context.Context
	Logger       logger.MomentoLogger
	IncludeTypes map[string]bool
}

func (rh *requestHandler) GetId() uuid.UUID {
	return rh.id
}

func (rh *requestHandler) GetRequest() interface{} {
	return rh.request
}

func (rh *requestHandler) GetMetadata() context.Context {
	return rh.metadata
}

func (rh *requestHandler) GetLogger() logger.MomentoLogger {
	return rh.logger
}

func (rh *requestHandler) GetIncludeTypes() map[string]bool {
	return rh.includeTypes
}

func (rh *requestHandler) OnRequest() {}

func (rh *requestHandler) OnResponse(_ interface{}) {}

func NewRequestHandler(props HandlerProps) RequestHandler {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return &requestHandler{
		id: id,
		logger: props.Logger,
		includeTypes: props.IncludeTypes,
		request: props.Request,
		metadata: props.Metadata,
	}
}
