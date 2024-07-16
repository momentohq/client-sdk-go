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

// PreviewStorageClient PREVIEW Momento Storage Client
//
// WARNING: the API for this client is not yet stable and may change without notice.
// Please contact Momento if you would like to try this preview.
type PreviewStorageClient interface {
	Logger() logger.MomentoLogger

	// CreateStore creates a new store if it does not exist.
	CreateStore(ctx context.Context, request *CreateStoreRequest) (responses.CreateStoreResponse, error)
	// DeleteStore deletes a store and all the items within it.
	DeleteStore(ctx context.Context, request *DeleteStoreRequest) (responses.DeleteStoreResponse, error)
	// ListStores lists all the stores.
	ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, error)
	// Get retrieves a value from a store.
	Get(ctx context.Context, request *StorageGetRequest) (*responses.StorageGetResponse, error)
	// Put sets a value in a store.
	Put(ctx context.Context, request *StoragePutRequest) (responses.StoragePutResponse, error)
	// Delete removes a value from a store.
	Delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, error)
	// Close closes the client.
	Close()
}

type defaultPreviewStorageClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	storageDataClients []*storageDataClient
	logger             logger.MomentoLogger
}

// NewPreviewStorageClient creates a new PreviewStorageClient with the provided configuration and credential provider.
//
// WARNING: the API for this client is not yet stable and may change without notice.
// Please contact Momento if you would like to try this preview.
func NewPreviewStorageClient(storageConfiguration config.StorageConfiguration, credentialProvider auth.CredentialProvider) (PreviewStorageClient, error) {
	if storageConfiguration.GetClientSideTimeout() < 1 {
		return nil, NewMomentoError(momentoerrors.InvalidArgumentError, "request timeout must be greater than 0", nil)
	}
	client := &defaultPreviewStorageClient{
		credentialProvider: credentialProvider,
		logger:             storageConfiguration.GetLoggerFactory().GetLogger("store-client"),
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
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	numChannels := storageConfiguration.GetNumGrpcChannels()
	if numChannels < 1 {
		numChannels = 1
	}
	dataClients := make([]*storageDataClient, 0)

	for i := uint32(0); i < numChannels; i++ {
		storeDataClient, err := newStorageDataClient(&models.StorageDataClientRequest{
			CredentialProvider: credentialProvider,
			Configuration:      storageConfiguration,
		})
		if err != nil {
			return nil, convertMomentoSvcErrorToCustomerError(err)
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

func (c defaultPreviewStorageClient) Logger() logger.MomentoLogger {
	return c.logger
}

func (c defaultPreviewStorageClient) CreateStore(ctx context.Context, request *CreateStoreRequest) (responses.CreateStoreResponse, error) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	err := c.controlClient.CreateStore(ctx, &models.CreateStoreRequest{
		StoreName: request.StoreName,
	})
	if err != nil {
		if err.Code() == AlreadyExistsError {
			c.logger.Info("Store with name '%s' already exists, skipping", request.StoreName)
			return &responses.CreateStoreAlreadyExists{}, nil
		}
		c.logger.Warn("Error creating cache '%s': %s", request.StoreName, err.Message())
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &responses.CreateStoreSuccess{}, nil
}

func (c defaultPreviewStorageClient) DeleteStore(ctx context.Context, request *DeleteStoreRequest) (responses.DeleteStoreResponse, error) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	// TODO: figure out how to retrieve metadata from control plane calls
	// var header, trailer metadata.MD
	err := c.controlClient.DeleteStore(
		ctx,
		&models.DeleteStoreRequest{
			StoreName: request.StoreName,
		})
	if err != nil {
		// TODO: remove this once delete store accepts the metadata CallOptions and
		// returns metadata that can be used to differentiate between the not found errors.
		// Currently the default is CacheNotFoundError
		if err.Code() == CacheNotFoundError {
			return nil, NewMomentoError(StoreNotFoundError, "Store not found", nil)
		}
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &responses.DeleteStoreSuccess{}, nil
}

func (c defaultPreviewStorageClient) ListStores(ctx context.Context, request *ListStoresRequest) (responses.ListStoresResponse, error) {
	resp, err := c.controlClient.ListStores(ctx, &models.ListStoresRequest{})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return responses.NewListStoresSuccess(resp.NextToken, resp.Stores), nil
}

func (c defaultPreviewStorageClient) Delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, error) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, err
	}

	response, err := c.getNextStorageDataClient().delete(ctx, request)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return response, nil
}

func (c defaultPreviewStorageClient) Get(ctx context.Context, request *StorageGetRequest) (*responses.StorageGetResponse, error) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, err
	}

	resp, err := c.getNextStorageDataClient().get(ctx, request)
	// Item not found errors are being converted to NotFound responses in the data client
	if err != nil {
		return resp, convertMomentoSvcErrorToCustomerError(err)
	}

	return resp, nil
}

func (c defaultPreviewStorageClient) Put(ctx context.Context, request *StoragePutRequest) (responses.StoragePutResponse, error) {
	if err := isStoreNameValid(request.StoreName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.Key, "Key"); err != nil {
		return nil, err
	}
	// Doing a quick explicit null check instead of reimplementing the nest of interfaces involved
	// in the cache client's `prepareValue` null check.
	if request.Value == nil {
		return nil, NewMomentoError(momentoerrors.InvalidArgumentError, "Value cannot be nil", nil)
	}

	resp, err := c.getNextStorageDataClient().put(ctx, request)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return resp, nil
}

func (c defaultPreviewStorageClient) Close() {
	for _, dataClient := range c.storageDataClients {
		dataClient.Close()
	}
	c.controlClient.Close()
}

func isStoreNameValid(storeName string) error {
	if len(strings.TrimSpace(storeName)) < 1 {
		return NewMomentoError(momentoerrors.InvalidArgumentError, "Store name cannot be empty", nil)
	}
	return nil
}
