package responses

// StorePutResponse is the base response type for a store put request.
type StorePutResponse interface {
	isStorePutResponse()
}

// StorePutSuccess indicates a successful store put request.
type StorePutSuccess struct{}

func (StorePutSuccess) isStorePutResponse() {}
