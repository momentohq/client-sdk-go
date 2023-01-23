// Package momento represents API ScsClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.

package momento

import (
	"context"
	
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/internal/utility"
)

const defaultRequestTimeout = uint32(5)

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
	credentialProvider    auth.CredentialProvider
	controlClient         *services.ScsControlClient
	dataClient            *services.ScsDataClient
	defaultTtlSeconds     uint32
	defaultRequestTimeout uint32
}

// Option returns a function to configure various options for ScsClient.
type Option func(*DefaultScsClient) MomentoError

// WithRequestTimeout returns an Option that and set user-specified request timeout and  can be chained with other builder methods.
func WithRequestTimeout(requestTimeout uint32) Option {
	return func(c *DefaultScsClient) MomentoError {
		if requestTimeout == 0 {
			return NewMomentoError(
				momentoerrors.InvalidArgumentError,
				"request timeout must be greater than zero",
				nil,
			)
		}
		c.defaultRequestTimeout = requestTimeout
		return nil
	}
}

// NewSimpleCacheClient returns a new ScsClient with provided authToken, defaultTtlSeconds, and opts arguments.
func NewSimpleCacheClient(credentialProvider auth.CredentialProvider, defaultTtlSeconds uint32, opts ...Option) (ScsClient, error) {
	client := &DefaultScsClient{
		credentialProvider: credentialProvider,
		defaultTtlSeconds:  defaultTtlSeconds,
	}

	// Loop through all user passed options before building up internal clients
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *House as the argument
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	requestTimeoutToUse := defaultRequestTimeout
	if client.defaultRequestTimeout > 0 {
		requestTimeoutToUse = client.defaultRequestTimeout
	}

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: credentialProvider,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	// TODO: just pass authProvider for token and endpoint
	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		CredentialProvider:    credentialProvider,
		DefaultTtlSeconds:     defaultTtlSeconds,
		RequestTimeoutSeconds: requestTimeoutToUse,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.dataClient = dataClient
	client.controlClient = controlClient

	return client, nil
}

// Create a new cache in your Momento account.
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

// Deletes a cache and all the items within your Momento account.
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

// Lists all caches in your Momento account.
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

// Creates a Momento signing key in your Momento account
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

// Revokes a Momento signing key in your Momento account, all tokens signed by which will be invalid
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

// Lists all Momento signing keys in your Momento account
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

// Stores an item in cache.
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

// Retrieve an item from the cache.
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

// Closes the client.
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
