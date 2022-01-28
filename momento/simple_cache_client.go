package momento

import (
	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
	"github.com/momentohq/client-sdk-go/internal/resolver"
	"github.com/momentohq/client-sdk-go/internal/scsmanagers"
	"github.com/momentohq/client-sdk-go/momento/requests"
	"github.com/momentohq/client-sdk-go/momento/responses"
)

type ScsClient struct {
	authToken         string
	defaultTtlSeconds uint32
	controlClient     *scsmanagers.ScsControlClient
	dataClient        *scsmanagers.ScsDataClient
}

func SimpleCacheClient(request requests.SimpleCacheClientRequest) (*ScsClient, error) {
	endpoints, err := resolver.Resolve(&internalRequests.ResolveRequest{
		AuthToken: request.AuthToken,
	})
	if err != nil {
		return nil, err
	}
	controlClient, err := scsmanagers.NewScsControlClient(&internalRequests.ControlClientRequest{
		AuthToken: request.AuthToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, err
	}
	dataClient, err := scsmanagers.NewScsDataClient(&internalRequests.DataClientRequest{
		AuthToken:         request.AuthToken,
		Endpoint:          endpoints.CacheEndpoint,
		DefaultTtlSeconds: request.DefaultTtlSeconds,
	})
	if err != nil {
		return nil, err
	}
	return &ScsClient{authToken: request.AuthToken, defaultTtlSeconds: request.DefaultTtlSeconds, controlClient: controlClient, dataClient: dataClient}, nil
}

func (client *ScsClient) CreateCache(request *requests.CreateCacheRequest) error {
	return client.controlClient.CreateCache(request)
}

func (client *ScsClient) DeleteCache(request *requests.DeleteCacheRequest) error {
	return client.controlClient.DeleteCache(request)
}

func (client *ScsClient) ListCaches(request *requests.ListCachesRequest) (*responses.ListCachesResponse, error) {
	return client.controlClient.ListCaches(request)
}

func (client *ScsClient) Set(request *requests.CacheSetRequest) (*responses.SetCacheResponse, error) {
	return client.dataClient.Set(request)
}

func (client *ScsClient) Get(request *requests.CacheGetRequest) (*responses.GetCacheResponse, error) {
	return client.dataClient.Get(request)
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
