package responses

type SortedSetPutElementsResponse interface {
	isSortedSetPutResponse()
}

type SortedSetPutSuccess struct{}

func (SortedSetPutSuccess) isSortedSetPutResponse() {}
