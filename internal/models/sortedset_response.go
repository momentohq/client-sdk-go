package models

type CacheResult int32

const (
	Hit  CacheResult = 2
	Miss CacheResult = 3
)

type SortedSetScore struct {
	Result CacheResult
	Score  float64
}

type SortedSetFetchResponse interface {
	isSortedSetFetchResponse()
}

// SortedSetFetchMissing Miss Response to a cache SortedSetFetch api request.
type SortedSetFetchMissing struct{}

func (_ SortedSetFetchMissing) isSortedSetFetchResponse() {}

// SortedSetFetchFound Hit Response to a cache SortedSetFetch api request.
type SortedSetFetchFound struct {
	Elements []*SortedSetElement
}

func (_ SortedSetFetchFound) isSortedSetFetchResponse() {}

type SortedSetGetScoreResponse interface {
	isSortedSetGetScoreResponse()
}

// SortedSetGetScoreMissing Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoreMissing struct{}

func (_ SortedSetGetScoreMissing) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreFound Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoreFound struct {
	Elements []*SortedSetScore
}

func (_ SortedSetGetScoreFound) isSortedSetGetScoreResponse() {}

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMissing Miss Response to a cache SortedSetRnk api request.
type SortedSetGetRankMissing struct{}

func (_ SortedSetGetRankMissing) isSortedSetGetRankResponse() {}

// SortedSetGetRankFound Hit Response to a cache SortedSetRank api request.
type SortedSetGetRankFound struct {
	Rank   uint64
	Status CacheResult
}

func (_ SortedSetGetRankFound) isSortedSetGetRankResponse() {}
