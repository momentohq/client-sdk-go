package scsmanager

import (
	"context"
	"fmt"
	"time"

	grpcmanagers "github.com/momentohq/client-sdk-go/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/protos"
	responses "github.com/momentohq/client-sdk-go/responses"
	utility "github.com/momentohq/client-sdk-go/utility"
)

const CONTROL_PORT = ":443"
const CONTROL_CTX_TIMEOUT = 10 * time.Second

type ScsControlClient struct {
	grpcManager grpcmanagers.ControlGrpcManager
	client      pb.ScsControlClient
}

func NewScsControlClient(authToken string, endPoint string) (*ScsControlClient, error) {
	newEndPoint := fmt.Sprint(endPoint, CONTROL_PORT)
	cm, err := grpcmanagers.NewControlGrpcManager(authToken, newEndPoint)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsControlClient(cm.Conn)
	return &ScsControlClient{grpcManager: cm, client: client}, nil
}

func (scc *ScsControlClient) Close() error {
	return scc.grpcManager.Close()
}

func (cc *ScsControlClient) CreateCache(cacheName string) error {
	if utility.IsCacheNameValid(cacheName) {
		request := pb.CreateCacheRequest{CacheName: cacheName}
		ctx, _ := context.WithTimeout(context.Background(), CONTROL_CTX_TIMEOUT)
		_, err := cc.client.CreateCache(ctx, &request)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cache name cannot be empty")
}

func (cc *ScsControlClient) DeleteCache(cacheName string) error {
	if utility.IsCacheNameValid(cacheName) {
		request := pb.DeleteCacheRequest{CacheName: cacheName}
		ctx, _ := context.WithTimeout(context.Background(), CONTROL_CTX_TIMEOUT)
		_, err := cc.client.DeleteCache(ctx, &request)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cache name cannot be empty")
}

func (cc *ScsControlClient) ListCaches(nextToken ...string) (*responses.ListCachesResponse, error) {
	defaultToken := ""
	if len(nextToken) != 0 {
		defaultToken = nextToken[0]
	}
	request := pb.ListCachesRequest{NextToken: defaultToken}
	ctx, _ := context.WithTimeout(context.Background(), CONTROL_CTX_TIMEOUT)
	resp, err := cc.client.ListCaches(ctx, &request)
	if err != nil {
		return nil, err
	}
	return responses.NewListCacheResponse(resp), nil
}
