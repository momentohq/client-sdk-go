package requests

type SimpleCacheClientRequest struct {
	AuthToken         string
	DefaultTtlSeconds uint32
}

type CreateCacheRequest struct {
	CacheName string
}

type DeleteCacheRequest struct {
	CacheName string
}

type ListCachesRequest struct {
	NextToken string
}

type CacheSetRequest struct {
	CacheName  string
	Key        interface{}
	Value      interface{}
	TtlSeconds uint32
}

type CacheGetRequest struct {
	CacheName string
	Key       interface{}
}
