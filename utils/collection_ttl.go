package utils

import "time"

type CollectionTtl struct {
	Ttl        time.Duration
	RefreshTtl bool
}

func FromCacheTtl() CollectionTtl {
	return CollectionTtl{RefreshTtl: true}
}

func Of(ttl time.Duration) CollectionTtl {
	return CollectionTtl{Ttl: ttl}
}

func RefreshTtlIfProvided(ttl ...time.Duration) CollectionTtl {
	if len(ttl) > 0 {
		return CollectionTtl{Ttl: ttl[0], RefreshTtl: true}
	}
	return CollectionTtl{RefreshTtl: false}
}

func WithRefreshTtlOnUpdates(currentTtl CollectionTtl) CollectionTtl {
	return CollectionTtl{
		Ttl:        currentTtl.Ttl,
		RefreshTtl: true,
	}
}

func WithNoRefreshTtlOnUpdates(currentTtl CollectionTtl) CollectionTtl {
	return CollectionTtl{
		Ttl:        currentTtl.Ttl,
		RefreshTtl: false,
	}
}
