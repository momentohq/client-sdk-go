package responses

// DeleteStoreResponse is the base response type for a delete store request.
type DeleteStoreResponse interface {
	isDeleteStoreResponse()
}

// DeleteStoreSuccess indicates a successful delete store request.
type DeleteStoreSuccess struct{}

func (DeleteStoreSuccess) isDeleteStoreResponse() {}
