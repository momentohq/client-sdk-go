package responses

type SortedSetPutElementResponse interface {
	isSortedSetPutElementResponse()
}

type SortedSetPutElementSuccess struct{}

func (SortedSetPutElementSuccess) isSortedSetPutElementResponse() {}
