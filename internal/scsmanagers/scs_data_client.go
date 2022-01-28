package scsmanagers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
	"github.com/momentohq/client-sdk-go/internal/scserrors"
	"github.com/momentohq/client-sdk-go/internal/utility"
	"github.com/momentohq/client-sdk-go/momento/requests"
	"github.com/momentohq/client-sdk-go/momento/responses"

	"google.golang.org/grpc/metadata"
)

const CachePort = ":443"
const CacheCtxTimeout = 10 * time.Second

type ScsDataClient struct {
	grpcManager       *grpcmanagers.DataGrpcManager
	dataClient        pb.ScsClient
	defaultTtlSeconds uint32
}

func NewScsDataClient(request *internalRequests.DataClientRequest) (*ScsDataClient, error) {
	dataManager, err := grpcmanagers.NewDataGrpcManager(&internalRequests.DataGrpcManagerRequest{
		AuthToken: request.AuthToken,
		Endpoint:  fmt.Sprint(request.Endpoint, CachePort),
	})
	if err != nil {
		return nil, err
	}
	return &ScsDataClient{grpcManager: dataManager, dataClient: pb.NewScsClient(dataManager.Conn), defaultTtlSeconds: request.DefaultTtlSeconds}, nil
}

func (client *ScsDataClient) Close() error {
	return client.grpcManager.Close()
}

func (client *ScsDataClient) Set(request *requests.CacheSetRequest) (*responses.SetCacheResponse, error) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, scserrors.InvalidInputError("cache name cannot be empty")
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
	resp, err := client.dataClient.Set(metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)), &pb.SetRequest{CacheKey: byteKey, CacheBody: byteValue, TtlMilliseconds: itemTtlMils})
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	return responses.NewSetCacheResponse(resp, byteValue), nil
}

func (client *ScsDataClient) Get(request *requests.CacheGetRequest) (*responses.GetCacheResponse, error) {
	if !utility.IsCacheNameValid(request.CacheName) {
		return nil, scserrors.InvalidInputError("cache name cannot be empty")
	}
	byteKey, err := asBytes(request.Key, "Unsupported type for key: ")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
	defer cancel()
	resp, err := client.dataClient.Get(metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)), &pb.GetRequest{CacheKey: byteKey})
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	newResp, err := responses.NewGetCacheResponse(resp)
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
		return nil, scserrors.InvalidInputError(fmt.Sprintf("%s %s", message, reflect.TypeOf(data).String()))
	}
}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
