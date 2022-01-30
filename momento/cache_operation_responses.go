package momento

import (
	"github.com/momentohq/client-sdk-go/internal/models"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

func NewListCacheResponse(resp *pb.ListCachesResponse) *ListCachesResponse {
	var caches = []CacheInfo{}
	for _, cache := range resp.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{nextToken: resp.NextToken, caches: caches}
}

func (resp *ListCachesResponse) NextToken() string {
	return resp.nextToken
}

func (resp *ListCachesResponse) Caches() []CacheInfo {
	return resp.caches
}

type CacheInfo struct {
	name string
}

func (ci CacheInfo) Name() string {
	return ci.name
}

func NewCacheInfo(cache *pb.Cache) CacheInfo {
	return CacheInfo{name: cache.CacheName}
}

const (
	HIT  string = "HIT"
	MISS string = "MISS"
)

type GetCacheResponse struct {
	value  []byte
	result string
}

func NewGetCacheResponse(resp *pb.GetResponse) (*GetCacheResponse, error) {
	var result string
	if resp.Result == pb.ECacheResult_Hit {
		result = HIT
	} else if resp.Result == pb.ECacheResult_Miss {
		result = MISS
	} else {
		return nil, models.ConvertEcacheResult(models.ConvertEcacheResultRequest{
			ECacheResult: resp.Result,
			Message:      resp.Message,
			OpName:       "GET",
		})
	}
	return &GetCacheResponse{value: resp.CacheBody, result: result}, nil
}

func (resp *GetCacheResponse) StringValue() string {
	if resp.result == HIT {
		return string(resp.value)
	}
	return ""
}

func (resp *GetCacheResponse) ByteValue() []byte {
	if resp.result == HIT {
		return resp.value
	}
	return nil
}

func (resp *GetCacheResponse) Result() string {
	return resp.result
}

type SetCacheResponse struct {
	value []byte
}

func NewSetCacheResponse(resp *pb.SetResponse, value []byte) *SetCacheResponse {
	return &SetCacheResponse{value: value}
}

func (resp *SetCacheResponse) StringValue() string {
	return string(resp.value)
}

func (resp *SetCacheResponse) ByteValue() []byte {
	return resp.value
}
