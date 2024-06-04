package responses

// StorageDeleteResponse is the base response type for a store delete request.
type StorageDeleteResponse interface {
	isStorageDeleteResponse()
}

// StorageDeleteSuccess indicates a successful store delete request.
type StorageDeleteSuccess struct{}

func (StorageDeleteSuccess) isStorageDeleteResponse() {}
