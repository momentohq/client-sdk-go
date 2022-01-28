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
	grpcManager   *grpcmanagers.ControlGrpcManager
	controlClient pb.ScsControlClient
}

func NewScsControlClient(request *internalRequests.ControlClientRequest) (*ScsControlClient, error) {
	controlManager, err := grpcmanagers.NewControlGrpcManager(&internalRequests.ControlGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, ControlPort),
	})
	if err != nil {
		return nil, err
	}
	return &ScsControlClient{grpcManager: controlManager, controlClient: pb.NewScsControlClient(controlManager.Conn)}, nil
}

func (client *ScsControlClient) Close() error {
	return client.grpcManager.Close()
}

func (client *ScsControlClient) CreateCache(request *requests.CreateCacheRequest) error {
	if !utility.IsCacheNameValid(request.CacheName) {
		return fmt.Errorf("cache name cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	_, err := client.controlClient.CreateCache(ctx, &pb.CreateCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return scserrors.GrpcErrorConverter(err)
	}
	return nil
}

func (client *ScsControlClient) DeleteCache(request *requests.DeleteCacheRequest) error {
	if !utility.IsCacheNameValid(request.CacheName) {
		return fmt.Errorf("cache name cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	_, err := client.controlClient.DeleteCache(ctx, &pb.DeleteCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return scserrors.GrpcErrorConverter(err)
	}
	return nil
}

func (client *ScsControlClient) ListCaches(request *requests.ListCachesRequest) (*responses.ListCachesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	resp, err := client.controlClient.ListCaches(ctx, &pb.ListCachesRequest{NextToken: request.NextToken})
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return responses.NewListCacheResponse(resp), nil
}
