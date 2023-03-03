package responses

type ListConcatenateBackResponse interface {
	isListConcatenateBackResponse()
}

type ListConcatenateBackSuccess struct {
	listLength uint32
}

func (ListConcatenateBackSuccess) isListConcatenateBackResponse() {}

func (resp ListConcatenateBackSuccess) ListLength() uint32 {
	return resp.listLength
}

func NewListConcatenateBackSuccess(listLength uint32) *ListConcatenateBackSuccess {
	return &ListConcatenateBackSuccess{listLength: listLength}
}
