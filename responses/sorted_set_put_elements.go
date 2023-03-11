package responses

// SortedSetPutElementsResponse is the base response type for a sorted set put elements request.
type SortedSetPutElementsResponse interface {
	isSortedSetPutElementsResponse()
}

// SortedSetPutElementsSuccess indicates a successful sorted set put elements request.
type SortedSetPutElementsSuccess struct{}

func (SortedSetPutElementsSuccess) isSortedSetPutElementsResponse() {}
