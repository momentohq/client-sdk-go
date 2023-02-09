package incubating

type SortedSetElement struct {
	Name  []byte
	Score float64
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
	Elements []SortedSetScoreElement
}

func (_ SortedSetGetScoreFound) isSortedSetGetScoreResponse() {}

type SortedSetScoreElement interface {
	isSortedSetScoreElement()
}

type SortedSetScoreHit struct {
	Score float64
}

func (_ SortedSetScoreHit) isSortedSetScoreElement() {}

type SortedSetScoreMiss struct{}

func (_ SortedSetScoreMiss) isSortedSetScoreElement() {}

type SortedSetScoreInvalid struct{}

func (_ SortedSetScoreInvalid) isSortedSetScoreElement() {}

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMissing Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMissing struct{}

func (_ SortedSetGetRankMissing) isSortedSetGetRankResponse() {}

// SortedSetGetRankFound Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankFound struct {
	Element SortedSetRankElement
}

func (_ SortedSetGetRankFound) isSortedSetGetRankResponse() {}

type SortedSetRankElement interface {
	isSortedSetRankElement()
}

type SortedSetRankHit struct {
	Rank uint64
}

func (_ SortedSetRankHit) isSortedSetRankElement() {}

type SortedSetRankMiss struct{}

func (_ SortedSetRankMiss) isSortedSetRankElement() {}
