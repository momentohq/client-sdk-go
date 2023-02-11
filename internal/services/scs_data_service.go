package services

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc/metadata"
)

const defaultRequestTimeout = 5 * time.Second

type ScsDataClient struct {
	grpcManager    *grpcmanagers.DataGrpcManager
	grpcClient     pb.ScsClient
	defaultTtl     time.Duration
	requestTimeout time.Duration
	endpoint       string
}

func (c ScsDataClient) RequestTimeout() time.Duration { return c.requestTimeout }
func (c ScsDataClient) GrpcClient() pb.ScsClient      { return c.grpcClient }

func NewScsDataClient(request *models.DataClientRequest) (*ScsDataClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewUnaryDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})
	if err != nil {
		return nil, err
	}
	var timeout time.Duration
	if request.Configuration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.Configuration.GetClientSideTimeout()
	}
	return &ScsDataClient{
		grpcManager:    dataManager,
		grpcClient:     pb.NewScsClient(dataManager.Conn),
		defaultTtl:     request.DefaultTtl,
		requestTimeout: timeout,
		endpoint:       request.CredentialProvider.GetCacheEndpoint(),
	}, nil
}

func (client *ScsDataClient) Endpoint() string {
	return client.endpoint
}

func (client *ScsDataClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsDataClient) Delete(ctx context.Context, request *models.CacheDeleteRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.grpcClient.Delete(
		metadata.NewOutgoingContext(ctx, client.CreateNewMetadata(request.CacheName)),
		&pb.XDeleteRequest{CacheKey: request.Key},
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (ScsDataClient) CreateNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
