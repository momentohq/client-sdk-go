package responses

type SortedSetRemoveElementsResponse interface {
	isSortedSetRemoveResponse()
}

type SortedSetRemoveSuccess struct{}

func (SortedSetRemoveSuccess) isSortedSetRemoveResponse() {}
