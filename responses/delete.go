package responses

// DeleteResponse is the base response type for a delete request.
type DeleteResponse interface {
	MomentoCacheResponse
	isDeleteResponse()
}

// DeleteSuccess indicates a successful delete request.
type DeleteSuccess struct{}

func (DeleteSuccess) isDeleteResponse() {}
