package responses

type SortedSetPutElementsResponse interface {
	isSortedSetPutElementsResponse()
}

type SortedSetPutElementsSuccess struct{}

func (SortedSetPutElementsSuccess) isSortedSetPutElementsResponse() {}
