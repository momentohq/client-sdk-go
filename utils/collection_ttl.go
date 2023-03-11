package utils

import "time"

// CollectionTtl represents the desired behavior for managing the Ttl on collection objects (dictionaries, lists, sets, and sorted sets) in your cache.
// For cache operations that modify a collection, there are a few things to consider.
// The first time the collection is created, we need to set a Ttl on it.
// For subsequent operations that modify the collection you may choose to update the Ttl in order to prolong
// the life of the cached collection object, or you may choose to leave the Ttl unmodified in order to ensure that the collection expires at the original TTL.
// The default behavior is to refresh the Ttl (to prolong the life of the collection) each time it is written.
// This behavior can be modified by calling WithNoRefreshTtlOnUpdates.
type CollectionTtl struct {
	// Ttl is the time.Duration after which the cached collection should be expired from the cache.
	// If CollectionTtl is not provided, the default Ttl that was passed to a momento.CacheClient instance is used.
	Ttl time.Duration
	// If RefreshTtl is true, then the collection's Ttl will be refreshed (to prolong the life of the collection) on every update.
	// If false, then the collection's Ttl will only be set when the collection is initially created.
	RefreshTtl bool
}

// FromCacheTtl is the default way to handle Ttls for collections.
// The default Ttl that was specified when constructing momento.CacheClient will be used,
// and the Ttl for the collection will be refreshed any time the collection is modified.
func FromCacheTtl() CollectionTtl {
	return CollectionTtl{RefreshTtl: true}
}

// Of constructs a CollectionTtl with the specified time.Duration.
// Ttl for the collection will be refreshed any time the collection is modified.
func Of(ttl time.Duration) CollectionTtl {
	return CollectionTtl{Ttl: ttl}
}

// RefreshTtlIfProvided constructs a CollectionTtl with the specified time.Duration.
// Will only refresh if the Ttl is provided.
func RefreshTtlIfProvided(ttl ...time.Duration) CollectionTtl {
	if len(ttl) > 0 {
		return CollectionTtl{Ttl: ttl[0], RefreshTtl: true}
	}
	return CollectionTtl{RefreshTtl: false}
}

// WithRefreshTtlOnUpdates specifies that the TTL for the collection should be refreshed
// when the collection is modified.  (This is the default behavior.)
func WithRefreshTtlOnUpdates(currentTtl CollectionTtl) CollectionTtl {
	return CollectionTtl{
		Ttl:        currentTtl.Ttl,
		RefreshTtl: true,
	}
}

// WithNoRefreshTtlOnUpdates specifies that the TTL for the collection should not be refreshed
// when the collection is modified.  Use this if you want to ensure
// that your collection expires at the originally specified time, even if you make modifications to the value of the collection.
func WithNoRefreshTtlOnUpdates(currentTtl CollectionTtl) CollectionTtl {
	return CollectionTtl{
		Ttl:        currentTtl.Ttl,
		RefreshTtl: false,
	}
}
