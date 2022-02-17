package momento

import (
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/resolver"
	"github.com/momentohq/client-sdk-go/internal/services"
)

const defaultRequestTimeout = uint32(5)

type ScsClient interface {
	CreateCache(request *CreateCacheRequest) error
	DeleteCache(request *DeleteCacheRequest) error
	ListCaches(request *ListCachesRequest) (*ListCachesResponse, error)

	Set(request *CacheSetRequest) (*SetCacheResponse, error)
	Get(request *CacheGetRequest) (*GetCacheResponse, error)

	Close()
}

type DefaultScsClient struct {
	authToken             string
	controlClient         *services.ScsControlClient
	dataClient            *services.ScsDataClient
	defaultTtlSeconds     uint32
	defaultRequestTimeout uint32
}

type Option func(*DefaultScsClient) MomentoError

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

func NewSimpleCacheClient(authToken string, defaultTtlSeconds uint32, opts ...Option) (ScsClient, error) {
	endpoints, err := resolver.Resolve(&models.ResolveRequest{
		AuthToken: authToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	client := &DefaultScsClient{
		authToken:         authToken,
		defaultTtlSeconds: defaultTtlSeconds,
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
		AuthToken: authToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		AuthToken:             authToken,
		Endpoint:              endpoints.CacheEndpoint,
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

func (c *DefaultScsClient) CreateCache(request *CreateCacheRequest) error {
	err := c.controlClient.CreateCache(&models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

func (c *DefaultScsClient) DeleteCache(request *DeleteCacheRequest) error {
	err := c.controlClient.DeleteCache(&models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

func (c *DefaultScsClient) ListCaches(request *ListCachesRequest) (*ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(&models.ListCachesRequest{
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

func (c *DefaultScsClient) Set(request *CacheSetRequest) (*SetCacheResponse, error) {
	ttlToUse := c.defaultTtlSeconds
	if request.TtlSeconds._ttl != nil {
		ttlToUse = *request.TtlSeconds._ttl
	}
	rsp, err := c.dataClient.Set(&models.CacheSetRequest{
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

func (c *DefaultScsClient) Get(request *CacheGetRequest) (*GetCacheResponse, error) {
	rsp, err := c.dataClient.Get(&models.CacheGetRequest{
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
