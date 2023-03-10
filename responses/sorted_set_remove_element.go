package responses

// SortedSetRemoveElementResponse is the base response type for a sorted set remove element request.
type SortedSetRemoveElementResponse interface {
	isSortedSetRemoveElementResponse()
}

// SortedSetRemoveElementSuccess indicates a successful sorted set remove element request.
type SortedSetRemoveElementSuccess struct{}

func (SortedSetRemoveElementSuccess) isSortedSetRemoveElementResponse() {}
