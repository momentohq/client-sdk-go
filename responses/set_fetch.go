package responses

type SetFetchResponse interface {
	isSetFetchResponse()
}

type SetFetchHit struct {
	elements       [][]byte
	elementsString []string
}

func (SetFetchHit) isSetFetchResponse() {}

func (resp SetFetchHit) ValueString() []string {
	if resp.elementsString == nil {
		for _, value := range resp.elements {
			resp.elementsString = append(resp.elementsString, string(value))
		}
	}
	return resp.elementsString
}

func (resp SetFetchHit) ValueByte() [][]byte {
	return resp.elements
}

type SetFetchMiss struct{}

func (SetFetchMiss) isSetFetchResponse() {}

func NewSetFetchHit(elements [][]byte) *SetFetchHit {
	return &SetFetchHit{
		elements: elements,
	}
}
