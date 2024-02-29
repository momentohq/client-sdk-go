package responses

// SetIfAbsentResponse is the base response type for a SetIfAbsent request
type SetIfAbsentResponse interface {
	isSetIfAbsentResponse()
}

// SetIfAbsentNotStored indicates a successful set request where the value was already present.
type SetIfAbsentNotStored struct{}

func (SetIfAbsentNotStored) isSetIfAbsentResponse() {}

// SetIfAbsentStored indicates a successful set request where the value was stored.
type SetIfAbsentStored struct{}

func (SetIfAbsentStored) isSetIfAbsentResponse() {}
