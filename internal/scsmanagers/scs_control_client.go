package scsmanagers

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
	"github.com/momentohq/client-sdk-go/internal/scserrors"
	"github.com/momentohq/client-sdk-go/internal/utility"
	"github.com/momentohq/client-sdk-go/momento/requests"
	"github.com/momentohq/client-sdk-go/momento/responses"
)

const ControlPort = ":443"
const ControlCtxTimeout = 10 * time.Second

type ScsControlClient struct {
	grpcManager *grpcmanagers.ControlGrpcManager
	client      pb.ScsControlClient
}

func NewScsControlClient(ccr internalRequests.ControlClientRequest) (*ScsControlClient, error) {
	newEndpoint := fmt.Sprint(ccr.Endpoint, ControlPort)
	controlGrpcManagerRequest := internalRequests.ControlGrpcManagerRequest{
		AuthToken: ccr.AuthToken,
		Endpoint:  newEndpoint,
	}
	cm, err := grpcmanagers.NewControlGrpcManager(controlGrpcManagerRequest)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsControlClient(cm.Conn)
	return &ScsControlClient{grpcManager: cm, client: client}, nil
}

func (cc *ScsControlClient) Close() error {
	return cc.grpcManager.Close()
}

func (cc *ScsControlClient) CreateCache(ccr requests.CreateCacheRequest) error {
	if utility.IsCacheNameValid(ccr.CacheName) {
		request := pb.CreateCacheRequest{CacheName: ccr.CacheName}
		ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
		defer cancel()
		_, err := cc.client.CreateCache(ctx, &request)
		if err != nil {
			return scserrors.GrpcErrorConverter(err)
		}
		return nil
	}
	return fmt.Errorf("cache name cannot be empty")
}

func (cc *ScsControlClient) DeleteCache(dcr requests.DeleteCacheRequest) error {
	if utility.IsCacheNameValid(dcr.CacheName) {
		request := pb.DeleteCacheRequest{CacheName: dcr.CacheName}
		ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
		defer cancel()
		_, err := cc.client.DeleteCache(ctx, &request)
		if err != nil {
			return scserrors.GrpcErrorConverter(err)
		}
		return nil
	}
	return fmt.Errorf("cache name cannot be empty")
}

func (cc *ScsControlClient) ListCaches(lcr requests.ListCachesRequest) (*responses.ListCachesResponse, error) {
	request := pb.ListCachesRequest{NextToken: lcr.NextToken}
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	resp, err := cc.client.ListCaches(ctx, &request)
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return responses.NewListCacheResponse(resp), nil
}
