package incubating

type SortedSetElement struct {
	Name  []byte
	Score float64
}
type SortedSetFetchResponse interface {
	isSortedSetFetchResponse()
}

// SortedSetFetchMiss Miss Response to a cache SortedSetFetch api request.
type SortedSetFetchMiss struct{}

func (SortedSetFetchMiss) isSortedSetFetchResponse() {}

// SortedSetFetchHit Hit Response to a cache SortedSetFetch api request.
type SortedSetFetchHit struct {
	Elements []*SortedSetElement
}

func (SortedSetFetchHit) isSortedSetFetchResponse() {}

type SortedSetGetScoreResponse interface {
	isSortedSetGetScoreResponse()
}

// SortedSetGetScoreMiss Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoreMiss struct{}

func (SortedSetGetScoreMiss) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreHit Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoreHit struct {
	Elements []SortedSetScoreElement
}

func (SortedSetGetScoreHit) isSortedSetGetScoreResponse() {}

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

// SortedSetGetRankMiss Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMiss struct{}

func (SortedSetGetRankMiss) isSortedSetGetRankResponse() {}

// SortedSetGetRankHit Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankHit struct {
	Element SortedSetRankElement
}

func (SortedSetGetRankHit) isSortedSetGetRankResponse() {}

type SortedSetRankElement interface {
	isSortedSetRankElement()
}

type SortedSetRankHit struct {
	Rank uint64
}

func (SortedSetRankHit) isSortedSetRankElement() {}

type SortedSetRankMiss struct{}

func (SortedSetRankMiss) isSortedSetRankElement() {}

type SortedSetIncrementResponse interface {
	isSortedSetIncrementResponse()
}
type SortedSetIncrementResponseSuccess struct {
	Value float64
}

func (_ SortedSetIncrementResponseSuccess) isSortedSetIncrementResponse() {}
