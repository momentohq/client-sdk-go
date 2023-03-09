package responses

// SortedSetRemoveResponse is a base response type for a sorted set remove request.
type SortedSetRemoveResponse interface {
	isSortedSetRemoveResponse()
}

// SortedSetRemoveSuccess indicates a successful sorted set remove request.
type SortedSetRemoveSuccess struct{}

func (SortedSetRemoveSuccess) isSortedSetRemoveResponse() {}
