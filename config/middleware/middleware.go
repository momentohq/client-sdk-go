package middleware

import (
	"context"
	"fmt"
	"reflect"

	"github.com/momentohq/client-sdk-go/config/retry"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/config/logger"
)

type middleware struct {
	logger       logger.MomentoLogger
	includeTypes map[string]bool
}

type Middleware interface {
	GetLogger() logger.MomentoLogger
	GetBaseRequestHandler(theRequest interface{}, requestName string, resourceName string, metadata map[string]string) (RequestHandler, error)
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
	theRequest interface{}, requestName string, resourceName string, metadata map[string]string,
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
			Metadata:     metadata,
			Request:      theRequest,
			RequestName:  requestName,
			ResourceName: resourceName,
			Logger:       mw.GetLogger(),
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
	id           uuid.UUID
	logger       logger.MomentoLogger
	request      interface{}
	requestName  string
	resourceName string
}

// RequestHandler is an interface that represents the capabilities of a middleware request handler.
// Custom request handlers will generally only need to implement the OnRequest and OnResponse methods.
type RequestHandler interface {
	GetId() uuid.UUID
	GetRequest() interface{}
	GetRequestName() string
	GetResourceName() string
	GetLogger() logger.MomentoLogger
	OnRequest(theRequest interface{}) (interface{}, error)
	OnMetadata(map[string]string) map[string]string
	OnResponse(theResponse interface{}, err error) (interface{}, error)
}

type HandlerProps struct {
	Request      interface{}
	RequestName  string
	ResourceName string
	Metadata     map[string]string
	Logger       logger.MomentoLogger
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

func (rh *requestHandler) GetLogger() logger.MomentoLogger {
	return rh.logger
}

func (rh *requestHandler) GetResourceName() string {
	return rh.resourceName
}

// OnRequest is called before the request is made to the backend. It can be used to modify the request object or
// return an error to halt the request. If the method is used to modify the request, the new request object returned here
// must be the same type as the original request object, and an error is returned if this is not the case. Returning nil
// from this method leaves the request unchanged. Returning an error halts the request and returns a ClientSdkError to
// the caller.
func (rh *requestHandler) OnRequest(_ interface{}) (interface{}, error) {
	return nil, nil
}

// OnMetadata is called before the request is made to the backend. It receives the current request
// metadata map, which it may modify and return. Returning nil from this method leaves the metadata unchanged.
func (rh *requestHandler) OnMetadata(map[string]string) map[string]string {
	return nil
}

// OnResponse is called after the response is received from the backend. It is passed the response object, which can
// be cast to the appropriate response type for further inspection. Returning nil from this method leaves the response
// unchanged. Returning an error will replace the current error, if any, with the new error, but response processing
// will continue so later middlewares can inspect and handle the error as appropriate. A response handler may return
// nil for the error to indicate that the current error should be nulled. If there is an error after
// all the middlewares have run, the error will be converted to a Momento service error and returned to the caller.
// Returning a new response object will replace the current response object. The new response must be the same
// type as the original response object, and an error is returned if this is not the case.
//
//		func (rh *myRequestHandler) OnResponse(_ interface{}, err error) (interface{}, error) {
//		  switch r := theResponse.(type) {
//		  case *responses.ListPushFrontSuccess:
//		    fmt.Printf("pushed to front of list whose length is now %d\n", r.ListLength())
//		  case *responses.ListPushBackSuccess:
//		    fmt.Printf("pushed to back of list whose length is now %d\n", r.ListLength())
//		  }
//	   return nil, err
//		}
func (rh *requestHandler) OnResponse(_ interface{}, err error) (interface{}, error) {
	return nil, err
}

func NewRequestHandler(props HandlerProps) RequestHandler {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return &requestHandler{
		id:           id,
		logger:       props.Logger,
		request:      props.Request,
		requestName:  props.RequestName,
		resourceName: props.ResourceName,
	}
}
