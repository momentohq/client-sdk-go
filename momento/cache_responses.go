package momento

// ListCachesResponse Output of the List caches operation.
type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

// NextToken Next Page Token returned by Simple Cache Service along with the list of caches.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListCachesRequest.
func (resp *ListCachesResponse) NextToken() string {
	return resp.nextToken
}

// Caches Returns all caches.
func (resp *ListCachesResponse) Caches() []CacheInfo {
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
