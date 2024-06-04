package responses

// StoreDeleteResponse is the base response type for a store delete request.
type StoreDeleteResponse interface {
	isStoreDeleteResponse()
}

// StoreDeleteSuccess indicates a successful store delete request.
type StoreDeleteSuccess struct{}

func (StoreDeleteSuccess) isStoreDeleteResponse() {}
