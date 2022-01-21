package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/momentohq/client_sdk_go/protos"
)

type scsControlClient struct {
	GrpcManager		controlGrpcManager
	Client			pb.ScsControlClient
}

func NewScsControlClient(authToken string, endPoint string) (*scsControlClient, error) {
	cm, err := NewControlGrpcManager(authToken, endPoint)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsControlClient(cm.Conn)
	return &scsControlClient{GrpcManager: cm, Client: client}, nil
}

func (scc *scsControlClient) Close() error {
	return scc.GrpcManager.Close()
}

func (cc *scsControlClient) CreateCache(cacheName string) error {
	if isCacheNameValid(cacheName) {
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

func (cc *scsControlClient) DeleteCache(cacheName string) error {
	if isCacheNameValid(cacheName) {
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


func (cc *scsControlClient) ListCaches(nextToken ...string) (*listCachesResponse, error) {
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
	return NewListCacheResponse(resp), nil
}

func isCacheNameValid(cacheName string) bool {
	return len(strings.TrimSpace(cacheName)) != 0
}
