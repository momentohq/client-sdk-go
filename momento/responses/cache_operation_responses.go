package responses

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
	"github.com/momentohq/client-sdk-go/internal/utility"
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

func NewGetCacheResponse(gr *pb.GetResponse) (*GetCacheResponse, error) {
	var result string
	if gr.Result == pb.ECacheResult_Hit {
		result = HIT
	} else if gr.Result == pb.ECacheResult_Miss {
		result = MISS
	} else {
		resultRequest := internalRequests.ConvertEcacheResultRequest{
			ECacheResult: gr.Result,
			Message:      gr.Message,
			OpName:       "GET",
		}
		return nil, utility.ConvertEcacheResult(resultRequest)
	}
	return &GetCacheResponse{value: gr.CacheBody, result: result}, nil
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
	value []byte
}

func NewSetCacheResponse(sr *pb.SetResponse, value []byte) *SetCacheResponse {
	return &SetCacheResponse{value: value}
}

func (scr *SetCacheResponse) StringValue() string {
	return string(scr.value)
}

func (scr *SetCacheResponse) ByteValue() []byte {
	return scr.value
}
