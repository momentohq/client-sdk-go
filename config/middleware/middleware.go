package middleware

import (
	"context"
	"fmt"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/retry"
	"google.golang.org/grpc"
	"reflect"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/config/logger"
)

type middleware struct {
	logger       logger.MomentoLogger
	includeTypes map[string]bool
}

type Middleware interface {
	GetLogger() logger.MomentoLogger
	GetBaseRequestHandler(theRequest interface{}, requestName string, resourceType internal.ClientType, resourceName string, metadata map[string]string) (RequestHandler, error)
	GetRequestHandler(baseRequestHandler RequestHandler) (RequestHandler, error)
	GetIncludeTypes() map[string]bool
}

type RetryMiddleware interface {
	Middleware
	AddUnaryRetryInterceptor(s retry.Strategy) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
}

type Props struct {
	Logger       logger.MomentoLogger
	IncludeTypes []interface{}
}

func (mw *middleware) GetBaseRequestHandler(
	theRequest interface{}, requestName string, resourceType internal.ClientType, resourceName string, metadata map[string]string,
) (RequestHandler, error) {
	allowedTypes := mw.GetIncludeTypes()
	if allowedTypes != nil {
		requestType := reflect.TypeOf(theRequest)
		if _, ok := allowedTypes[requestType.String()]; !ok {
			return nil, fmt.Errorf("request type %T not in includeTypes", theRequest)
		}
	}

	// Return the "base" request handler. User request handlers will be composed on top of this.
	return NewRequestHandler(
		HandlerProps{
			Metadata: metadata,
			Request:  theRequest,
			RequestName: requestName,
			ResourceType: resourceType,
			ResourceName: resourceName,
			Logger:   mw.GetLogger(),
		},
	), nil
}

func (mw *middleware) GetLogger() logger.MomentoLogger {
	return mw.logger
}

// GetRequestHandler returns a new RequestHandler that wraps the provided baseRequestHandler.
// Custom middlewares should implement this method to accept a base request handler and use
// it to compose a custom RequestHandler:
//
//	func (mw *myMiddleware) GetRequestHandler(baseHandler RequestHandler) (RequestHandler, error) {
//	  return NewMyMiddlewareRequestHandler(baseHandler, mw.myField), nil
//	}
func (mw *middleware) GetRequestHandler(_ RequestHandler) (RequestHandler, error) {
	return nil, fmt.Errorf("GetRequestHandler not implemented in middleware")
}

func (mw *middleware) GetIncludeTypes() map[string]bool {
	return mw.includeTypes
}

// NewMiddleware creates a new middleware with a logger and optional list of request types it should handle.
// If the IncludeTypes are omitted or empty, all request types will be processed. For example, to limit processing
// to only requests of type *momento.SetRequest and *momento.GetRequest, pass the following slice as the IncludeTypes:
//
//	[]interface{}{&momento.SetRequest{}, &momento.GetRequest{}}
//
// Custom middleware implementations can use this constructor to create and store the base middleware:
//
//	type myMiddleware struct {
//	  middleware.Middleware
//	  requestCount atomic.Int64
//	}
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
	id       uuid.UUID
	logger   logger.MomentoLogger
	request  interface{}
	requestName string
	resourceType internal.ClientType
	resourceName string
	metadata map[string]string
}

// RequestHandler is an interface that represents the capabilities of a middleware request handler.
// Custom request handlers will generally only need to implement the OnRequest and OnResponse methods.
type RequestHandler interface {
	GetId() uuid.UUID
	GetRequest() interface{}
	GetRequestName() string
	GetResourceType() internal.ClientType
	GetResourceName() string
	GetMetadata() map[string]string
	GetLogger() logger.MomentoLogger
	OnRequest()
	OnResponse(theResponse interface{}, err error) (interface{}, error)
}

type HandlerProps struct {
	Request  interface{}
	RequestName string
	ResourceType internal.ClientType
	ResourceName string
	Metadata map[string]string
	Logger   logger.MomentoLogger
}

func (rh *requestHandler) GetId() uuid.UUID {
	return rh.id
}

func (rh *requestHandler) GetRequest() interface{} {
	return rh.request
}

func (rh *requestHandler) GetRequestName() string {
	return rh.requestName
}

func (rh *requestHandler) GetMetadata() map[string]string {
	return rh.metadata
}

func (rh *requestHandler) GetLogger() logger.MomentoLogger {
	return rh.logger
}

func (rh *requestHandler) GetResourceType() internal.ClientType {
	return rh.resourceType
}

func (rh *requestHandler) GetResourceName() string {
	return rh.resourceName
}

// OnRequest is called before the request is made to the backend.
func (rh *requestHandler) OnRequest() {}

// OnResponse is called after the response is received from the backend. It is passed the response object, which can
// be cast to the appropriate response type for further inspection:
//
//	func (rh *myRequestHandler) OnResponse(theResponse interface{}) {
//	  switch r := theResponse.(type) {
//	  case *responses.ListPushFrontSuccess:
//	    fmt.Printf("pushed to front of list whose length is now %d\n", r.ListLength())
//	  case *responses.ListPushBackSuccess:
//	    fmt.Printf("pushed to back of list whose length is now %d\n", r.ListLength())
//	}
func (rh *requestHandler) OnResponse(theResponse interface{}, err error) (interface{}, error) {
	return theResponse, err
}

func NewRequestHandler(props HandlerProps) RequestHandler {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return &requestHandler{
		id:       id,
		logger:   props.Logger,
		request:  props.Request,
		requestName: props.RequestName,
		resourceType: props.ResourceType,
		resourceName: props.ResourceName,
		metadata: props.Metadata,
	}
}
