package cacheclient

import (
	rl "github.com/momentohq/client-sdk-go/resolver"
	rs "github.com/momentohq/client-sdk-go/responses"
	cc "github.com/momentohq/client-sdk-go/scsmanagers"
)

type simpleCacheClient struct {
	authToken         string
	defaultTtlSeconds uint32
	controlClient     *cc.ScsControlClient
	dataClient        *cc.ScsDataClient
}

func SimpleCacheClient(authToken string, defaultTtlSeconds uint32) (*simpleCacheClient, error) {
	endPoints, err := rl.Resolve(authToken)
	if err != nil {
		return nil, err
	}
	ctEndPoint := endPoints.ContorlEndPoint
	cEndPoint := endPoints.CacheEndPoint
	controlClient, ctErr := cc.NewScsControlClient(authToken, ctEndPoint)
	if ctErr != nil {
		return nil, ctErr
	}
	dataClient, cErr := cc.NewScsDataClient(authToken, cEndPoint, defaultTtlSeconds)
	if cErr != nil {
		return nil, cErr
	}
	return &simpleCacheClient{authToken: authToken, defaultTtlSeconds: defaultTtlSeconds, controlClient: controlClient, dataClient: dataClient}, nil
}

func (scc *simpleCacheClient) CreateCache(cacheName string) error {
	return scc.controlClient.ScsCreateCache(cacheName)
}

func (scc *simpleCacheClient) DeleteCache(cacheName string) error {
	return scc.controlClient.ScsDeleteCache(cacheName)
}

func (scc *simpleCacheClient) ListCaches(nextToken ...string) (*rs.ListCachesResponse, error) {
	return scc.controlClient.ScsListCaches(nextToken...)
}

func (scc *simpleCacheClient) Set(cacheName string, key interface{}, value interface{}, ttlSeconds ...uint32) (*rs.SetCacheResponse, error) {
	return scc.dataClient.ScsSet(cacheName, key, value, ttlSeconds...)
}

func (scc *simpleCacheClient) Get(cacheName string, key interface{}) (*rs.GetCacheResponse, error) {
	return scc.dataClient.ScsGet(cacheName, key)
}

func (scc *simpleCacheClient) Close() (error, error) {
	ccErr := scc.controlClient.Close()
	dErr := scc.dataClient.Close()
	if ccErr != nil || dErr != nil {
		if ccErr != nil {
			return ccErr, nil
		} else if dErr != nil {
			return nil, dErr
		} else {
			return ccErr, dErr
		}
	}
	return nil, nil
}
