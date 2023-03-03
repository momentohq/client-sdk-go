package responses

type SortedSetElement struct {
	Value []byte
	Score float64
}
type SortedSetFetchResponse interface {
	isSortedSetFetchResponse()
}

// SortedSetFetchMiss Miss Response to a cache SortedSetFetch api request.
type SortedSetFetchMiss struct{}

func (SortedSetFetchMiss) isSortedSetFetchResponse() {}

// SortedSetFetchHit Hit Response to a cache SortedSetFetch api request.
type SortedSetFetchHit struct {
	Elements []*SortedSetElement
}

func (SortedSetFetchHit) isSortedSetFetchResponse() {}

func NewSortedSetFetchHit(elements []*SortedSetElement) *SortedSetFetchHit {
	return &SortedSetFetchHit{Elements: elements}
}
