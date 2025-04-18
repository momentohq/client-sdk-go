package responses

// ListRemoveValueResponse is the base type for a list remove value request.
type ListRemoveValueResponse interface {
	MomentoCacheResponse
	isListRemoveValueResponse()
}

// ListRemoveValueSuccess indicates a successful list remove value request.
type ListRemoveValueSuccess struct{}

func (ListRemoveValueSuccess) isListRemoveValueResponse() {}
