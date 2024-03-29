package responses

type SortedSetBytesElement struct {
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
	elements []SortedSetBytesElement
}

func (SortedSetFetchHit) isSortedSetFetchResponse() {}

// NewSortedSetFetchHit returns a new SortedSetFetchHit containing the supplied elements.
func NewSortedSetFetchHit(elements []SortedSetBytesElement) *SortedSetFetchHit {
	if elements == nil {
		elements = []SortedSetBytesElement{}
	}
	return &SortedSetFetchHit{elements: elements}
}

// ValueStringElements returns the elements as an array of objects, each containing a `value` and `score` field.
// The value is a utf-8 string, decoded from the underlying byte array, and the score is a number.
func (r SortedSetFetchHit) ValueStringElements() []SortedSetStringElement {
	elementsString := make([]SortedSetStringElement, 0, len(r.elements))

	for _, element := range r.elements {
		stringElement := SortedSetStringElement{
			Value: string(element.Value),
			Score: element.Score,
		}
		elementsString = append(elementsString, stringElement)
	}

	return elementsString
}

// ValueBytesElements returns the elements as an array of objects, each containing a `value` and `score` field.
// The value is a byte array, and the score is a number.
func (r SortedSetFetchHit) ValueBytesElements() []SortedSetBytesElement {
	return r.elements
}
