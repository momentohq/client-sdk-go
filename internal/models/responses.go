package models

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type CreateCacheRequest struct {
	CacheName string
}

type DeleteCacheRequest struct {
	CacheName string
}

type ListCachesRequest struct {
	NextToken string
}

type ListCachesResponse struct {
	NextToken string
	Caches    []CacheInfo
}

func NewListCacheResponse(resp *pb.ListCachesResponse) *ListCachesResponse {
	var caches = []CacheInfo{}
	for _, cache := range resp.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{NextToken: resp.NextToken, Caches: caches}
}

type CacheInfo struct {
	Name string
}

func NewCacheInfo(cache *pb.Cache) CacheInfo {
	return CacheInfo{Name: cache.CacheName}
}

const (
	HIT  string = "HIT"
	MISS string = "MISS"
)

type CacheGetRequest struct {
	CacheName string
	Key       interface{}
}

type GetCacheResponse struct {
	Value  []byte
	Result string
}

func NewGetCacheResponse(resp *pb.GetResponse) (*GetCacheResponse, error) {
	var result string
	if resp.Result == pb.ECacheResult_Hit {
		result = HIT
	} else if resp.Result == pb.ECacheResult_Miss {
		result = MISS
	} else {
		return nil, ConvertEcacheResult(ConvertEcacheResultRequest{
			ECacheResult: resp.Result,
			Message:      resp.Message,
			OpName:       "GET",
		})
	}
	return &GetCacheResponse{Value: resp.CacheBody, Result: result}, nil
}

type CacheSetRequest struct {
	CacheName  string
	Key        interface{}
	Value      interface{}
	TtlSeconds uint32
}

type SetCacheResponse struct {
	Value []byte
}

func NewSetCacheResponse(resp *pb.SetResponse, value []byte) *SetCacheResponse {
	return &SetCacheResponse{Value: value}
}
