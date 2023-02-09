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

func (SortedSetFetchMissing) isSortedSetFetchResponse() {}

// SortedSetFetchFound Hit Response to a cache SortedSetFetch api request.
type SortedSetFetchFound struct {
	Elements []*SortedSetElement
}

func (SortedSetFetchFound) isSortedSetFetchResponse() {}

type SortedSetGetScoreResponse interface {
	isSortedSetGetScoreResponse()
}

// SortedSetGetScoreMissing Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoreMissing struct{}

func (SortedSetGetScoreMissing) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreFound Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoreFound struct {
	Elements []SortedSetScoreElement
}

func (SortedSetGetScoreFound) isSortedSetGetScoreResponse() {}

type SortedSetScoreElement interface {
	isSortedSetScoreElement()
}

type SortedSetScoreHit struct {
	Score float64
}

func (SortedSetScoreHit) isSortedSetScoreElement() {}

type SortedSetScoreMiss struct{}

func (SortedSetScoreMiss) isSortedSetScoreElement() {}

type SortedSetScoreInvalid struct{}

func (SortedSetScoreInvalid) isSortedSetScoreElement() {}

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMissing Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMissing struct{}

func (SortedSetGetRankMissing) isSortedSetGetRankResponse() {}

// SortedSetGetRankFound Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankFound struct {
	Element SortedSetRankElement
}

func (SortedSetGetRankFound) isSortedSetGetRankResponse() {}

type SortedSetRankElement interface {
	isSortedSetRankElement()
}

type SortedSetRankHit struct {
	Rank uint64
}

func (SortedSetRankHit) isSortedSetRankElement() {}

type SortedSetRankMiss struct{}

func (SortedSetRankMiss) isSortedSetRankElement() {}
