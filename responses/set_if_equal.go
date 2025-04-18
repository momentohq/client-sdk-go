package responses

// SetIfEqualResponse is the base response type for a SetIfEqual request
type SetIfEqualResponse interface {
	MomentoCacheResponse
	isSetIfEqualResponse()
}

// SetIfEqualNotStored indicates a successful set request where the value was already present.
type SetIfEqualNotStored struct{}

func (SetIfEqualNotStored) isSetIfEqualResponse() {}

// SetIfEqualStored indicates a successful set request where the value was stored.
type SetIfEqualStored struct{}

func (SetIfEqualStored) isSetIfEqualResponse() {}
