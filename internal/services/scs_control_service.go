package services

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/utility"
)

const ControlPort = ":443"
const ControlCtxTimeout = 10 * time.Second

type ScsControlClient struct {
	grpcManager   *grpcmanagers.ScsControlGrpcManager
	controlClient pb.ScsControlClient
}

func NewScsControlClient(request *models.ControlClientRequest) (*ScsControlClient, error) {
	controlManager, err := grpcmanagers.NewScsControlGrpcManager(&models.ControlGrpcManagerRequest{
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

func (client *ScsControlClient) CreateCache(request *models.CreateCacheRequest) error {
	if !utility.IsCacheNameValid(request.CacheName) {
		return fmt.Errorf("cache name cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	_, err := client.controlClient.CreateCache(ctx, &pb.CreateCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return momentoerrors.GrpcErrorConverter(err)
	}
	return nil
}

func (client *ScsControlClient) DeleteCache(request *models.DeleteCacheRequest) error {
	if !utility.IsCacheNameValid(request.CacheName) {
		return fmt.Errorf("cache name cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	_, err := client.controlClient.DeleteCache(ctx, &pb.DeleteCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return momentoerrors.GrpcErrorConverter(err)
	}
	return nil
}

func (client *ScsControlClient) ListCaches(request *models.ListCachesRequest) (*models.ListCachesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ControlCtxTimeout)
	defer cancel()
	resp, err := client.controlClient.ListCaches(ctx, &pb.ListCachesRequest{NextToken: request.NextToken})
	if err != nil {
		return nil, momentoerrors.GrpcErrorConverter(err)
	}
	return models.NewListCacheResponse(resp), nil
}
