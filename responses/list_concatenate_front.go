package responses

type ListConcatenateFrontResponse interface {
	isListConcatenateFrontResponse()
}

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
