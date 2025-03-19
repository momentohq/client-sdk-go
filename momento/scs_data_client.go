package momento

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

const defaultRequestTimeout = 5 * time.Second
const defaultEagerConnectTimeout = 30 * time.Second

type scsDataClient struct {
	grpcManager         *grpcmanagers.DataGrpcManager
	grpcClient          pb.ScsClient
	defaultTtl          time.Duration
	requestTimeout      time.Duration
	endpoint            string
	eagerConnectTimeout time.Duration
	loggerFactory       logger.MomentoLoggerFactory
	middleware          []middleware.Middleware
}

func newScsDataClient(request *models.DataClientRequest, eagerConnectTimeout time.Duration) (*scsDataClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewUnaryDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		RetryStrategy:      request.Configuration.GetRetryStrategy(),
		ReadConcern:        request.Configuration.GetReadConcern(),
		GrpcConfiguration:  request.Configuration.GetTransportStrategy().GetGrpcConfig(),
		Middleware:         request.Configuration.GetMiddleware(),
	})
	if err != nil {
		return nil, err
	}
	var timeout time.Duration
	if request.Configuration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.Configuration.GetClientSideTimeout()
	}
	return &scsDataClient{
		grpcManager:         dataManager,
		grpcClient:          pb.NewScsClient(dataManager.Conn),
		defaultTtl:          request.DefaultTtl,
		requestTimeout:      timeout,
		endpoint:            request.CredentialProvider.GetCacheEndpoint(),
		eagerConnectTimeout: eagerConnectTimeout,
		loggerFactory:       request.Configuration.GetLoggerFactory(),
		middleware:          request.Configuration.GetMiddleware(),
	}, nil
}

func (client scsDataClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func deepCopyMap(original map[string]string) map[string]string {
	newMap := make(map[string]string, len(original))
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

func (client scsDataClient) applyMiddlewareRequestHandlers(
	r requester, req interface{}, requestMetadata map[string]string,
) ([]middleware.RequestHandler, interface{}, map[string]string, error) {
	middlewareRequestHandlers := make([]middleware.RequestHandler, 0, len(client.middleware))
	for _, mw := range client.middleware {
		// An error here means the middleware is configured to skip this type of request, so we
		// don't add it to the list of request handlers to call on response.
		newBaseHandler, err := mw.GetBaseRequestHandler(r, r.requestName(), r.cacheName())
		if err != nil {
			continue
		}

		// If the middleware is allowed to handle this request type, we use the base handler
		// to compose a more specific handler off of. An error here means something actually went wrong,
		// so we return it.
		newHandler, err := mw.GetRequestHandler(newBaseHandler)
		if err != nil {
			return nil, nil, nil, momentoerrors.NewMomentoSvcErr(momentoerrors.ClientSdkError, err.Error(), err)
		}

		// Call the request handler OnRequest method and then add the handler to list of handlers to
		// call OnResponse on when the response comes back.
		newReq, err := newHandler.OnRequest(req)
		if err != nil {
			return nil, nil, nil, momentoerrors.NewMomentoSvcErr(momentoerrors.ClientSdkError, err.Error(), err)
		}
		if newReq != nil {
			if reflect.TypeOf(newReq) != reflect.TypeOf(req) {
				return nil, nil, nil, NewMomentoError(
					ClientSdkError,
					fmt.Sprintf("middleware request handler %T OnRequest returned an invalid request", newHandler),
					nil,
				)
			}
			req = newReq
		}

		newMd := newHandler.OnMetadata(deepCopyMap(requestMetadata))
		if newMd != nil {
			requestMetadata = newMd
		}

		middlewareRequestHandlers = append(middlewareRequestHandlers, newHandler)
	}

	return middlewareRequestHandlers, req, requestMetadata, nil
}

// Iterate over the middleware request handlers in reverse order, giving them a chance to
// inspect the response and error results. Any error returned from the middleware OnResponse()
// method will be immediately returned as the actual error, skipping any outstanding response handlers.
// If none of the response handlers return an error, the original error (if any) will be returned after
// it is converted to a Momento service error.
func (client scsDataClient) applyMiddlewareResponseHandlers(
	r requester,
	middlewareRequestHandlers []middleware.RequestHandler,
	resp grpcResponse,
	requestError error,
	responseMetadata []metadata.MD,
) (grpcResponse, error) {
	var newResp interface{}
	for i := len(middlewareRequestHandlers) - 1; i >= 0; i-- {
		var requestHandlerError error
		rh := middlewareRequestHandlers[i]
		newResp, requestHandlerError = rh.OnResponse(resp, requestError)
		if requestHandlerError != nil {
			return nil, momentoerrors.ConvertSvcErr(requestHandlerError, responseMetadata...)
		}

		// The request handler returned a nil error, so we nil out the original error.
		requestError = nil

		if newResp != nil {
			err := r.validateResponseType(newResp.(grpcResponse))
			if err != nil {
				return nil, NewMomentoError(
					ClientSdkError,
					fmt.Sprintf("middleware request handler %T OnResponse returned an invalid response", rh),
					nil,
				)
			}
			resp = newResp.(grpcResponse)
		}
	}

	// If there were no middleware to process, and we were passed in an error,
	// we need to convert it to a Momento service error before returning it back.
	if requestError != nil {
		requestError = momentoerrors.ConvertSvcErr(requestError, responseMetadata...)
	}
	return resp, requestError
}

func (client scsDataClient) makeRequest(ctx context.Context, r requester) error {
	if _, err := prepareCacheName(r); err != nil {
		return err
	}

	req, err := r.initGrpcRequest(client)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	var middlewareRequestHandlers []middleware.RequestHandler
	requestMetadata := make(map[string]string)
	middlewareRequestHandlers, req, requestMetadata, err = client.applyMiddlewareRequestHandlers(r, req, requestMetadata)
	if err != nil {
		return err
	}

	requestContext := internal.CreateCacheRequestContextFromMetadataMap(ctx, r.cacheName(), requestMetadata)
	resp, responseMetadata, requestError := r.makeGrpcRequest(req, requestContext, client)

	resp, err = client.applyMiddlewareResponseHandlers(r, middlewareRequestHandlers, resp, requestError, responseMetadata)
	if err != nil {
		return err
	}

	if err := r.interpretGrpcResponse(resp); err != nil {
		return err
	}

	return nil
}

func (client scsDataClient) Connect() error {
	timeout := defaultEagerConnectTimeout
	if client.eagerConnectTimeout > 0 {
		timeout = client.eagerConnectTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := client.grpcManager.Connect(ctx)
	if err != nil {
		client.grpcManager.Close()
	}
	return err
}
