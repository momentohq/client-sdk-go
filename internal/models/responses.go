package models

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type ListCachesResponse struct {
	NextToken string
	Caches    []responses.CacheInfo
}

func NewListCacheResponse(resp *pb.XListCachesResponse) *ListCachesResponse {
	var caches []responses.CacheInfo
	for _, cache := range resp.Cache {
		caches = append(caches, NewCacheInfo(cache))
	}
	return &ListCachesResponse{NextToken: resp.NextToken, Caches: caches}
}

type CacheInfo struct {
	Name string
}

func NewCacheInfo(cache *pb.XCache) responses.CacheInfo {
	return responses.NewCacheInfo(cache.CacheName)
}

type TopicSubscribeResponse struct{}

type TopicPublishResponse struct{}
