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

type CacheGetResponse interface {
	isCacheGetResponse()
}

// CacheGetMiss Miss Response to a cache Get api request.
type CacheGetMiss struct{}

func (_ CacheGetMiss) isCacheGetResponse() {}

// CacheGetHit Hit Response to a cache Get api request.
type CacheGetHit struct {
	value []byte
}

func (_ CacheGetHit) isCacheGetResponse() {}

// ValueString Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp CacheGetHit) ValueString() string {
	return string(resp.value)
}

// ValueByte Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp CacheGetHit) ValueByte() []byte {
	return resp.value
}
