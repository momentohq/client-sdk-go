package responses

// StoragePutResponse is the base response type for a store put request.
type StoragePutResponse interface {
	isStorePutResponse()
}

// StoragePutSuccess indicates a successful store put request.
type StoragePutSuccess struct{}

func (StoragePutSuccess) isStorePutResponse() {}
