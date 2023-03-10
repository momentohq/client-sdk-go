package responses

// DeleteResponse is the base response type for a delete request.
type DeleteResponse interface {
	isDeleteResponse()
}

// DeleteSuccess indicates a successful delete request.
type DeleteSuccess struct{}

func (DeleteSuccess) isDeleteResponse() {}
