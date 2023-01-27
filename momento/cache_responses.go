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

type cacheGetResponseTypes string

const (
	hit  cacheGetResponseTypes = "HIT"
	miss cacheGetResponseTypes = "MISS"
)

// CacheGetResponse Base type for possible responses a cache GET can return. Miss || Hit
type CacheGetResponse struct {
	responseType cacheGetResponseTypes
	value        []byte
}

// IsHit returns true if successfully fetched request item from cache otherwise returns false
func (r *CacheGetResponse) IsHit() bool {
	if r.responseType == hit {
		return true
	}
	return false
}

// AsHit returns CacheGetHitResponse pointer if successfully fetched request item otherwise returns nil
func (r *CacheGetResponse) AsHit() *CacheGetHitResponse {
	if r.IsHit() {
		return &CacheGetHitResponse{
			value: r.value,
		}
	}
	return nil
}

func (r *CacheGetResponse) IsMiss() bool {
	if r.responseType == miss {
		return true
	}
	return false
}
func (r *CacheGetResponse) AsMiss() *CacheGetMissResponse {
	if r.IsMiss() {
		return &CacheGetMissResponse{}
	}
	return nil
}

// CacheGetMissResponse Miss Response to a cache Get api request.
type CacheGetMissResponse struct{}

// CacheGetHitResponse Hit Response to a cache Get api request.
type CacheGetHitResponse struct {
	value []byte
}

// StringValue Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp *CacheGetHitResponse) StringValue() string {
	return string(resp.value)
}

// ByteValue Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp *CacheGetHitResponse) ByteValue() []byte {
	return resp.value
}
