package responses

// SortedSetPutElementResponse is the base response type for a sorted set put element request.
type SortedSetPutElementResponse interface {
	MomentoCacheResponse
	isSortedSetPutElementResponse()
}

// SortedSetPutElementSuccess indicates a successful sorted set put element request.
type SortedSetPutElementSuccess struct{}

func (SortedSetPutElementSuccess) isSortedSetPutElementResponse() {}
