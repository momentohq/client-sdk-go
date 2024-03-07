package responses

// SetIfAbsentOrEqualResponse is the base response type for a SetIfAbsentOrEqual request
type SetIfAbsentOrEqualResponse interface {
	isSetIfAbsentOrEqualResponse()
}

// SetIfAbsentOrEqualNotStored indicates a successful set request where the value was already present.
type SetIfAbsentOrEqualNotStored struct{}

func (SetIfAbsentOrEqualNotStored) isSetIfAbsentOrEqualResponse() {}

// SetIfAbsentOrEqualStored indicates a successful set request where the value was stored.
type SetIfAbsentOrEqualStored struct{}

func (SetIfAbsentOrEqualStored) isSetIfAbsentOrEqualResponse() {}
