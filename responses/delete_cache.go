package responses

type DeleteCacheResponse interface {
	isDeleteCacheResponse()
}

type DeleteCacheSuccess struct{}

func (DeleteCacheSuccess) isDeleteCacheResponse() {}
