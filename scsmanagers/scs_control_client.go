package scsmanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	gm "github.com/momentohq/client-sdk-go/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/protos"
	rs "github.com/momentohq/client-sdk-go/responses"
)

type scsControlClient struct {
	GrpcManager		gm.ControlGrpcManager
	Client			pb.ScsControlClient
}

func NewScsControlClient(authToken string, endPoint string) (*scsControlClient, error) {
	cm, err := gm.NewControlGrpcManager(authToken, endPoint)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsControlClient(cm.Conn)
	return &scsControlClient{GrpcManager: cm, Client: client}, nil
}

func (scc *scsControlClient) close() error {
	return scc.GrpcManager.Close()
}

func (cc *scsControlClient) ScsCreateCache(cacheName string) error {
	if isCacheNameValid(cacheName) {
		request := pb.CreateCacheRequest{CacheName: cacheName}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := cc.Client.CreateCache(ctx, &request)
		if err != nil {
			return err
		}
		return nil
	}
	defer cc.close()
	return fmt.Errorf("cache name cannot be empty")
}

func (cc *scsControlClient) ScsDeleteCache(cacheName string) error {
	if isCacheNameValid(cacheName) {
		request := pb.DeleteCacheRequest{CacheName: cacheName}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := cc.Client.DeleteCache(ctx, &request)
		if err != nil {
			return err
		}
		return nil
	}
	defer cc.close()
	return fmt.Errorf("cache name cannot be empty")
}


func (cc *scsControlClient) ScsListCaches(nextToken ...string) (*rs.ListCachesResponse, error) {
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
	defer cc.close()
	return rs.NewListCacheResponse(resp), nil
}

func isCacheNameValid(cacheName string) bool {
	return len(strings.TrimSpace(cacheName)) != 0
}
