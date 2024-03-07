package responses

// SetIfNotEqualResponse is the base response type for a SetIfNotEqual request
type SetIfNotEqualResponse interface {
	isSetIfNotEqualResponse()
}

// SetIfNotEqualNotStored indicates a successful set request where the value was already present.
type SetIfNotEqualNotStored struct{}

func (SetIfNotEqualNotStored) isSetIfNotEqualResponse() {}

// SetIfNotEqualStored indicates a successful set request where the value was stored.
type SetIfNotEqualStored struct{}

func (SetIfNotEqualStored) isSetIfNotEqualResponse() {}
