package momento

// Initializes ListCacheResponse to handle list cache response.
type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

// Returns next token.
func (resp *ListCachesResponse) NextToken() string {
	return resp.nextToken
}

// Returns all caches.
func (resp *ListCachesResponse) Caches() []CacheInfo {
	return resp.caches
}

// Initializes CacheInfo to handle caches returned from list cache operation.
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
