package responses

type SortedSetElement struct {
	Value []byte
	Score float64
}

type SortedSetStringElement struct {
	Value string
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
	elements []*SortedSetElement
}

func (SortedSetFetchHit) isSortedSetFetchResponse() {}

func NewSortedSetFetchHit(elements []*SortedSetElement) *SortedSetFetchHit {
	return &SortedSetFetchHit{elements: elements}
}

func (r SortedSetFetchHit) ValueStringElements() []*SortedSetStringElement {
	elementsString := make([]*SortedSetStringElement, 0, len(r.elements))

	for _, element := range r.elements {
		stringElement := &SortedSetStringElement{
			Value: string(element.Value),
			Score: element.Score,
		}
		elementsString = append(elementsString, stringElement)
	}

	return elementsString
}

func (r SortedSetFetchHit) ValueByteElements() []*SortedSetElement {
	return r.elements
}
