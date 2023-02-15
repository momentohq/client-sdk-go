package momento

///// Reponses

type CreateCacheResponse interface {
	isCreateCacheResponse()
}

type CreateCacheSuccess struct{}

func (CreateCacheSuccess) isCreateCacheResponse() {}

///// Request

type CreateCacheRequest struct {
	// string used to create a cache.
	CacheName string
}
