// Package momento represents API ScsClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.

package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/internal/utility"
)

// ScsClient wraps lower level cache control and data operations.
type ScsClient interface {
	CreateCache(ctx context.Context, request *CreateCacheRequest) error
	DeleteCache(ctx context.Context, request *DeleteCacheRequest) error
	ListCaches(ctx context.Context, request *ListCachesRequest) (*ListCachesResponse, error)

	CreateSigningKey(ctx context.Context, request *CreateSigningKeyRequest) (*CreateSigningKeyResponse, error)
	RevokeSigningKey(ctx context.Context, request *RevokeSigningKeyRequest) error
	ListSigningKeys(ctx context.Context, request *ListSigningKeysRequest) (*ListSigningKeysResponse, error)

	Set(ctx context.Context, request *CacheSetRequest) (*SetCacheResponse, error)
	Get(ctx context.Context, request *CacheGetRequest) (*GetCacheResponse, error)
	Delete(ctx context.Context, request *CacheDeleteRequest) error

	Close()
}

// DefaultScsClient represents all information needed for momento client to enable cache control and data operations.
type DefaultScsClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	dataClient         *services.ScsDataClient
	defaultTtlSeconds  uint32
}

type SimpleCacheClientProps struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTtlSeconds  uint32
}

// NewSimpleCacheClient returns a new ScsClient with provided authToken, DefaultTtlSeconds, and opts arguments.
func NewSimpleCacheClient(props *SimpleCacheClientProps) (ScsClient, error) {
	if props.Configuration.GetClientSideTimeoutMillis() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must not be 0", nil)
	}
	client := &DefaultScsClient{
		credentialProvider: props.CredentialProvider,
		defaultTtlSeconds:  props.DefaultTtlSeconds,
	}

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
		DefaultTtlSeconds:  props.DefaultTtlSeconds,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.dataClient = dataClient
	client.controlClient = controlClient

	return client, nil
}

