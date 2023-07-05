package responses

// SortedSetLengthResponse is the base response type for a sorted set length request.
type SortedSetLengthResponse interface {
	isSortedSetLengthResponse()
}

// SortedSetLengthHit indicates a sorted set length request was a hit.
type SortedSetLengthHit struct {
	value uint32
}

func (SortedSetLengthHit) isSortedSetLengthResponse() {}

// Length returns the length of the sorted set.
func (resp SortedSetLengthHit) Length() uint32 {
	return resp.value
}

// SortedSetLengthMiss indicates a sorted set length request was a miss.
type SortedSetLengthMiss struct{}

func (SortedSetLengthMiss) isSortedSetLengthResponse() {}

// NewSortedSetLengthHit returns a new SortedSetLengthHit containing the supplied value.
func NewSortedSetLengthHit(value uint32) *SortedSetLengthHit {
	return &SortedSetLengthHit{value: value}
}
