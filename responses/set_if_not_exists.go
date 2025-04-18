package responses

// SetIfNotExistsResponse is the base response type for a SetIfNotExists request
type SetIfNotExistsResponse interface {
	MomentoCacheResponse
	isSetIfNotExistsResponse()
}

// SetIfNotExistsNotStored indicates a successful set request where the value was already present.
type SetIfNotExistsNotStored struct{}

func (SetIfNotExistsNotStored) isSetIfNotExistsResponse() {}

// SetIfNotExistsStored indicates a successful set request where the value was stored.
type SetIfNotExistsStored struct{}

func (SetIfNotExistsStored) isSetIfNotExistsResponse() {}
