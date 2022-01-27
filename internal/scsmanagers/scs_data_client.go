package scsmanagers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/responses"
	"github.com/momentohq/client-sdk-go/internal/scserrors"
	"github.com/momentohq/client-sdk-go/internal/utility"
	"google.golang.org/grpc/metadata"
)

const CachePort = ":443"
const CacheCtxTimeout = 10 * time.Second

type ScsDataClient struct {
	grpcManager       *grpcmanagers.DataGrpcManager
	client            pb.ScsClient
	defaultTtlSeconds uint32
}

func NewScsDataClient(authToken string, endPoint string, defaultTtlSeconds uint32) (*ScsDataClient, error) {
	newEndPoint := fmt.Sprint(endPoint, CachePort)
	cm, err := grpcmanagers.NewDataGrpcManager(authToken, newEndPoint)
	if err != nil {
		return nil, err
	}
	client := pb.NewScsClient(cm.Conn)
	er := isTtlValid(defaultTtlSeconds)
	if er != nil {
		cm.Conn.Close()
		return nil, er
	}
	return &ScsDataClient{grpcManager: cm, client: client, defaultTtlSeconds: defaultTtlSeconds}, nil
}

func (dc *ScsDataClient) Close() error {
	return dc.grpcManager.Close()
}

func (dc *ScsDataClient) Set(cacheName string, key interface{}, value interface{}, ttlSeconds ...uint32) (*responses.SetCacheResponse, error) {
	if utility.IsCacheNameValid(cacheName) {
		byteKey, errAsBytesKey := asBytes(key, "Unsupported type for key: ")
		if errAsBytesKey != nil {
			return nil, errAsBytesKey
		}
		byteValue, errAsBytesValue := asBytes(value, "Unsupported type for value: ")
		if errAsBytesValue != nil {
			return nil, errAsBytesValue
		}
		var itemTtlMils uint32
		if len(ttlSeconds) == 0 {
			itemTtlMils = dc.defaultTtlSeconds * 1000
		} else {
			err := isTtlValid(ttlSeconds[0])
			if err != nil {
				return nil, err
			} else {
				itemTtlMils = ttlSeconds[0] * 1000
			}
		}
		request := pb.SetRequest{CacheKey: byteKey, CacheBody: byteValue, TtlMilliseconds: itemTtlMils}
		ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
		defer cancel()
		md := createNewMetadata(cacheName)
		resp, errSet := dc.client.Set(metadata.NewOutgoingContext(ctx, md), &request)
		if errSet != nil {
			return nil, scserrors.GrpcErrorConverter(errSet)
		}
		newResp := responses.NewSetCacheResponse(resp)
		return newResp, nil
	}
	return nil, scserrors.InvalidInputError("cache name cannot be empty")
}

func (dc *ScsDataClient) Get(cacheName string, key interface{}) (*responses.GetCacheResponse, error) {
	if utility.IsCacheNameValid(cacheName) {
		byteKey, errAsBytes := asBytes(key, "Unsupported type for key: ")
		if errAsBytes != nil {
			return nil, errAsBytes
		}
		request := pb.GetRequest{CacheKey: byteKey}
		ctx, cancel := context.WithTimeout(context.Background(), CacheCtxTimeout)
		defer cancel()
		md := createNewMetadata(cacheName)
		resp, getErr := dc.client.Get(metadata.NewOutgoingContext(ctx, md), &request)
		if getErr != nil {
			return nil, scserrors.GrpcErrorConverter(getErr)
		}
		newResp, er := responses.NewGetCacheResponse(resp)
		if er != nil {
			return nil, er
		}
		return newResp, nil
	}
	return nil, scserrors.InvalidInputError("cache name cannot be empty")
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

func isTtlValid(ttlSeconds interface{}) error {
	if (reflect.TypeOf(ttlSeconds).String() != "uint32") || (reflect.ValueOf(ttlSeconds).Interface().(uint32) < uint32(0)) {
		return scserrors.InvalidInputError("ttl seconds must be a non-negative integer")
	}
	return nil
}

func createNewMetadata(cacheName string) metadata.MD {
	return metadata.Pairs("cache", cacheName)
}
