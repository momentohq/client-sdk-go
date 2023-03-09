package responses

// ListPushFrontResponse is a base type for a list push front request.
type ListPushFrontResponse interface {
	isListPushFrontResponse()
}

// ListPushFrontSuccess indicates a successful list push front request.
type ListPushFrontSuccess struct {
	value uint32
}

func (ListPushFrontSuccess) isListPushFrontResponse() {}

// ListLength returns the new length of the list after the push operation.
func (resp ListPushFrontSuccess) ListLength() uint32 {
	return resp.value
}

// NewListPushFrontSuccess returns a new ListPushFrontSuccess contains value.
func NewListPushFrontSuccess(value uint32) *ListPushFrontSuccess {
	return &ListPushFrontSuccess{value: value}
}
