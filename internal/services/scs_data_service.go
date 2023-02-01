package services

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc/metadata"
)

const defaultRequestTimeout = 5 * time.Second

type ScsDataClient struct {
	grpcManager       *grpcmanagers.DataGrpcManager
	grpcClient        pb.ScsClient
	defaultTtlSeconds uint64
	requestTimeout    time.Duration
	endpoint          string
}

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
		grpcManager:       dataManager,
		grpcClient:        pb.NewScsClient(dataManager.Conn),
		defaultTtlSeconds: uint64(request.DefaultTtlSeconds),
		requestTimeout:    timeout,
		endpoint:          request.CredentialProvider.GetCacheEndpoint(),
	}, nil
}

func (client *ScsDataClient) Endpoint() string {
	return client.endpoint
}

func (client *ScsDataClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsDataClient) Set(ctx context.Context, request *models.CacheSetRequest) momentoerrors.MomentoSvcErr {
	itemTtlMils := client.defaultTtlSeconds * 1000
	if request.TtlSeconds > 0 {
		itemTtlMils = uint64(request.TtlSeconds * 1000)
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.grpcClient.Set(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSetRequest{
			CacheKey:        request.Key,
			CacheBody:       request.Value,
			TtlMilliseconds: itemTtlMils,
		},
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsDataClient) Get(ctx context.Context, request *models.CacheGetRequest) (models.CacheGetResponse, momentoerrors.MomentoSvcErr) {
	// Execute request
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	resp, err := client.grpcClient.Get(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XGetRequest{CacheKey: request.Key},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	// Convert from grpc struct to internal struct
	if resp.Result == pb.ECacheResult_Hit {
		return &models.CacheGetHit{Value: resp.CacheBody}, nil
	} else if resp.Result == pb.ECacheResult_Miss {
		return &models.CacheGetMiss{}, nil
	} else {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InternalServerError,
			fmt.Sprintf(
				"CacheService returned an unexpected result: %v for operation: %s with message: %s",
				resp.Result, "GET", resp.Message,
			),
			nil,
		)
	}
}

func (client *ScsDataClient) Delete(ctx context.Context, request *models.CacheDeleteRequest) momentoerrors.MomentoSvcErr {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.grpcClient.Delete(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XDeleteRequest{CacheKey: request.Key},
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil

}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
