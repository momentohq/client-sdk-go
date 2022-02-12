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
const DeafultDataCtxTimeout = 5 * time.Second

type ScsDataClient struct {
	grpcManager       *grpcmanagers.ScsDataGrpcManager
	dataClient        pb.ScsClient
	defaultTtlSeconds uint32
	dataCtxTimeout    time.Duration
}

func NewScsDataClient(request *models.DataClientRequest) (*ScsDataClient, momentoerrors.MomentoBaseError) {
	dataManager, err := grpcmanagers.NewScsDataGrpcManager(&models.DataGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, CachePort),
	})
	if err != nil {
		return nil, momentoerrors.ConvertError(err)
	}
	var timeout time.Duration
	if request.DataCtxTimeout == nil {
		timeout = time.Duration(DeafultDataCtxTimeout)
	} else {
		timeout = time.Duration(*request.DataCtxTimeout)
	}
	return &ScsDataClient{
		grpcManager:       dataManager,
		dataClient:        pb.NewScsClient(dataManager.Conn),
		defaultTtlSeconds: request.DefaultTtlSeconds,
		dataCtxTimeout:    timeout,
	}, nil
}

func (client *ScsDataClient) Close() error {
	return client.grpcManager.Close()
}

func (client *ScsDataClient) Set(request *models.CacheSetRequest) (*models.SetCacheResponse, momentoerrors.MomentoBaseError) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoBaseError("InvalidArgumentError", "Cache name cannot be empty")
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
	ctx, cancel := context.WithTimeout(context.Background(), client.dataCtxTimeout)
	defer cancel()
	resp, setErr := client.dataClient.Set(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XSetRequest{
			CacheKey:        byteKey,
			CacheBody:       byteValue,
			TtlMilliseconds: itemTtlMils,
		},
	)
	if err != nil {
		return nil, momentoerrors.ConvertError(setErr)
	}
	return models.NewSetCacheResponse(resp, byteValue), nil
}

func (client *ScsDataClient) Get(request *models.CacheGetRequest) (*models.GetCacheResponse, momentoerrors.MomentoBaseError) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoBaseError("InvalidArgumentError", "Cache name cannot be empty")
	}
	byteKey, err := asBytes(request.Key, "Unsupported type for key: ")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.dataCtxTimeout)
	defer cancel()
	resp, getErr := client.dataClient.Get(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XGetRequest{CacheKey: byteKey},
	)
	if getErr != nil {
		return nil, momentoerrors.ConvertError(getErr)
	}
	newResp, respErr := models.NewGetCacheResponse(resp)
	if respErr != nil {
		return nil, momentoerrors.ConvertError(respErr)
	}
	return newResp, nil

}

func asBytes(data interface{}, message string) ([]byte, momentoerrors.MomentoBaseError) {
	switch data.(type) {
	case string:
		return []byte(reflect.ValueOf(data).String()), nil
	case []byte:
		return reflect.ValueOf(data).Bytes(), nil
	default:
		return nil, momentoerrors.NewMomentoBaseError("InvalidArgumentError", fmt.Sprintf("%s %s", message, reflect.TypeOf(data).String()))
	}
}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
