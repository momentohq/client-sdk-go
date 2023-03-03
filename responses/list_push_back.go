package responses

type ListPushBackResponse interface {
	isListPushBackResponse()
}

type ListPushBackSuccess struct {
	value uint32
}

func (ListPushBackSuccess) isListPushBackResponse() {}

func (resp ListPushBackSuccess) ListLength() uint32 {
	return resp.value
}

func NewListPushBackSuccess(value uint32) *ListPushBackSuccess {
	return &ListPushBackSuccess{value: value}
}
