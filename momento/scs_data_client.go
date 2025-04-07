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
	r requester, requestMetadata map[string]string,
) ([]middleware.RequestHandler, requester, map[string]string, error) {
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
			return nil, nil, nil, err
		}

		// Call the request handler OnRequest method and then add the handler to list of handlers to
		// call OnResponse on when the response comes back.
		newReq, err := newHandler.OnRequest(r)
		if err != nil {
			return nil, nil, nil, err
		}
		if newReq != nil {
			if reflect.TypeOf(newReq) != reflect.TypeOf(r) {
				return nil, nil, nil, NewMomentoError(
					ClientSdkError,
					fmt.Sprintf("middleware request handler %T OnRequest returned an invalid request", newHandler),
					nil,
				)
			}
			r = newReq.(requester)
		}

		newMd := newHandler.OnMetadata(deepCopyMap(requestMetadata))
		if newMd != nil {
			requestMetadata = newMd
		}

		middlewareRequestHandlers = append(middlewareRequestHandlers, newHandler)
	}

	return middlewareRequestHandlers, r, requestMetadata, nil
}

// Iterate over the middleware request handlers in reverse order, giving them a chance to
// inspect the response and error results. Any error returned from the middleware OnResponse()
// method will be immediately returned as the actual error, skipping any outstanding response handlers.
func (client scsDataClient) applyMiddlewareResponseHandlers(
	middlewareRequestHandlers []middleware.RequestHandler,
	resp interface{},
	responseMetadata []metadata.MD,
) (interface{}, error) {
	for i := len(middlewareRequestHandlers) - 1; i >= 0; i-- {
		var requestHandlerError error
		rh := middlewareRequestHandlers[i]
		newResp, err := rh.OnResponse(resp)
		if err != nil {
			return nil, momentoerrors.ConvertSvcErr(requestHandlerError, responseMetadata...)
		}
		if newResp != nil {
			resp = newResp
		}
	}
	return resp, nil
}

func (client scsDataClient) makeRequest(ctx context.Context, r requester) (interface{}, error) {
	if _, err := prepareCacheName(r); err != nil {
		return nil, err
	}

	var middlewareRequestHandlers []middleware.RequestHandler
	requestMetadata := make(map[string]string)
	var err error
	middlewareRequestHandlers, r, requestMetadata, err = client.applyMiddlewareRequestHandlers(r, requestMetadata)
	if err != nil {
		return nil, err
	}

	req, err := r.initGrpcRequest(client)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()


	requestContext := internal.CreateCacheRequestContextFromMetadataMap(ctx, r.cacheName(), requestMetadata)
	resp, responseMetadata, err := r.makeGrpcRequest(req, requestContext, client)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err, responseMetadata...)
	}

	momentoResp, err := r.interpretGrpcResponse(resp)
	if err != nil {
		return nil, err
	}

	momentoResp, err = client.applyMiddlewareResponseHandlers(middlewareRequestHandlers, momentoResp, responseMetadata)
	if err != nil {
		return nil, err
	}

	return momentoResp, nil
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
