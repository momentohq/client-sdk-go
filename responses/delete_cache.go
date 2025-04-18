package responses

// DeleteCacheResponse is the base response type for a delete cache request.
type DeleteCacheResponse interface {
	MomentoCacheResponse
	isDeleteCacheResponse()
}

// DeleteCacheSuccess indicates a successful delete cache request.
type DeleteCacheSuccess struct{}

func (DeleteCacheSuccess) isDeleteCacheResponse() {}
