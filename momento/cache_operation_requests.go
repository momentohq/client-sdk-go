package momento

type CreateCacheRequest struct {
	// string used to create a cache.
	CacheName string
}

type DeleteCacheRequest struct {
	// string cache name to delete.
	CacheName string
}

type ListCachesRequest struct {
	// Token to continue paginating through the list. It's used to handle large paginated lists.
	NextToken string
}

type CreateSigningKeyRequest struct {
	// Ttl in minutes before the Momento signing key expires
	TtlMinutes uint32
}

type RevokeSigningKeyRequest struct {
	// Id of Momento signing key to revoke
	KeyId string
}

type ListSigningKeysRequest struct {
	// Token to continue paginating through the list. It's used to handle large paginated lists.
	NextToken string
}

type CacheSetRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key interface{}
	// string ot byte value to be stored.
	Value interface{}
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	TtlSeconds TimeToLive
}

// Helper function that returns an initialized TimeToLive containing a pointer to ttl.
func TTL(ttl uint32) TimeToLive {
	pInt := &ttl
	return TimeToLive{
		_ttl: pInt,
	}
}

// TimeToLive provies a structure to hold a uint32 pointer of ttl.
type TimeToLive struct {
	_ttl *uint32
}

type CacheGetRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string ot byte key to be used to retrieve the item.
	Key interface{}
}

type CacheDeleteRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string ot byte key to be used to retrieve the item.
	Key interface{}
}
