package momento

import (
	"context"
	"strings"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/responses"
)

type PreviewStoreClient interface {
	// CreateStore creates a new store if it does not exist.
	CreateStore(ctx context.Context, request *CreateStoreRequest) (responses.CreateStoreResponse, momentoerrors.MomentoSvcErr)
	// DeleteStore deletes a store and all the items within it.
	DeleteStore(ctx context.Context, request *DeleteStoreRequest) (responses.DeleteStoreResponse, momentoerrors.MomentoSvcErr)
	// ListStores lists all the stores.
	ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, momentoerrors.MomentoSvcErr)
	// Get retrieves a value from a store.
	Get(ctx context.Context, request *StoreGetRequest) (responses.StoreGetResponse, momentoerrors.MomentoSvcErr)
	// Put sets a value in a store.
	Put(ctx context.Context, request *StorePutRequest) (responses.StorePutResponse, momentoerrors.MomentoSvcErr)
	// Delete removes a value from a store.
	Delete(ctx context.Context, request *StoreDeleteRequest) (responses.StoreDeleteResponse, momentoerrors.MomentoSvcErr)
	// Close closes the client.
	Close()
}

type defaultPreviewStoreClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	storeDataClient    *storeDataClient
	log                logger.MomentoLogger
}

// NewPreviewStoreClient creates a new PreviewStoreClient with the provided configuration and credential provider.
func NewPreviewStoreClient(storeConfiguration config.StoreConfiguration, credentialProvider auth.CredentialProvider) (PreviewStoreClient, error) {
	if storeConfiguration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must be greater than 0", nil)
	}
	client := &defaultPreviewStoreClient{
		credentialProvider: credentialProvider,
		log:                storeConfiguration.GetLoggerFactory().GetLogger("store-client"),
	}

	controlConfig := config.NewCacheConfiguration(&config.ConfigurationProps{
		TransportStrategy: storeConfiguration.GetTransportStrategy(),
		LoggerFactory:     storeConfiguration.GetLoggerFactory(),
	})
	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: credentialProvider,
		Configuration:      controlConfig,
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

func (c defaultPreviewStoreClient) CreateStore(ctx context.Context, request *CreateStoreRequest) (responses.CreateStoreResponse, momentoerrors.MomentoSvcErr) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	err := c.controlClient.CreateStore(ctx, &models.CreateStoreRequest{
		StoreName: request.StoreName,
	})
	if err != nil {
		if err.Code() == AlreadyExistsError {
			c.log.Info("Store with name '%s' already exists, skipping", request.StoreName)
			return &responses.CreateStoreAlreadyExists{}, nil
		}
		c.log.Warn("Error creating cache '%s': %s", request.StoreName, err.Message())
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &responses.CreateStoreSuccess{}, nil
}

func (c defaultPreviewStoreClient) DeleteStore(ctx context.Context, request *DeleteStoreRequest) (responses.DeleteStoreResponse, momentoerrors.MomentoSvcErr) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	err := c.controlClient.DeleteStore(ctx, &models.DeleteStoreRequest{
		StoreName: request.StoreName,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &responses.DeleteStoreSuccess{}, nil
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
	if err := isStoreNameValid(request.StoreName); err != nil {
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
	if err := isStoreNameValid(request.StoreName); err != nil {
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
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, err.(MomentoError)
	}
	// Doing a quick explicit null check instead of reimplementing the nest of interfaces involved
	// in the cache client's `prepareValue` null check.
	if request.Value == nil {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Value cannot be nil", nil)
	}

	resp, err := c.storeDataClient.set(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c defaultPreviewStoreClient) Close() {
	c.storeDataClient.Close()
	c.controlClient.Close()
}

func isStoreNameValid(storeName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(storeName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Store name cannot be empty", nil)
	}
	return nil
}
