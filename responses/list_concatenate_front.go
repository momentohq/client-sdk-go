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

func (resp ListConcatenateFrontSuccess) ListLength() uint32 {
	return resp.listLength
}

func NewListConcatenateFrontSuccess(listLength uint32) *ListConcatenateFrontSuccess {
	return &ListConcatenateFrontSuccess{listLength: listLength}
}
