package momento

// Bytes generic interface to help users deal with passing us bytes
type Bytes interface{ asBytes() []byte }

// RawBytes plain old []byte
type RawBytes struct {
	Bytes []byte
}

func (r RawBytes) asBytes() []byte {
	return r.Bytes
}

// StringBytes string type that will be converted to []byte
type StringBytes struct {
	Text string
}

func (r StringBytes) asBytes() []byte {
	return []byte(r.Text)
}

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

type CacheSetRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Bytes
	// string ot byte value to be stored.
	Value Bytes
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	TTLSeconds TimeToLive
}

// TTL Helper function that returns an initialized TimeToLive containing a pointer to ttl.
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
	// string or byte key to be used to store item
	Key Bytes
}

type CacheDeleteRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to delete the item.
	Key Bytes
}
