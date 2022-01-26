package responses

import (
	pb "github.com/momentohq/client-sdk-go/protos"
	"github.com/momentohq/client-sdk-go/utility"
)

type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

func NewListCacheResponse(lcr *pb.ListCachesResponse) *ListCachesResponse {
	var caches = []CacheInfo{}
	for _, cache := range lcr.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{nextToken: lcr.NextToken, caches: caches}
}

func (lcr *ListCachesResponse) NextToken() string {
	return lcr.nextToken
}

func (lcr *ListCachesResponse) Caches() []CacheInfo {
	return lcr.caches
}

type CacheInfo struct {
	name string
}

func (ci CacheInfo) Name() string {
	return ci.name
}

func NewCacheInfo(ci *pb.Cache) CacheInfo {
	return CacheInfo{name: ci.CacheName}
}

const (
	HIT  string = "HIT"
	MISS string = "MISS"
)

type GetCacheResponse struct {
	value  []byte
	result string
}

func NewGetCacheResponse(gcr *pb.GetResponse) (*GetCacheResponse, error) {
	var result string
	if gcr.Result == pb.ECacheResult_Hit {
		result = HIT
	} else if gcr.Result == pb.ECacheResult_Miss {
		result = MISS
	} else {
		return nil, utility.ConvertEcacheResult(gcr.Result, gcr.Message, "GET")
	}
	return &GetCacheResponse{value: gcr.CacheBody, result: result}, nil
}

func (gcr *GetCacheResponse) StringValue() string {
	if gcr.result == HIT {
		return string(gcr.value)
	}
	return ""
}

func (gcr *GetCacheResponse) ByteValue() []byte {
	if gcr.result == HIT {
		return gcr.value
	}
	return nil
}

func (gcr *GetCacheResponse) Result() string {
	return gcr.result
}

type SetCacheResponse struct {
	result string
}

func NewSetCacheResponse(scr *pb.SetResponse) *SetCacheResponse {
	return &SetCacheResponse{result: scr.Result.String()}
}
