package utils

import "time"

type CollectionTTL struct {
	Ttl        time.Duration
	RefreshTtl bool
}

func FromCacheTtl() CollectionTTL {
	return CollectionTTL{RefreshTtl: true}
}

func Of(ttl time.Duration) CollectionTTL {
	return CollectionTTL{Ttl: ttl}
}

func RefreshTtlIfProvided(ttl ...time.Duration) CollectionTTL {
	if len(ttl) > 0 {
		return CollectionTTL{Ttl: ttl[0], RefreshTtl: true}
	}
	return CollectionTTL{RefreshTtl: false}
}

func WithRefreshTtlOnUpdates(currentCollectionTtl CollectionTTL) CollectionTTL {
	return CollectionTTL{
		Ttl:        currentCollectionTtl.Ttl,
		RefreshTtl: true,
	}
}

func WithNoRefreshTtlOnUpdates(currentCollectionTtl CollectionTTL) CollectionTTL {
	return CollectionTTL{
		Ttl:        currentCollectionTtl.Ttl,
		RefreshTtl: false,
	}
}
