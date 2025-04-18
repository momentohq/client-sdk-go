package responses

// SortedSetLengthByScoreResponse is the base response type for a sorted set length request.
type SortedSetLengthByScoreResponse interface {
	MomentoCacheResponse
	isSortedSetLengthByScoreResponse()
}

// SortedSetLengthByScoreHit indicates a sorted set length request was a hit.
type SortedSetLengthByScoreHit struct {
	value uint32
}

func (SortedSetLengthByScoreHit) isSortedSetLengthByScoreResponse() {}

// Length returns the length of the sorted set.
func (resp SortedSetLengthByScoreHit) Length() uint32 {
	return resp.value
}

// SortedSetLengthByScoreMiss indicates a sorted set length request was a miss.
type SortedSetLengthByScoreMiss struct{}

func (SortedSetLengthByScoreMiss) isSortedSetLengthByScoreResponse() {}

// NewSortedSetLengthByScoreHit returns a new SortedSetLengthByScoreHit containing the supplied value.
func NewSortedSetLengthByScoreHit(value uint32) *SortedSetLengthByScoreHit {
	return &SortedSetLengthByScoreHit{value: value}
}
