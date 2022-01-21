package cacheclient

import (
	cc "github.com/momentohq/client-sdk-go/scsmanagers"
)

type simpleCacheClient struct {
	authToken 			string
	defaultTtlSeconds 	uint32
}

func SimpleCacheClient(authToken string, defaultTtlSeconds uint32) simpleCacheClient {
	return simpleCacheClient{authToken: authToken, defaultTtlSeconds: defaultTtlSeconds}
} 

func (scc simpleCacheClient) CreateCache (cacheName string) error {
	endpoint := "control.cell-alpha-dev.preprod.a.momentohq.com:443"
	client, err := cc.NewScsControlClient(scc.authToken, endpoint)
	if err != nil {
		return err
	}
	e := client.ScsCreateCache(cacheName)
	if e != nil {
		return e
	}
	return nil
}