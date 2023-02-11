package models

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListCachesResponse struct {
	NextToken string
	Caches    []CacheInfo
}

func NewListCacheResponse(resp *pb.XListCachesResponse) *ListCachesResponse {
	var caches []CacheInfo
	for _, cache := range resp.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{NextToken: resp.NextToken, Caches: caches}
}

type CacheInfo struct {
	Name string
}

func NewCacheInfo(cache *pb.XCache) CacheInfo {
	return CacheInfo{Name: cache.CacheName}
}

type TopicSubscribeResponse struct{}

type TopicPublishResponse struct{}
