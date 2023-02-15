package momento

//////// Responses

type DeleteCacheResponse interface {
	isDeleteCacheResponse()
}

type DeleteCacheSuccess struct{}

func (DeleteCacheSuccess) isDeleteCacheResponse() {}

//////// Request

type DeleteCacheRequest struct {
	// string cache name to delete.
	CacheName string
}
