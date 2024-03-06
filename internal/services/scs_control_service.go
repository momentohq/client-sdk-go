package services

import (
	"context"
	"google.golang.org/grpc/connectivity"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

const ControlCtxTimeout = 60 * time.Second

type ScsControlClient struct {
	grpcManager *grpcmanagers.ScsControlGrpcManager
	grpcClient  pb.ScsControlClient
}

func NewScsControlClient(request *models.ControlClientRequest) (*ScsControlClient, momentoerrors.MomentoSvcErr) {
	controlManager, err := grpcmanagers.NewScsControlGrpcManager(&models.ControlGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		RetryStrategy:      request.Configuration.GetRetryStrategy(),
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &ScsControlClient{grpcManager: controlManager, grpcClient: pb.NewScsControlClient(controlManager.Conn)}, nil
}

func (client *ScsControlClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsControlClient) CreateCache(ctx context.Context, request *models.CreateCacheRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.CreateCache(ctx, &pb.XCreateCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) DeleteCache(ctx context.Context, request *models.DeleteCacheRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	_, err := client.grpcClient.DeleteCache(ctx, &pb.XDeleteCacheRequest{CacheName: request.CacheName})
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsControlClient) ListCaches(ctx context.Context, request *models.ListCachesRequest) (*models.ListCachesResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, ControlCtxTimeout)
	defer cancel()
	resp, err := client.grpcClient.ListCaches(ctx, &pb.XListCachesRequest{NextToken: request.NextToken})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return models.NewListCacheResponse(resp), nil
}

func (client *ScsControlClient) Connect() error {
	timeout := ControlCtxTimeout

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := grpcmanagers.Connect(ctx, client.grpcManager.Conn)
	return err
}

func (client *ScsControlClient) IsConnected() bool {
	return client.grpcManager.Conn.GetState() == connectivity.Ready
}
