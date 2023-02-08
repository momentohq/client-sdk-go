package momento

import "time"

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
	Ttl time.Duration
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
