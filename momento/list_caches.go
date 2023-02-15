package momento

////// Responses

type ListCachesResponse interface {
	isListCachesResponse()
}

// ListCachesResponse Output of the List caches operation.
type ListCachesSuccess struct {
	nextToken string
	caches    []CacheInfo
}

func (ListCachesSuccess) isListCachesResponse() {}

// NextToken Next Page Token returned by Simple Cache Service along with the list of caches.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListCachesRequest.
func (resp ListCachesSuccess) NextToken() string {
	return resp.nextToken
}

// Caches Returns all caches.
func (resp ListCachesSuccess) Caches() []CacheInfo {
	return resp.caches
}

// CacheInfo Information about a Cache.
type CacheInfo struct {
	name string
}

// Name Returns cache's name.
func (ci CacheInfo) Name() string {
	return ci.name
}

////// Request

type ListCachesRequest struct {
	// Token to continue paginating through the list. It's used to handle large paginated lists.
	NextToken string
}
