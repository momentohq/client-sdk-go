package responses

// ListConcatenateFrontResponse is a base reponse for a list concatenate front request.
type ListConcatenateFrontResponse interface {
	isListConcatenateFrontResponse()
}

// ListConcatenateFrontSuccess
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
