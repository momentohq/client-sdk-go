package responses

// ListCachesResponse is a base response type for a list caches request.
type ListCachesResponse interface {
	isListCachesResponse()
}

// ListCachesSuccess Output of the List caches operation.
type ListCachesSuccess struct {
	nextToken string
	caches    []CacheInfo
}

func (ListCachesSuccess) isListCachesResponse() {}

// NewListCachesSuccess returns a new ListCachesSuccess which indicates a successful list caches request.
func NewListCachesSuccess(nextToken string, caches []CacheInfo) *ListCachesSuccess {
	return &ListCachesSuccess{
		nextToken: nextToken,
		caches:    caches, //convertCacheInfo(caches),
	}
}

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

// NewCacheInfo returns new CacheInfo contains name.
func NewCacheInfo(name string) CacheInfo {
	return CacheInfo{name: name}
}
