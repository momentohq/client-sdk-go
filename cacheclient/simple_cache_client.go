package cacheclient

import (
	rl "github.com/momentohq/client-sdk-go/resolver"
	rs "github.com/momentohq/client-sdk-go/responses"
	cc "github.com/momentohq/client-sdk-go/scsmanagers"
)

type simpleCacheClient struct {
	authToken 			string
	defaultTtlSeconds 	uint32
	controlEndPoint 	string
	cacheEndPoint 		string
}

func SimpleCacheClient(authToken string, defaultTtlSeconds uint32) (*simpleCacheClient, error) {
	endPoints, err := rl.Resolve(authToken)
	if err != nil {
		return nil, err
	}
	return &simpleCacheClient{authToken: authToken, defaultTtlSeconds: defaultTtlSeconds, controlEndPoint: endPoints.ContorlEndPoint, cacheEndPoint: endPoints.CacheEndPoint}, nil
} 

func (scc *simpleCacheClient) CreateCache (cacheName string) error {
	endpoint := scc.controlEndPoint
	client, err := cc.NewScsControlClient(scc.authToken, endpoint)
	if err != nil {
		return err
	}
	return client.ScsCreateCache(cacheName)
}

func (scc *simpleCacheClient) DeleteCache(cacheName string) error {
	endpoint := scc.controlEndPoint
	client, err := cc.NewScsControlClient(scc.authToken, endpoint)
	if err != nil {
		return err
	}
	return client.ScsDeleteCache(cacheName)
}

func (scc *simpleCacheClient) ListCaches (nextToken ...string) (*rs.ListCachesResponse, error) {
	endpoint := scc.controlEndPoint
	client, err := cc.NewScsControlClient(scc.authToken, endpoint)
	if err != nil {
		return nil, err
	}
	return client.ScsListCaches(nextToken...)
}

func (scc *simpleCacheClient) Set (cacheName string, key interface{}, value interface{}, ttlSeconds ...uint32) (*rs.SetCacheResponse, error) {
	endpoint := scc.cacheEndPoint
	client, err := cc.NewScsDataClient(scc.authToken, endpoint, scc.defaultTtlSeconds)
	if err != nil {
		return nil, err
	}
	return client.ScsSet(cacheName, key, value, ttlSeconds...)
}

func (scc *simpleCacheClient) Get (cacheName string, key interface{}) (*rs.GetCacheResponse, error) {
	endpoint := scc.cacheEndPoint
	client, err := cc.NewScsDataClient(scc.authToken, endpoint, scc.defaultTtlSeconds)
	if err != nil {
		return nil, err
	}
	return client.ScsGet(cacheName, key)
}