package responses

// SortedSetIncrementScoreResponse is the base response type for a sorted set increment score request.
type SortedSetIncrementScoreResponse interface {
	isSortedSetIncrementResponse()
}

// SortedSetIncrementScoreSuccess indicates a successful sorted set increment score request.
type SortedSetIncrementScoreSuccess float64

func (SortedSetIncrementScoreSuccess) isSortedSetIncrementResponse() {}

// Score returns the new score of the element after incrementing.
func (r SortedSetIncrementScoreSuccess) Score() float64 { return float64(r) }
