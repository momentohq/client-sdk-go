package responses

type SortedSetPutResponse interface {
	isSortedSetPutResponse()
}

type SortedSetPutSuccess struct{}

func (SortedSetPutSuccess) isSortedSetPutResponse() {}
