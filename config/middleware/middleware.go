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
	GetRequestHandler() RequestHandler
	GetIncludeTypes() map[string]bool
}

type Props struct {
	Logger       logger.MomentoLogger
	IncludeTypes []interface{}
}

func (mw *middleware) GetLogger() logger.MomentoLogger {
	return mw.logger
}

func (mw *middleware) GetRequestHandler() RequestHandler {
	return mw.requestHandler
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
	includeTypes map[string]bool
	metadata     context.Context
}

type RequestHandler interface {
	GetId() uuid.UUID
	GetRequest() interface{}
	GetLogger() logger.MomentoLogger
	GetIncludeTypes() map[string]bool
	OnRequest(theRequest interface{}, metadata context.Context) error
	OnResponse(theResponse interface{})
}

type HandlerProps struct {
	Logger       logger.MomentoLogger
	IncludeTypes map[string]bool
}

func (rh *requestHandler) GetId() uuid.UUID {
	return rh.id
}

func (rh *requestHandler) GetRequest() interface{} {
	return rh.request
}

func (rh *requestHandler) GetLogger() logger.MomentoLogger {
	return rh.logger
}

func (rh *requestHandler) GetIncludeTypes() map[string]bool {
	return rh.includeTypes
}

// OnRequest checks if the request type is one that the middleware is configured to handle.
// If it is, it sets the request and metadata fields on the request handler. Otherwise, it
// returns an error. Middleware request handlers may call this method in their OnRequest
// implementations to ensure they are only handling requests they are configured to handle.
// If the request handlers return the error, they will be omitted from request and response handling.
func (rh *requestHandler) OnRequest(theRequest interface{}, metadata context.Context) error {
	allowedTypes := rh.GetIncludeTypes()
	if allowedTypes != nil  {
		requestType := reflect.TypeOf(theRequest)
		if _, ok := allowedTypes[requestType.String()]; !ok {
			return fmt.Errorf("request type %T not in includeTypes", theRequest)
		}
	}
	rh.request = theRequest
	rh.metadata = metadata
	return nil
}

func (rh *requestHandler) OnResponse(_ interface{}) {}

func NewRequestHandler(props HandlerProps) RequestHandler {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return &requestHandler{id: id, logger: props.Logger, includeTypes: props.IncludeTypes}
}