// CreateCache Create a new cache in your Momento account.
// The following are possible errors that can be returned:
// InvalidArgumentError: If provided CacheName is empty.
// AlreadyExistsError: If cache with the given name already exists.
// ClientSdkError: For any SDK checks that fail.
func (c *DefaultScsClient) CreateCache(ctx context.Context, request *CreateCacheRequest) error {
	err := c.controlClient.CreateCache(ctx, &models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

// DeleteCache Deletes a cache and all the items within your Momento account.
// The following are possible errors that can be returned:
// InvalidArgumentError: If provided CacheName is empty.
// NotFoundError: If an attempt is made to delete a MomentoCache that doesn't exist.
// ClientSdkError: For any SDK checks that fail.
func (c *DefaultScsClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) error {
	err := c.controlClient.DeleteCache(ctx, &models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

// ListCaches Lists all caches in your Momento account.
// The following is a possible error that can be returned:
// AuthenticationError: If the provided Momento Auth Token is invalid.
func (c *DefaultScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (*ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(ctx, &models.ListCachesRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &ListCachesResponse{
		nextToken: rsp.NextToken,
		caches:    convertCacheInfo(rsp.Caches),
	}, nil
}

// CreateSigningKey Creates a Momento signing key in your Momento account
// The following are possible errors that can be returned:
// AuthenticationError: If the provided Momento Auth Token is invalid.
// InvalidArgumentError: If provided TtlMinutes is negative
// ClientSdkError: For any SDK checks that fail.
func (c *DefaultScsClient) CreateSigningKey(ctx context.Context, request *CreateSigningKeyRequest) (*CreateSigningKeyResponse, error) {
	rsp, err := c.controlClient.CreateSigningKey(ctx, c.dataClient.Endpoint(), &models.CreateSigningKeyRequest{
		TtlMinutes: request.TtlMinutes,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &CreateSigningKeyResponse{
		keyId:     rsp.KeyId,
		endpoint:  rsp.Endpoint,
		key:       rsp.Key,
		expiresAt: rsp.ExpiresAt,
	}, nil
}

// RevokeSigningKey Revokes a Momento signing key in your Momento account, all tokens signed by which will be invalid
// The following are possible errors that can be returned:
// AuthenticationError: If the provided Momento Auth Token is invalid.
// ClientSdkError: For any SDK checks that fail.
func (c *DefaultScsClient) RevokeSigningKey(ctx context.Context, request *RevokeSigningKeyRequest) error {
	err := c.controlClient.RevokeSigningKey(ctx, &models.RevokeSigningKeyRequest{
		KeyId: request.KeyId,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

// ListSigningKeys Lists all Momento signing keys in your Momento account
// The following are possible errors that can be returned:
// AuthenticationError: If the provided Momento Auth Token is invalid.
// ClientSdkError: For any SDK checks that fail.
func (c *DefaultScsClient) ListSigningKeys(ctx context.Context, request *ListSigningKeysRequest) (*ListSigningKeysResponse, error) {
	rsp, err := c.controlClient.ListSigningKeys(ctx, c.dataClient.Endpoint(), &models.ListSigningKeysRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &ListSigningKeysResponse{
		nextToken:   rsp.NextToken,
		signingKeys: convertSigningKey(rsp.SigningKeys),
	}, nil
}

// Set Stores an item in cache.
// The following are possible errors that can be returned:
// InvalidArgumentError: If provided CacheName is empty.
// NotFoundError: If the cache with the given name doesn't exist.
// InternalServerError: If server encountered an unknown error while trying to store the item.
func (c *DefaultScsClient) Set(ctx context.Context, request *CacheSetRequest) (*SetCacheResponse, error) {
	ttlToUse := c.defaultTtlSeconds
	if request.TtlSeconds._ttl != nil {
		ttlToUse = *request.TtlSeconds._ttl
	}
	err := utility.IsKeyValid(request.Key)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	err = utility.IsValueValid(request.Value)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	rsp, err := c.dataClient.Set(ctx, &models.CacheSetRequest{
		CacheName:  request.CacheName,
		Key:        request.Key,
		Value:      request.Value,
		TtlSeconds: ttlToUse,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &SetCacheResponse{
		value: rsp.Value,
	}, nil
}

// Get Retrieve an item from the cache.
// The following are possible errors that can be returned:
// InvalidArgumentError: If provided CacheName is empty.
// NotFoundError: If the cache with the given name doesn't exist.
// InternalServerError: If server encountered an unknown error while trying to store the item.
func (c *DefaultScsClient) Get(ctx context.Context, request *CacheGetRequest) (*GetCacheResponse, error) {
	err := utility.IsKeyValid(request.Key)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	rsp, err := c.dataClient.Get(ctx, &models.CacheGetRequest{
		CacheName: request.CacheName,
		Key:       request.Key,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &GetCacheResponse{
		value:  rsp.Value,
		result: rsp.Result,
	}, nil
}

// Delete an item from the cache.
// The following are possible errors that can be returned:
// InvalidArgumentError: If provided CacheName is empty.
// NotFoundError: If the cache with the given name doesn't exist.
// InternalServerError: If server encountered an unknown error while trying to delete the item.
func (c *DefaultScsClient) Delete(ctx context.Context, request *CacheDeleteRequest) error {
	err := utility.IsKeyValid(request.Key)
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	err = c.dataClient.Delete(ctx, &models.CacheDeleteRequest{
		CacheName: request.CacheName,
		Key:       request.Key,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

// Close Closes the client.
func (c *DefaultScsClient) Close() {
	defer c.controlClient.Close()
	defer c.dataClient.Close()
}

func convertMomentoSvcErrorToCustomerError(e momentoerrors.MomentoSvcErr) MomentoError {
	if e == nil {
		return nil
	}
	return NewMomentoError(e.Code(), e.Message(), e.OriginalErr())
}

func convertCacheInfo(i []models.CacheInfo) []CacheInfo {
	var convertedList []CacheInfo
	for _, c := range i {
		convertedList = append(convertedList, CacheInfo{
			name: c.Name,
		})
	}
	return convertedList
}

func convertSigningKey(sk []models.SigningKey) []SigningKey {
	var convertedList []SigningKey
	for _, s := range sk {
		convertedList = append(convertedList, SigningKey{
			keyId:     s.KeyId,
			endpoint:  s.Endpoint,
			expiresAt: s.ExpiresAt,
		})
	}
	return convertedList
}
