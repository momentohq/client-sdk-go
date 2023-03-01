package utils

import "time"

type CollectionTtl struct {
	Ttl        time.Duration
	RefreshTtl *bool
}

func FromCacheTtl() CollectionTtl {
	t := true
	return CollectionTtl{RefreshTtl: &t}
}

func Of(ttl time.Duration) CollectionTtl {
	return CollectionTtl{Ttl: ttl}
}

func RefreshTtlIfProvided(ttl ...time.Duration) CollectionTtl {
	if len(ttl) > 0 {
		t := true
		return CollectionTtl{Ttl: ttl[0], RefreshTtl: &t}
	}
	f := false
	return CollectionTtl{RefreshTtl: &f}
}

func WithRefreshTtlOnUpdates(currentTtl CollectionTtl) CollectionTtl {
	t := true
	return CollectionTtl{
		Ttl:        currentTtl.Ttl,
		RefreshTtl: &t,
	}
}

func WithNoRefreshTtlOnUpdates(currentTtl CollectionTtl) CollectionTtl {
	f := false
	return CollectionTtl{
		Ttl:        currentTtl.Ttl,
		RefreshTtl: &f,
	}
}
