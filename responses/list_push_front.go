package responses

type ListPushFrontResponse interface {
	isListPushFrontResponse()
}

type ListPushFrontSuccess struct {
	value uint32
}

func (ListPushFrontSuccess) isListPushFrontResponse() {}

func (resp ListPushFrontSuccess) ListLength() uint32 {
	return resp.value
}

func NewListPushFrontSuccess(value uint32) *ListPushFrontSuccess {
	return &ListPushFrontSuccess{value: value}
}
