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
	requestId := globalRequestId.Add(1)

	if _, err := prepareCacheName(r); err != nil {
		return err
	}

	if err := r.initGrpcRequest(client); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateCacheMetadata(ctx, r.cacheName())

	for _, mw := range client.middleware {
		mw.OnRequest(requestId, r, requestMetadata)
	}

	_, responseMetadata, err := r.makeGrpcRequest(requestMetadata, client)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err, responseMetadata...)
	}

	if err := r.interpretGrpcResponse(); err != nil {
		return err
	}

	for _, mw := range client.middleware {
		mw.OnResponse(requestId, r.getResponse())
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
