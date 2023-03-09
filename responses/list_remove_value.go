package responses

// ListRemoveValueResponse is a base type for a list remove value request.
type ListRemoveValueResponse interface {
	isListRemoveValueResponse()
}

// ListRemoveValueSuccess indicates a successful list remove value request.
type ListRemoveValueSuccess struct{}

func (ListRemoveValueSuccess) isListRemoveValueResponse() {}
