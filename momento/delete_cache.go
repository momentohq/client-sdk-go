package momento

type DeleteCacheRequest struct {
	// string cache name to delete.
	CacheName string
}

func (c DeleteCacheRequest) cacheName() string {
	return c.CacheName
}
