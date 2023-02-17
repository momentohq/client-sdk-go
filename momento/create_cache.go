package momento

///// Reponses

type CreateCacheResponse interface {
	isCreateCacheResponse()
}

type CreateCacheSuccess struct{}

func (CreateCacheSuccess) isCreateCacheResponse() {}

type CreateCacheAlreadyExists struct{}

func (CreateCacheAlreadyExists) isCreateCacheResponse() {}

///// Request

type CreateCacheRequest struct {
	// string used to create a cache.
	CacheName string
}
