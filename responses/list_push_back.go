package responses

// ListPushBackResponse is a base response type for a list push back request.
type ListPushBackResponse interface {
	isListPushBackResponse()
}

// ListPushBackSuccess indicates a successful lit push back request.
type ListPushBackSuccess struct {
	value uint32
}

func (ListPushBackSuccess) isListPushBackResponse() {}

// ListLength returns the new length of the list after the push operation.
func (resp ListPushBackSuccess) ListLength() uint32 {
	return resp.value
}

// NewListPushBackSuccess returns a new ListPushBackSuccess contains value.
func NewListPushBackSuccess(value uint32) *ListPushBackSuccess {
	return &ListPushBackSuccess{value: value}
}
