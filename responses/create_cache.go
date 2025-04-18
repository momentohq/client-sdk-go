package responses

// CreateCacheResponse is the base response type for a create cache request.
type CreateCacheResponse interface {
	MomentoCacheResponse
	isCreateCacheResponse()
}

// CreateCacheSuccess indicates a successful create cache request.
type CreateCacheSuccess struct{}

func (CreateCacheSuccess) isCreateCacheResponse() {}

// CreateCacheAlreadyExists indicates that the cache already exists, so there was nothing to do.
type CreateCacheAlreadyExists struct{}

func (CreateCacheAlreadyExists) isCreateCacheResponse() {}
