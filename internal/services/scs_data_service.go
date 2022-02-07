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

const CachePort = ":443"
const CacheCtxTimeout = 10 * time.Second

type ScsDataClient struct {
	grpcManager       *grpcmanagers.ScsDataGrpcManager
	dataClient        pb.ScsClient
	defaultTtlSeconds uint32
}

func NewScsDataClient(request *models.DataClientRequest) (*ScsDataClient, error) {
	dataManager, err := grpcmanagers.NewScsDataGrpcManager(&models.DataGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, CachePort),
	})
	if err != nil {
		return nil, err
	}
	return &ScsDataClient{
		grpcManager:       dataManager,
		dataClient:        pb.NewScsClient(dataManager.Conn),
		defaultTtlSeconds: request.DefaultTtlSeconds,
	}, nil
}

func (client *ScsDataClient) Close() error {
	return client.grpcManager.Close()
}

func (client *ScsDataClient) Set(request *models.CacheSetRequest) (*models.SetCacheResponse, error) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.InvalidInputError("cache name cannot be empty")
	}
	byteKey, err := asBytes(request.Key, "Unsupported type for key: ")
	if err != nil {
		return nil, err
	}
	byteValue, err := asBytes(request.Value, "Unsupported type for value: ")
	if err != nil {
		return nil, err
	}
	itemTtlMils := client.defaultTtlSeconds * 1000
	if request.TtlSeconds > 0 {
		itemTtlMils = request.TtlSeconds * 1000
	}
	ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
	defer cancel()
	resp, err := client.dataClient.Set(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSetRequest{
			CacheKey:        byteKey,
			CacheBody:       byteValue,
			TtlMilliseconds: itemTtlMils,
		},
	)
	if err != nil {
		return nil, momentoerrors.GrpcErrorConverter(err)
	}
	return models.NewSetCacheResponse(resp, byteValue), nil
}

func (client *ScsDataClient) Get(request *models.CacheGetRequest) (*models.GetCacheResponse, error) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.InvalidInputError("cache name cannot be empty")
	}
	byteKey, err := asBytes(request.Key, "Unsupported type for key: ")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
	defer cancel()
	resp, err := client.dataClient.Get(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XGetRequest{CacheKey: byteKey},
	)
	if err != nil {
		return nil, momentoerrors.GrpcErrorConverter(err)
	}
	newResp, err := models.NewGetCacheResponse(resp)
	if err != nil {
		return nil, err
	}
	return newResp, nil

}

func asBytes(data interface{}, message string) ([]byte, error) {
	switch data.(type) {
	case string:
		return []byte(reflect.ValueOf(data).String()), nil
	case []byte:
		return reflect.ValueOf(data).Bytes(), nil
	default:
		return nil, momentoerrors.InvalidInputError(fmt.Sprintf("%s %s", message, reflect.TypeOf(data).String()))
	}
}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
