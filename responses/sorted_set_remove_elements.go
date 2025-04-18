package responses

// SortedSetRemoveElementsResponse is the base response type for a sorted set put elements request.
type SortedSetRemoveElementsResponse interface {
	MomentoCacheResponse
	isSortedSetRemoveElementsResponse()
}

// SortedSetRemoveElementsSuccess indicates a successful sorted set put elements request.
type SortedSetRemoveElementsSuccess struct{}

func (SortedSetRemoveElementsSuccess) isSortedSetRemoveElementsResponse() {}
