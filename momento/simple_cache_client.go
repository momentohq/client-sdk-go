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

func SimpleCacheClient(ccr requests.SimpleCacheClientRequest) (*ScsClient, error) {
	endpoints, err := resolver.Resolve(&internalRequests.ResolveRequest{
		AuthToken: ccr.AuthToken,
	})
	if err != nil {
		return nil, err
	}
	controlClient, err := scsmanagers.NewScsControlClient(&internalRequests.ControlClientRequest{
		AuthToken: ccr.AuthToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, err
	}
	dataClient, err := scsmanagers.NewScsDataClient(&internalRequests.DataClientRequest{
		AuthToken:         ccr.AuthToken,
		Endpoint:          endpoints.CacheEndpoint,
		DefaultTtlSeconds: ccr.DefaultTtlSeconds,
	})
	if err != nil {
		return nil, err
	}
	return &ScsClient{authToken: ccr.AuthToken, defaultTtlSeconds: ccr.DefaultTtlSeconds, controlClient: controlClient, dataClient: dataClient}, nil
}

func (scc *ScsClient) CreateCache(ccr *requests.CreateCacheRequest) error {
	return scc.controlClient.CreateCache(ccr)
}

func (scc *ScsClient) DeleteCache(dcr *requests.DeleteCacheRequest) error {
	return scc.controlClient.DeleteCache(dcr)
}

func (scc *ScsClient) ListCaches(lcr *requests.ListCachesRequest) (*responses.ListCachesResponse, error) {
	return scc.controlClient.ListCaches(lcr)
}

func (scc *ScsClient) Set(csr *requests.CacheSetRequest) (*responses.SetCacheResponse, error) {
	return scc.dataClient.Set(csr)
}

func (scc *ScsClient) Get(cgr *requests.CacheGetRequest) (*responses.GetCacheResponse, error) {
	return scc.dataClient.Get(cgr)
}

func (scc *ScsClient) Close() error {
	ccErr := scc.controlClient.Close()
	dErr := scc.dataClient.Close()
	if ccErr != nil || dErr != nil {
		if ccErr != nil {
			return ccErr
		} else if dErr != nil {
			return dErr
		}
	}
	return nil
}
