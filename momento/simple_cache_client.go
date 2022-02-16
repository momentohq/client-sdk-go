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

	Close() error
}

type DefaultScsClient struct {
	authToken         string
	defaultTtlSeconds uint32
	controlClient     *services.ScsControlClient
	dataClient        *services.ScsDataClient
}

func NewSimpleCacheClient(request *SimpleCacheClientRequest) (ScsClient, error) {
	endpoints, err := resolver.Resolve(&models.ResolveRequest{
		AuthToken: request.AuthToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		AuthToken: request.AuthToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	requestTimeoutToUse := defaultRequestTimeout
	if request.RequestTimeoutSeconds > 1 {
		requestTimeoutToUse = request.RequestTimeoutSeconds
	}

	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		AuthToken:         request.AuthToken,
		Endpoint:          endpoints.CacheEndpoint,
		DefaultTtlSeconds: request.DefaultTtlSeconds,
		RequestTimeout:    requestTimeoutToUse,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}
	return &DefaultScsClient{
		authToken:         request.AuthToken,
		defaultTtlSeconds: request.DefaultTtlSeconds,
		controlClient:     controlClient,
		dataClient:        dataClient,
	}, nil
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
	rsp, err := c.dataClient.Set(&models.CacheSetRequest{
		CacheName:  request.CacheName,
		Key:        request.Key,
		Value:      request.Value,
		TtlSeconds: request.TtlSeconds,
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

func (c *DefaultScsClient) Close() error {
	err := c.controlClient.Close()
	err = c.dataClient.Close()
	return convertMomentoSvcErrorToCustomerError(err)
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
