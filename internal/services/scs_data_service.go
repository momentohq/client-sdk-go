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

func NewScsDataClient(request *models.DataClientRequest) (*ScsDataClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewScsDataGrpcManager(&models.DataGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, CachePort),
	})
	if err != nil {
		return nil, err
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

func (client *ScsDataClient) Set(request *models.CacheSetRequest) (*models.SetCacheResponse, momentoerrors.MomentoSvcErr) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty")
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
		itemTtlMils = request.TtlSeconds * 1000
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.dataCtxTimeout)
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
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return models.NewSetCacheResponse(resp, byteValue), nil
}

func (client *ScsDataClient) Get(request *models.CacheGetRequest) (*models.GetCacheResponse, momentoerrors.MomentoSvcErr) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty")
	}
	byteKey, momentoSvcErr := asBytes(request.Key, "Unsupported type for key: ")
	if momentoSvcErr != nil {
		return nil, momentoSvcErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.dataCtxTimeout)
	defer cancel()
	resp, err := client.dataClient.Get(
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

func asBytes(data interface{}, message string) ([]byte, momentoerrors.MomentoSvcErr) {
	switch data.(type) {
	case string:
		return []byte(reflect.ValueOf(data).String()), nil
	case []byte:
		return reflect.ValueOf(data).Bytes(), nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, fmt.Sprintf("%s %s", message, reflect.TypeOf(data).String()))
	}
}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
