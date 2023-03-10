package responses

type SortedSetRemoveElementsResponse interface {
	isSortedSetRemoveElementsResponse()
}

type SortedSetRemoveElementsSuccess struct{}

func (SortedSetRemoveElementsSuccess) isSortedSetRemoveElementsResponse() {}
