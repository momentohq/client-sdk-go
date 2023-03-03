package responses

type ListPopBackResponse interface {
	isListPopBackResponse()
}

type ListPopBackHit struct {
	value []byte
}

func (ListPopBackHit) isListPopBackResponse() {}

func (resp ListPopBackHit) ValueByte() []byte {
	return resp.value
}

func (resp ListPopBackHit) ValueString() string {
	return string(resp.value)
}

type ListPopBackMiss struct{}

func (ListPopBackMiss) isListPopBackResponse() {}

func NewListPopBackHit(value []byte) *ListPopBackHit {
	return &ListPopBackHit{value: value}
}
