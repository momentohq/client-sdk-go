package responses

type SortedSetRemoveElementResponse interface {
	isSortedSetRemoveElementResponse()
}

type SortedSetRemoveElementSuccess struct{}

func (SortedSetRemoveElementSuccess) isSortedSetRemoveElementResponse() {}
