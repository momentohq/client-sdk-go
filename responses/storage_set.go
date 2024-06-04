package responses

// StorageSetResponse is the base response type for a store put request.
type StorageSetResponse interface {
	isStorePutResponse()
}

// StorageSetSuccess indicates a successful store put request.
type StorageSetSuccess struct{}

func (StorageSetSuccess) isStorePutResponse() {}
