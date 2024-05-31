package momento

import (
	"context"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/retry"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/responses"
	"strings"
)

type PreviewStoreClient interface {
	CreateStore(ctx context.Context, request CreateStoreRequest) (responses.CreateStoreResponse, momentoerrors.MomentoSvcErr)
	DeleteStore(ctx context.Context, request DeleteStoreRequest) (responses.DeleteStoreResponse, momentoerrors.MomentoSvcErr)
	ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, momentoerrors.MomentoSvcErr)
	Get(ctx context.Context, request *StoreGetRequest) (responses.StoreGetResponse, momentoerrors.MomentoSvcErr)
	Put(ctx context.Context, request *StorePutRequest) (responses.StorePutResponse, momentoerrors.MomentoSvcErr)
	Delete(ctx context.Context, request *StoreDeleteRequest) (responses.StoreDeleteResponse, momentoerrors.MomentoSvcErr)

	Close()
}

type defaultPreviewStoreClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	storeDataClient    *storeDataClient
	log                logger.MomentoLogger
}

func NewPreviewStoreClient(storeConfiguration config.StoreConfiguration, credentialProvider auth.CredentialProvider) (PreviewStoreClient, error) {
	client := &defaultPreviewStoreClient{
		credentialProvider: credentialProvider,
		log:                storeConfiguration.GetLoggerFactory().GetLogger("store-client"),
	}

	// TODO: this is a bit of a mess
	cfg := config.NewCacheConfiguration(&config.ConfigurationProps{
		TransportStrategy: storeConfiguration.GetTransportStrategy(),
		LoggerFactory:     storeConfiguration.GetLoggerFactory(),
		RetryStrategy:     retry.NewNeverRetryStrategy(),
		NumGrpcChannels:   1,
		ReadConcern:       config.BALANCED,
	})
	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: credentialProvider,
		// TODO: ummmmm, this isn't great
		Configuration: cfg,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	storeDataClient, err := newStoreDataClient(&models.StoreDataClientRequest{
		CredentialProvider: credentialProvider,
		Configuration:      storeConfiguration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.controlClient = controlClient
	client.storeDataClient = storeDataClient

	return client, nil
}

func (c defaultPreviewStoreClient) Close() {
	c.storeDataClient.close()
}

func (c defaultPreviewStoreClient) CreateStore(ctx context.Context, request CreateStoreRequest) (responses.CreateStoreResponse, momentoerrors.MomentoSvcErr) {
	if err := isCacheNameValid(request.StoreName); err != nil {
		return nil, err
	}

	err := c.controlClient.CreateStore(ctx, &models.CreateStoreRequest{
		StoreName: request.StoreName,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return responses.CreateStoreSuccess{}, nil
}

func (c defaultPreviewStoreClient) DeleteStore(ctx context.Context, request DeleteStoreRequest) (responses.DeleteStoreResponse, momentoerrors.MomentoSvcErr) {
	if err := isCacheNameValid(request.StoreName); err != nil {
		return nil, err
	}

	err := c.controlClient.DeleteStore(ctx, &models.DeleteStoreRequest{
		StoreName: request.StoreName,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return responses.DeleteStoreSuccess{}, nil
}

func (c defaultPreviewStoreClient) ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, momentoerrors.MomentoSvcErr) {
	resp, err := c.controlClient.ListStores(ctx, &models.ListStoresRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return responses.NewListStoresSuccess(resp.NextToken, resp.Stores), nil
}

func (c defaultPreviewStoreClient) Delete(ctx context.Context, request *StoreDeleteRequest) (responses.StoreDeleteResponse, momentoerrors.MomentoSvcErr) {
	if err := isCacheNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	response, err := c.storeDataClient.delete(ctx, request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return response, nil
}

func (c defaultPreviewStoreClient) Get(ctx context.Context, request *StoreGetRequest) (responses.StoreGetResponse, momentoerrors.MomentoSvcErr) {
	if err := isCacheNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	resp, err := c.storeDataClient.get(ctx, request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	return resp, nil
}

func (c defaultPreviewStoreClient) Put(ctx context.Context, request *StorePutRequest) (responses.StorePutResponse, momentoerrors.MomentoSvcErr) {
	if err := isCacheNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	resp, err := c.storeDataClient.set(ctx, request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return resp, nil
}

func isStoreNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Store name cannot be empty", nil)
	}
	return nil
}
