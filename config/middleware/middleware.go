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

func NewMiddleware(props Props) Middleware {
	var includeTypeMap map[string]bool
	if props.IncludeTypes != nil {
		includeTypeMap = make(map[string]bool)
		for _, t := range props.IncludeTypes {
			includeTypeMap[reflect.TypeOf(t).String()] = true
		}
	} else {
		includeTypeMap = nil
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
