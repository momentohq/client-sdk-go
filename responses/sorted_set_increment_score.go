package responses

type SortedSetIncrementScoreResponse interface {
	isSortedSetIncrementResponse()
}
type SortedSetIncrementScoreSuccess float64

func (SortedSetIncrementScoreSuccess) isSortedSetIncrementResponse() {}

func (r SortedSetIncrementScoreSuccess) Score() float64 { return float64(r) }
