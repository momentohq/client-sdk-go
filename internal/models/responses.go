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

type ListStoresResponse struct {
	NextToken string
	Stores    []responses.StoreInfo
}

func NewListStoresResponse(resp *pb.XListStoresResponse) *ListStoresResponse {
	var stores []responses.StoreInfo
	for _, store := range resp.Store {
		stores = append(stores, NewStoreInfo(store))
	}
	return &ListStoresResponse{NextToken: resp.NextToken, Stores: stores}
}

type StoreInfo struct {
	Name string
}

func NewStoreInfo(store *pb.XStore) responses.StoreInfo {
	return responses.NewStoreInfo(store.StoreName)
}
