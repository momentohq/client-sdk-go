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
	// Validate input
	if err := utility.IsKeyValid(request.Key); err != nil {
		return err
	}
	if err := utility.IsCacheNameValid(request.CacheName); err != nil {
		return err
	}
	byteKey, svcErr := utility.EncodeKey(request.Key)
	if svcErr != nil {
		return svcErr
	}
	byteValue, svcErr := utility.EncodeValue(request.Value)
	if svcErr != nil {
		return svcErr
	}

	itemTtlMils := client.defaultTtlSeconds * 1000
	if request.TtlSeconds > 0 {
		itemTtlMils = uint64(request.TtlSeconds * 1000)
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.grpcClient.Set(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSetRequest{
			CacheKey:        byteKey,
			CacheBody:       byteValue,
			TtlMilliseconds: itemTtlMils,
		},
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil
}

func (client *ScsDataClient) Get(ctx context.Context, request *models.CacheGetRequest) (models.CacheGetResponse, momentoerrors.MomentoSvcErr) {

	// Validate input
	if err := utility.IsKeyValid(request.Key); err != nil {
		return nil, err
	}
	if err := utility.IsCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	key, svcErr := utility.EncodeKey(request.Key)
	if svcErr != nil {
		return nil, svcErr
	}

	// Execute request
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	resp, err := client.grpcClient.Get(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XGetRequest{CacheKey: key},
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
	if err := utility.IsCacheNameValid(request.CacheName); err != nil {
		return err
	}
	byteKey, svcErr := utility.EncodeKey(request.Key)
	if svcErr != nil {
		return svcErr
	}

	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	_, err := client.grpcClient.Delete(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XDeleteRequest{CacheKey: byteKey},
	)
	if err != nil {
		return momentoerrors.ConvertSvcErr(err)
	}
	return nil

}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
