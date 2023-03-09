package responses

// SortedSetPutResponse is a base reponse type for a sorted set put request.
type SortedSetPutResponse interface {
	isSortedSetPutResponse()
}

// SortedSetPutSuccess indicates a successful sorted set put request.
type SortedSetPutSuccess struct{}

func (SortedSetPutSuccess) isSortedSetPutResponse() {}
