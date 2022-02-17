package momento

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
	TtlSeconds TimeToLive
}

func TTL(ttl uint32) TimeToLive {
	pInt := &ttl
	return TimeToLive{
		_ttl: pInt,
	}
}

type TimeToLive struct {
	_ttl *uint32
}

type CacheGetRequest struct {
	CacheName string
	Key       interface{}
}
