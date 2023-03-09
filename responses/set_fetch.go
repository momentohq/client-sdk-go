package responses

// SetFetchResponse is a base response type for a set fetch request.
type SetFetchResponse interface {
	isSetFetchResponse()
}

// SetFetchHit indicates a set fetch request was a hit.
type SetFetchHit struct {
	elements       [][]byte
	elementsString []string
}

func (SetFetchHit) isSetFetchResponse() {}

// ValueString returns the data as a Set whose values are utf-8 strings, decoded from the underlying byte arrays.
func (resp SetFetchHit) ValueString() []string {
	if resp.elementsString == nil {
		for _, value := range resp.elements {
			resp.elementsString = append(resp.elementsString, string(value))
		}
	}
	return resp.elementsString
}

// ValueByte returns the data as a Set whose values are byte arrays.
func (resp SetFetchHit) ValueByte() [][]byte {
	return resp.elements
}

// SetFetchMiss indicates a set fetch request was a miss.
type SetFetchMiss struct{}

func (SetFetchMiss) isSetFetchResponse() {}

// NewSetFetchHit returns a new SetFetchHit contains elements.
func NewSetFetchHit(elements [][]byte) *SetFetchHit {
	return &SetFetchHit{
		elements: elements,
	}
}
