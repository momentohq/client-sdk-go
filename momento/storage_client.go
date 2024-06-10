package momento

import (
	"context"
	"strings"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/responses"
)

var storageDataClientCount uint64

type PreviewStorageClient interface {
	// CreateStore creates a new store if it does not exist.
	CreateStore(ctx context.Context, request *CreateStoreRequest) (responses.CreateStoreResponse, momentoerrors.MomentoSvcErr)
	// DeleteStore deletes a store and all the items within it.
	DeleteStore(ctx context.Context, request *DeleteStoreRequest) (responses.DeleteStoreResponse, momentoerrors.MomentoSvcErr)
	// ListStores lists all the stores.
	ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, momentoerrors.MomentoSvcErr)
	// Get retrieves a value from a store.
	Get(ctx context.Context, request *StorageGetRequest) (responses.StorageGetResponse, momentoerrors.MomentoSvcErr)
	// Set sets a value in a store.
	Set(ctx context.Context, request *StorageSetRequest) (responses.StorageSetResponse, momentoerrors.MomentoSvcErr)
	// Delete removes a value from a store.
	Delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, momentoerrors.MomentoSvcErr)
	// Close closes the client.
	Close()
}

type defaultPreviewStorageClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	storageDataClients []*storageDataClient
	log                logger.MomentoLogger
}

// NewPreviewStorageClient creates a new PreviewStorageClient with the provided configuration and credential provider.
func NewPreviewStorageClient(storageConfiguration config.StorageConfiguration, credentialProvider auth.CredentialProvider) (PreviewStorageClient, error) {
	if storageConfiguration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must be greater than 0", nil)
	}
	client := &defaultPreviewStorageClient{
		credentialProvider: credentialProvider,
		log:                storageConfiguration.GetLoggerFactory().GetLogger("store-client"),
	}

	controlConfig := config.NewCacheConfiguration(&config.ConfigurationProps{
		TransportStrategy: storageConfiguration.GetTransportStrategy(),
		LoggerFactory:     storageConfiguration.GetLoggerFactory(),
	})
	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: credentialProvider,
		Configuration:      controlConfig,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	numChannels := storageConfiguration.GetNumGrpcChannels()
	if numChannels < 1 {
		numChannels = 1
	}
	dataClients := make([]*storageDataClient, numChannels)

	for i := uint32(0); i < numChannels; i++ {
		storeDataClient, err := newStorageDataClient(&models.StorageDataClientRequest{
			CredentialProvider: credentialProvider,
			Configuration:      storageConfiguration,
		})
		if err != nil {
			return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
		}
		dataClients = append(dataClients, storeDataClient)
	}

	client.controlClient = controlClient
	client.storageDataClients = dataClients

	return client, nil
}

func (c defaultPreviewStorageClient) getNextStorageDataClient() *storageDataClient {
	nextClientIndex := atomic.AddUint64(&storageDataClientCount, 1)
	dataClient := c.storageDataClients[nextClientIndex%uint64(len(c.storageDataClients))]
	return dataClient
}

func (c defaultPreviewStorageClient) CreateStore(ctx context.Context, request *CreateStoreRequest) (responses.CreateStoreResponse, momentoerrors.MomentoSvcErr) {
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

func (c defaultPreviewStorageClient) DeleteStore(ctx context.Context, request *DeleteStoreRequest) (responses.DeleteStoreResponse, momentoerrors.MomentoSvcErr) {
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

func (c defaultPreviewStorageClient) ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, momentoerrors.MomentoSvcErr) {
	resp, err := c.controlClient.ListStores(ctx, &models.ListStoresRequest{})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return responses.NewListStoresSuccess(resp.NextToken, resp.Stores), nil
}

func (c defaultPreviewStorageClient) Delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, momentoerrors.MomentoSvcErr) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	response, err := c.getNextStorageDataClient().delete(ctx, request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return response, nil
}

func (c defaultPreviewStorageClient) Get(ctx context.Context, request *StorageGetRequest) (responses.StorageGetResponse, momentoerrors.MomentoSvcErr) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	resp, err := c.getNextStorageDataClient().get(ctx, request)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	return resp, nil
}

func (c defaultPreviewStorageClient) Set(ctx context.Context, request *StorageSetRequest) (responses.StorageSetResponse, momentoerrors.MomentoSvcErr) {
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

	resp, err := c.getNextStorageDataClient().set(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c defaultPreviewStorageClient) Close() {
	for _, dataClient := range c.storageDataClients {
		dataClient.Close()
	}
	c.controlClient.Close()
}

func isStoreNameValid(storeName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(storeName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Store name cannot be empty", nil)
	}
	return nil
}
