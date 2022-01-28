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
	client            pb.ScsClient
	defaultTtlSeconds uint32
}

func NewScsDataClient(dcr internalRequests.DataClientRequest) (*ScsDataClient, error) {
	newEndpoint := fmt.Sprint(dcr.Endpoint, CachePort)
	dataGrpcManagerRequest := internalRequests.DataGrpcManagerRequest{
		AuthToken: dcr.AuthToken,
		Endpoint:  newEndpoint,
	}
	cm, err := grpcmanagers.NewDataGrpcManager(dataGrpcManagerRequest)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsClient(cm.Conn)
	return &ScsDataClient{grpcManager: cm, client: client, defaultTtlSeconds: dcr.DefaultTtlSeconds}, nil
}

func (dc *ScsDataClient) Close() error {
	return dc.grpcManager.Close()
}

func (dc *ScsDataClient) Set(csr requests.CacheSetRequest) (*responses.SetCacheResponse, error) {
	if !utility.IsCacheNameValid(csr.CacheName) {
		return nil, scserrors.InvalidInputError("cache name cannot be empty")
	}
	var defaultTtlMils = dc.defaultTtlSeconds * 1000
	var ttlMils = csr.TtlSeconds * 1000
	byteKey, err := asBytes(csr.Key, "Unsupported type for key: ")
	if err != nil {
		return nil, err
	}
	byteValue, err := asBytes(csr.Value, "Unsupported type for value: ")
	if err != nil {
		return nil, err
	}
	var itemTtlMils uint32
	if csr.TtlSeconds == 0 {
		itemTtlMils = defaultTtlMils
	} else {
		itemTtlMils = ttlMils

	}
	request := pb.SetRequest{CacheKey: byteKey, CacheBody: byteValue, TtlMilliseconds: itemTtlMils}
	ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
	defer cancel()
	md := createNewMetadata(csr.CacheName)
	resp, err := dc.client.Set(metadata.NewOutgoingContext(ctx, md), &request)
	if err != nil {
		return nil, scserrors.GrpcErrorConverter(err)
	}
	newResp := responses.NewSetCacheResponse(resp, byteValue)
	return newResp, nil
}

func (dc *ScsDataClient) Get(cgr requests.CacheGetRequest) (*responses.GetCacheResponse, error) {
	if !utility.IsCacheNameValid(cgr.CacheName) {
		return nil, scserrors.InvalidInputError("cache name cannot be empty")
	}
	byteKey, err := asBytes(cgr.Key, "Unsupported type for key: ")
	if err != nil {
		return nil, err
	}
	request := pb.GetRequest{CacheKey: byteKey}
	ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
	defer cancel()
	md := createNewMetadata(cgr.CacheName)
	resp, err := dc.client.Get(metadata.NewOutgoingContext(ctx, md), &request)
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
