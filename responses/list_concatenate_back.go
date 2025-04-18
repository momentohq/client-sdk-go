package responses

// ListConcatenateBackResponse is the base response type for a list concatenate back request.
type ListConcatenateBackResponse interface {
	MomentoCacheResponse
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

// NewListConcatenateBackSuccess returns a new ListConcatenateBackSuccess containing the supplied length.
func NewListConcatenateBackSuccess(listLength uint32) *ListConcatenateBackSuccess {
	return &ListConcatenateBackSuccess{listLength: listLength}
}
