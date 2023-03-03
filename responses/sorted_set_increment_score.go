package responses

type SortedSetIncrementScoreResponse interface {
	isSortedSetIncrementResponse()
}
type SortedSetIncrementScoreSuccess struct {
	Value float64
}

func (SortedSetIncrementScoreSuccess) isSortedSetIncrementResponse() {}
