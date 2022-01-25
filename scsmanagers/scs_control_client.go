package scsmanager

import (
	"context"
	"fmt"
	"time"

	gm "github.com/momentohq/client-sdk-go/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/protos"
	rs "github.com/momentohq/client-sdk-go/responses"
	ut "github.com/momentohq/client-sdk-go/utility"
)

const CONTROL_PORT = ":443"

type ScsControlClient struct {
	GrpcManager		gm.ControlGrpcManager
	Client			pb.ScsControlClient
}

func NewScsControlClient(authToken string, endPoint string) (*ScsControlClient, error) {
	newEndPoint := fmt.Sprint(endPoint, CONTROL_PORT)
	cm, err := gm.NewControlGrpcManager(authToken, newEndPoint)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsControlClient(cm.Conn)
	return &ScsControlClient{GrpcManager: cm, Client: client}, nil
}

func (scc *ScsControlClient) Close() error {
	return scc.GrpcManager.Close()
}

func (cc *ScsControlClient) ScsCreateCache(cacheName string) error {
	if ut.IsCacheNameValid(cacheName) {
		request := pb.CreateCacheRequest{CacheName: cacheName}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := cc.Client.CreateCache(ctx, &request)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cache name cannot be empty")
}

func (cc *ScsControlClient) ScsDeleteCache(cacheName string) error {
	if ut.IsCacheNameValid(cacheName) {
		request := pb.DeleteCacheRequest{CacheName: cacheName}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := cc.Client.DeleteCache(ctx, &request)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cache name cannot be empty")
}


func (cc *ScsControlClient) ScsListCaches(nextToken ...string) (*rs.ListCachesResponse, error) {
	defaultToken := ""
	if len(nextToken) != 0 {
		defaultToken = nextToken[0]
	}
	request := pb.ListCachesRequest{NextToken: defaultToken}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cc.Client.ListCaches(ctx, &request)
	if err != nil {
		return nil, err
	}
	return rs.NewListCacheResponse(resp), nil
}
