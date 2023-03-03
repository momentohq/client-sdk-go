package responses

type SortedSetRemoveResponse interface {
	isSortedSetRemoveResponse()
}

type SortedSetRemoveSuccess struct{}

func (SortedSetRemoveSuccess) isSortedSetRemoveResponse() {}
