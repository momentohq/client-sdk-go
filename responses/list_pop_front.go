package responses

type ListPopFrontResponse interface {
	isListPopFrontResponse()
}

type ListPopFrontHit struct {
	value []byte
}

func (ListPopFrontHit) isListPopFrontResponse() {}

func (resp ListPopFrontHit) ValueByte() []byte {
	return resp.value
}

func (resp ListPopFrontHit) ValueString() string {
	return string(resp.value)
}

type ListPopFrontMiss struct{}

func (ListPopFrontMiss) isListPopFrontResponse() {}

func NewListPopFrontHit(value []byte) *ListPopFrontHit {
	return &ListPopFrontHit{value: value}
}
