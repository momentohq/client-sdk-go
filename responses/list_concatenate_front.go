package responses

// ListConcatenateFrontResponse is the base response for a list concatenate front request.
type ListConcatenateFrontResponse interface {
	isListConcatenateFrontResponse()
}

// ListConcatenateFrontSuccess returns a new ListConcatenateFrontSuccess containing the supplied length.
type ListConcatenateFrontSuccess struct {
	listLength uint32
}

func (ListConcatenateFrontSuccess) isListConcatenateFrontResponse() {}

// ListLength returns the length of the given list.
func (resp ListConcatenateFrontSuccess) ListLength() uint32 {
	return resp.listLength
}

// NewListConcatenateFrontSuccess returns a new ListConcatenateFrontSuccess containing the supplied length.
func NewListConcatenateFrontSuccess(listLength uint32) *ListConcatenateFrontSuccess {
	return &ListConcatenateFrontSuccess{listLength: listLength}
}
