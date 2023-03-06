package responses

type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMiss Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMiss struct{}

func (SortedSetGetRankMiss) isSortedSetGetRankResponse() {}

// SortedSetGetRankHit Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankHit uint64

func (SortedSetGetRankHit) isSortedSetGetRankResponse() {}

func (r SortedSetGetRankHit) Rank() uint64 { return uint64(r) }
