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

type CacheGetResponse interface {
	isCacheGetResponse()
}

// CacheGetMiss Miss Response to a cache Get api request.
type CacheGetMiss struct{}

func (_ CacheGetMiss) isCacheGetResponse() {}

// CacheGetHit Hit Response to a cache Get api request.
type CacheGetHit struct {
	Value []byte
}

func (_ CacheGetHit) isCacheGetResponse() {}

type TopicSubscribeResponse struct{}

type TopicPublishResponse struct{}
