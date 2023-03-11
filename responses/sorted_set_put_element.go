package responses

// SortedSetPutElementResponse is the base response type for a sorted set put element request.
type SortedSetPutElementResponse interface {
	isSortedSetPutElementResponse()
}

// SortedSetPutElementSuccess indicates a successful sorted set put element request.
type SortedSetPutElementSuccess struct{}

func (SortedSetPutElementSuccess) isSortedSetPutElementResponse() {}
