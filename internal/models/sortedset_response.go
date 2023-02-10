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

// SortedSetGetScoreMiss Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoreMiss struct{}

func (_ SortedSetGetScoreMiss) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreHit Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoreHit struct {
	Elements []*SortedSetScore
}

func (_ SortedSetGetScoreHit) isSortedSetGetScoreResponse() {}

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMiss Miss Response to a cache SortedSetRnk api request.
type SortedSetGetRankMiss struct{}

func (_ SortedSetGetRankMiss) isSortedSetGetRankResponse() {}

// SortedSetGetRankHit Hit Response to a cache SortedSetRank api request.
type SortedSetGetRankHit struct {
	Rank   uint64
	Status CacheResult
}

func (_ SortedSetGetRankHit) isSortedSetGetRankResponse() {}
