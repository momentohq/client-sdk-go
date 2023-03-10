package responses

// SortedSetRemoveElementsResponse is the base response type for a sorted set remove request.
type SortedSetRemoveElementsResponse interface {
	isSortedSetRemoveResponse()
}

// SortedSetRemoveSuccess indicates a successful sorted set remove request.
type SortedSetRemoveSuccess struct{}

func (SortedSetRemoveSuccess) isSortedSetRemoveResponse() {}
