package responses

import (
	pb "github.com/momentohq/client-sdk-go/protos"
)

type ListCachesResponse struct {
	NextToken	string
	Caches 		[]cacheInfo
}

func NewListCacheResponse(lcr *pb.ListCachesResponse) *ListCachesResponse {
	caches := []cacheInfo{}
	for _, cache := range lcr.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{NextToken: lcr.NextToken, Caches: caches}
}


type cacheInfo struct {
	Name	string
}

func NewCacheInfo(ci *pb.Cache) cacheInfo{
	return cacheInfo{Name: ci.CacheName}
}
