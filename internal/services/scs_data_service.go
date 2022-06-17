package services

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/utility"

	"google.golang.org/grpc/metadata"
)

const cachePort = ":443"
const defaultRequestTimeoutSeconds = 5

type ScsDataClient struct {
	grpcManager           *grpcmanagers.ScsDataGrpcManager
	grpcClient            pb.ScsClient
	defaultTtlSeconds     uint64
	requestTimeoutSeconds time.Duration
	endpoint              string
}

func NewScsDataClient(request *models.DataClientRequest) (*ScsDataClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewScsDataGrpcManager(&models.DataGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, cachePort),
	})
	if err != nil {
		return nil, err
	}
	var timeout time.Duration
	if request.RequestTimeoutSeconds < 1 {
		timeout = time.Duration(defaultRequestTimeoutSeconds) * time.Second
	} else {
		timeout = time.Duration(request.RequestTimeoutSeconds) * time.Second
	}
	return &ScsDataClient{
		grpcManager:           dataManager,
		grpcClient:            pb.NewScsClient(dataManager.Conn),
		defaultTtlSeconds:     uint64(request.DefaultTtlSeconds),
		requestTimeoutSeconds: timeout,
		endpoint:              request.Endpoint,
	}, nil
}

func (client *ScsDataClient) Endpoint() string {
	return client.endpoint
}

func (client *ScsDataClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}

func (client *ScsDataClient) Set(request *models.CacheSetRequest) (*models.SetCacheResponse, momentoerrors.MomentoSvcErr) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	byteKey, momentoSvcErr := asBytes(request.Key, "Unsupported type for key: ")
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	byteValue, momentoSvcErr := asBytes(request.Value, "Unsupported type for value: ")
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	itemTtlMils := client.defaultTtlSeconds * 1000
	if request.TtlSeconds > 0 {
		itemTtlMils = uint64(request.TtlSeconds * 1000)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.requestTimeoutSeconds)
	defer cancel()
	resp, err := client.grpcClient.Set(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSetRequest{
			CacheKey:        byteKey,
			CacheBody:       byteValue,
			TtlMilliseconds: itemTtlMils,
		},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return models.NewSetCacheResponse(resp, byteValue), nil
}

func (client *ScsDataClient) Get(request *models.CacheGetRequest) (*models.GetCacheResponse, momentoerrors.MomentoSvcErr) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	byteKey, momentoSvcErr := asBytes(request.Key, "Unsupported type for key: ")
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.requestTimeoutSeconds)
	defer cancel()
	resp, err := client.grpcClient.Get(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XGetRequest{CacheKey: byteKey},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	newResp, momentoSvcErr := models.NewGetCacheResponse(resp)
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	return newResp, nil

}

func (client *ScsDataClient) Delete(request *models.CacheDeleteRequest) (*models.DeleteCacheResponse, momentoerrors.MomentoSvcErr) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	byteKey, momentoSvcErr := asBytes(request.Key, "Unsupported type for key: ")
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.requestTimeoutSeconds)
	defer cancel()
	resp, err := client.grpcClient.Delete(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XDeleteRequest{CacheKey: byteKey},
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	newResp, momentoSvcErr := models.NewDeleteCacheResponse(resp)
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	return newResp, nil

}

func asBytes(data interface{}, message string) ([]byte, momentoerrors.MomentoSvcErr) {
	switch data.(type) {
	case string:
		return []byte(reflect.ValueOf(data).String()), nil
	case []byte:
		return reflect.ValueOf(data).Bytes(), nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, fmt.Sprintf("%s %s", message, reflect.TypeOf(data).String()), nil)
	}
}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
