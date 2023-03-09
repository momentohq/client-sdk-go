package responses

// SortedSetGetRankResponse is a base response type for a sorted set get rank request.
type SortedSetGetRankResponse interface {
	isSortedSetGetRankResponse()
}

// SortedSetGetRankMiss Miss Response to a cache SortedSetGetRank api request.
type SortedSetGetRankMiss struct{}

func (SortedSetGetRankMiss) isSortedSetGetRankResponse() {}

// SortedSetGetRankHit Hit Response to a cache SortedSetGetRank api request.
type SortedSetGetRankHit uint64

func (SortedSetGetRankHit) isSortedSetGetRankResponse() {}

// Rank returns the rank of the element in the sorted set.  Ranks start at 0.
func (r SortedSetGetRankHit) Rank() uint64 { return uint64(r) }
