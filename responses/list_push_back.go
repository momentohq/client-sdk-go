package responses

// ListPushBackResponse is the base response type for a list push back request.
type ListPushBackResponse interface {
	MomentoCacheResponse
	isListPushBackResponse()
}

// ListPushBackSuccess indicates a successful list push back request.
type ListPushBackSuccess struct {
	value uint32
}

func (ListPushBackSuccess) isListPushBackResponse() {}

// ListLength returns the new length of the list after the push operation.
func (resp ListPushBackSuccess) ListLength() uint32 {
	return resp.value
}

// NewListPushBackSuccess returns a new ListPushBackSuccess containing the supplied value.
func NewListPushBackSuccess(value uint32) *ListPushBackSuccess {
	return &ListPushBackSuccess{value: value}
}
