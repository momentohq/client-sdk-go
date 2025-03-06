package momento

import (
	"context"
	"sync/atomic"
	"time"

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

var globalRequestId atomic.Uint64

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

func (client scsDataClient) makeRequest(ctx context.Context, r requester) error {
	if _, err := prepareCacheName(r); err != nil {
		return err
	}

	if err := r.initGrpcRequest(client); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateCacheMetadata(ctx, r.cacheName())

	middlewareRequestHandlers := make([]middleware.RequestHandler, 0, len(client.middleware))
	for _, mw := range client.middleware {
		// An error here means the middleware is configured to skip this type of request, so we
		// don't add it to the list of request handlers to call on response.
		newBaseHandler, err := mw.GetBaseRequestHandler(r, requestMetadata)
		if err != nil {
			continue
		}

		// If the middleware is allowed to handle this request type, we use the base handler
		// to compose a more specific handler off of. An error here means something actually went wrong,
		// so we return it.
		newHandler, err := mw.GetRequestHandler(newBaseHandler)
		if err != nil {
			return momentoerrors.NewMomentoSvcErr(momentoerrors.ClientSdkError, err.Error(), err)
		}

		// Call the request handler OnRequest method and then add the handler to list of handlers to
		// call OnResponse on when the response comes back.
		newHandler.OnRequest()

		// Give the middleware a chance to modify the request metadata. If a middleware doesn't implement
		// GetMetadata, the base response handler will return the original metadata.
		requestMetadata = newHandler.GetMetadata()
		middlewareRequestHandlers = append(middlewareRequestHandlers, newHandler)
	}

	_, responseMetadata, err := r.makeGrpcRequest(requestMetadata, client)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err, responseMetadata...)
	}

	if err := r.interpretGrpcResponse(); err != nil {
		return err
	}

	for _, rh := range middlewareRequestHandlers {
		rh.OnResponse(r)
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
