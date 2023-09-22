package responses

// SetPopResponse is the base response type for a set pop request.
type SetPopResponse interface {
	isSetPopResponse()
}

// SetPopHit indicates a set pop request was a hit.
type SetPopHit struct {
	elements       [][]byte
	elementsString []string
}

func (SetPopHit) isSetPopResponse() {}

// ValueString returns the data as a Set whose values are utf-8 strings, decoded from the underlying byte arrays.
func (resp SetPopHit) ValueString() []string {
	if resp.elementsString == nil {
		for _, value := range resp.elements {
			resp.elementsString = append(resp.elementsString, string(value))
		}
	}
	return resp.elementsString
}

// ValueByte returns the data as a Set whose values are byte arrays.
func (resp SetPopHit) ValueByte() [][]byte {
	return resp.elements
}

// SetPopMiss indicates a set pop request was a miss.
type SetPopMiss struct{}

func (SetPopMiss) isSetPopResponse() {}

// NewSetPopHit returns a new SetPopHit containing the supplied elements.
func NewSetPopHit(elements [][]byte) *SetPopHit {
	return &SetPopHit{
		elements: elements,
	}
}
