package responses

// DeleteCacheResponse is a base response type for a delete cache request.
type DeleteCacheResponse interface {
	isDeleteCacheResponse()
}

// DeleteCacheSuccess indicates a successful delete cache request.
type DeleteCacheSuccess struct{}

func (DeleteCacheSuccess) isDeleteCacheResponse() {}
