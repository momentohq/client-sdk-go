package responses

// ListConcatenateBackResponse is a base response type for a list concatenate back request.
type ListConcatenateBackResponse interface {
	isListConcatenateBackResponse()
}

// ListConcatenateBackSuccess indicates a successful list concatenate back request.
type ListConcatenateBackSuccess struct {
	listLength uint32
}

func (ListConcatenateBackSuccess) isListConcatenateBackResponse() {}

// ListLength returns the length of the given list.
func (resp ListConcatenateBackSuccess) ListLength() uint32 {
	return resp.listLength
}

// NewListConcatenateBackSuccess returns a new ListConcatenateBackSuccess contains length.
func NewListConcatenateBackSuccess(listLength uint32) *ListConcatenateBackSuccess {
	return &ListConcatenateBackSuccess{listLength: listLength}
}
