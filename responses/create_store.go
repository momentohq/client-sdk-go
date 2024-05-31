package responses

// CreateStoreResponse is the base response type for a create store request.
type CreateStoreResponse interface {
	isCreateStoreResponse()
}

// CreateStoreSuccess indicates a successful create store request.
type CreateStoreSuccess struct{}

func (CreateStoreSuccess) isCreateStoreResponse() {}

// CreateStoreAlreadyExists indicates that the store already exists, so there was nothing to do.
type CreateStoreAlreadyExists struct{}

func (CreateStoreAlreadyExists) isCreateStoreResponse() {}
