package responses

type CreateCacheResponse interface {
	isCreateCacheResponse()
}

type CreateCacheSuccess struct{}

func (CreateCacheSuccess) isCreateCacheResponse() {}

type CreateCacheAlreadyExists struct{}

func (CreateCacheAlreadyExists) isCreateCacheResponse() {}
