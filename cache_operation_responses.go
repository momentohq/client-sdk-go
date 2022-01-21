package main

import (
	pb "github.com/momentohq/client_sdk_go/protos"
)

type listCachesResponse struct {
	NextToken	string
	Caches 		[]cacheInfo
}

func NewListCacheResponse(lcr *pb.ListCachesResponse) *listCachesResponse {
	caches := []cacheInfo{}
	for _, cache := range lcr.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &listCachesResponse{NextToken: lcr.NextToken, Caches: caches}
}


type cacheInfo struct {
	Name	string
}

func NewCacheInfo(ci *pb.Cache) cacheInfo{
	return cacheInfo{Name: ci.CacheName}
}
