package responses

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMiss Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMiss struct{}

func (SortedSetGetRankMiss) isSortedSetGetRankResponse() {}

// SortedSetGetRankHit Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankHit struct {
	Rank uint64
}

func (SortedSetGetRankHit) isSortedSetGetRankResponse() {}
