package momento

// Output of the List caches operation.
type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

// Next Page Token returned by Simple Cache Service along with the list of caches.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListCachesRequest.
func (resp *ListCachesResponse) NextToken() string {
	return resp.nextToken
}

// Returns all caches.
func (resp *ListCachesResponse) Caches() []CacheInfo {
	return resp.caches
}

// Information about the Cache.
type CacheInfo struct {
	name string
}

// Returns cache's name.
func (ci CacheInfo) Name() string {
	return ci.name
}

const (
	// Represents cache hit.
	HIT string = "HIT"
	// Represents cache miss.
	MISS string = "MISS"
)

// Initializes GetCacheResponse to handle gRPC get response.
type GetCacheResponse struct {
	value  []byte
	result string
}

// Returns value stored in cache as string if there was Hit. Returns an empty string otherwise.
func (resp *GetCacheResponse) StringValue() string {
	if resp.result == HIT {
		return string(resp.value)
	}
	return ""
}

// Returns value stored in cache as bytes if there was Hit. Returns nil otherwise.
func (resp *GetCacheResponse) ByteValue() []byte {
	if resp.result == HIT {
		return resp.value
	}
	return nil
}

// Returns get operation result such as HIT or MISS.
func (resp *GetCacheResponse) Result() string {
	return resp.result
}

// Initializes SetCacheResponse to handle gRPC set response.
type SetCacheResponse struct {
	value []byte
}

// Decodes and returns byte value set in cache to string.
func (resp *SetCacheResponse) StringValue() string {
	return string(resp.value)
}

// Returns byte value set in cache.
func (resp *SetCacheResponse) ByteValue() []byte {
	return resp.value
}
