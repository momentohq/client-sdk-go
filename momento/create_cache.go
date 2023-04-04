package momento

type CreateCacheRequest struct {
	// string used to create a cache.
	CacheName string
}

func (c CreateCacheRequest) cacheName() string {
	return c.CacheName
}
