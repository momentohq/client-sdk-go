package cacheclient

import (
	"github.com/momentohq/client-sdk-go/resolver"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/scsmanagers"
)

type simpleCacheClient struct {
	authToken         string
	defaultTtlSeconds uint32
	controlClient     *scsmanagers.ScsControlClient
	dataClient        *scsmanagers.ScsDataClient
}

func SimpleCacheClient(authToken string, defaultTtlSeconds uint32) (*simpleCacheClient, error) {
	endPoints, err := resolver.Resolve(authToken)
	if err != nil {
		return nil, err
	}
	ctEndPoint := endPoints.ContorlEndPoint
	cEndPoint := endPoints.CacheEndPoint
	controlClient, ctErr := scsmanagers.NewScsControlClient(authToken, ctEndPoint)
	if ctErr != nil {
		return nil, ctErr
	}
	dataClient, cErr := scsmanagers.NewScsDataClient(authToken, cEndPoint, defaultTtlSeconds)
	if cErr != nil {
		return nil, cErr
	}
	return &simpleCacheClient{authToken: authToken, defaultTtlSeconds: defaultTtlSeconds, controlClient: controlClient, dataClient: dataClient}, nil
}

func (scc *simpleCacheClient) CreateCache(cacheName string) error {
	return scc.controlClient.CreateCache(cacheName)
}

func (scc *simpleCacheClient) DeleteCache(cacheName string) error {
	return scc.controlClient.DeleteCache(cacheName)
}

func (scc *simpleCacheClient) ListCaches(nextToken ...string) (*responses.ListCachesResponse, error) {
	return scc.controlClient.ListCaches(nextToken...)
}

func (scc *simpleCacheClient) Set(cacheName string, key interface{}, value interface{}, ttlSeconds ...uint32) (*responses.SetCacheResponse, error) {
	return scc.dataClient.Set(cacheName, key, value, ttlSeconds...)
}

func (scc *simpleCacheClient) Get(cacheName string, key interface{}) (*responses.GetCacheResponse, error) {
	return scc.dataClient.Get(cacheName, key)
}

func (scc *simpleCacheClient) Close() error {
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
