package momento

import (
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/resolver"
	"github.com/momentohq/client-sdk-go/internal/services"
)

const DeafultDataCtxTimeout = 5

type ScsClient struct {
	authToken         string
	defaultTtlSeconds uint32
	controlClient     *services.ScsControlClient
	dataClient        *services.ScsDataClient
}

func SimpleCacheClient(request *SimpleCacheClientRequest) (*ScsClient, momentoerrors.MomentoBaseError) {
	endpoints, momentoErr := resolver.Resolve(&models.ResolveRequest{
		AuthToken: request.AuthToken,
	})
	if momentoErr != nil {
		return nil, momentoErr
	}
	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		AuthToken: request.AuthToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, momentoerrors.ConvertError(err)
	}
	err = validateRequestTimeout(request.RequestTimeoutSeconds)
	if err != nil {
		return nil, err
	}
	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		AuthToken:         request.AuthToken,
		Endpoint:          endpoints.CacheEndpoint,
		DefaultTtlSeconds: request.DefaultTtlSeconds,
		DataCtxTimeout:    request.RequestTimeoutSeconds,
	})
	if err != nil {
		return nil, momentoerrors.ConvertError(err)
	}
	return &ScsClient{
		authToken:         request.AuthToken,
		defaultTtlSeconds: request.DefaultTtlSeconds,
		controlClient:     controlClient,
		dataClient:        dataClient,
	}, nil
}

func (client *ScsClient) CreateCache(request *CreateCacheRequest) momentoerrors.MomentoBaseError {
	return client.controlClient.CreateCache(&models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
}

func (client *ScsClient) DeleteCache(request *DeleteCacheRequest) momentoerrors.MomentoBaseError {
	return client.controlClient.DeleteCache(&models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
}

func (client *ScsClient) ListCaches(request *ListCachesRequest) (*ListCachesResponse, momentoerrors.MomentoBaseError) {
	rsp, err := client.controlClient.ListCaches(&models.ListCachesRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, err
	}
	return &ListCachesResponse{
		nextToken: rsp.NextToken,
		caches:    convertCacheInfo(rsp.Caches),
	}, nil
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

func (client *ScsClient) Set(request *CacheSetRequest) (*SetCacheResponse, momentoerrors.MomentoBaseError) {
	rsp, err := client.dataClient.Set(&models.CacheSetRequest{
		CacheName:  request.CacheName,
		Key:        request.Key,
		Value:      request.Value,
		TtlSeconds: request.TtlSeconds,
	})
	if err != nil {
		return nil, err
	}
	return &SetCacheResponse{
		value: rsp.Value,
	}, nil
}

func (client *ScsClient) Get(request *CacheGetRequest) (*GetCacheResponse, momentoerrors.MomentoBaseError) {
	rsp, err := client.dataClient.Get(&models.CacheGetRequest{
		CacheName: request.CacheName,
		Key:       request.Key,
	})
	if err != nil {
		return nil, err
	}
	return &GetCacheResponse{
		value:  rsp.Value,
		result: rsp.Result,
	}, nil
}

func (client *ScsClient) Close() error {
	ccErr := client.controlClient.Close()
	dErr := client.dataClient.Close()
	if ccErr != nil || dErr != nil {
		if ccErr != nil {
			return ccErr
		} else if dErr != nil {
			return dErr
		}
	}
	return nil
}

func validateRequestTimeout(requestTimeout *uint32) (err momentoerrors.MomentoBaseError) {
	if requestTimeout != nil && *requestTimeout == 0 {
		return momentoerrors.NewMomentoBaseError("InvalidArgumentError", "Request timeout must be greater than zero.")
	}
	return nil
}
