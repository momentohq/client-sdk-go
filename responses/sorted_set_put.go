package responses

// SortedSetPutElementsResponse is the base reponse type for a sorted set put request.
type SortedSetPutElementsResponse interface {
	isSortedSetPutResponse()
}

// SortedSetPutSuccess indicates a successful sorted set put request.
type SortedSetPutSuccess struct{}

func (SortedSetPutSuccess) isSortedSetPutResponse() {}
