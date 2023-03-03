package responses

type ListRemoveValueResponse interface {
	isListRemoveValueResponse()
}

type ListRemoveValueSuccess struct{}

func (ListRemoveValueSuccess) isListRemoveValueResponse() {}
