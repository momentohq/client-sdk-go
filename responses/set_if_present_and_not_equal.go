package responses

// SetIfPresentAndNotEqualResponse is the base response type for a SetIfPresentAndNotEqual request
type SetIfPresentAndNotEqualResponse interface {
	MomentoCacheResponse
	isSetIfPresentAndNotEqualResponse()
}

// SetIfPresentAndNotEqualNotStored indicates a successful set request where the value was already present.
type SetIfPresentAndNotEqualNotStored struct{}

func (SetIfPresentAndNotEqualNotStored) isSetIfPresentAndNotEqualResponse() {}

// SetIfPresentAndNotEqualStored indicates a successful set request where the value was stored.
type SetIfPresentAndNotEqualStored struct{}

func (SetIfPresentAndNotEqualStored) isSetIfPresentAndNotEqualResponse() {}
