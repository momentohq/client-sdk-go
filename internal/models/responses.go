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

func NewListCacheResponse(resp *pb.XListCachesResponse) *ListCachesResponse {
	var caches = []CacheInfo{}
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

type CacheResult string

const (
	HIT  CacheResult = "HIT"
	MISS CacheResult = "MISS"
)

type CacheGetResponse struct {
	Value  []byte
	Result CacheResult
}

type CacheSetRequest struct {
	CacheName  string
	Key        interface{}
	Value      interface{}
	TtlSeconds uint32
}

type CacheDeleteRequest struct {
	CacheName string
	Key       interface{}
}

type TopicSubscribeRequest struct {
	CacheName string
	TopicName string
}

type TopicSubscribeResponse struct{}

type TopicPublishRequest struct {
	CacheName string
	TopicName string
	Value     string // TODO think about string vs byte more
}

type TopicPublishResponse struct{}
