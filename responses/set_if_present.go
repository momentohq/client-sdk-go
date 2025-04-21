package responses

// SetIfPresentResponse is the base response type for a SetIfPresent request
type SetIfPresentResponse interface {
	MomentoCacheResponse
	isSetIfPresentResponse()
}

// SetIfPresentNotStored indicates a successful set request where the value was already present.
type SetIfPresentNotStored struct{}

func (SetIfPresentNotStored) isSetIfPresentResponse() {}

// SetIfPresentStored indicates a successful set request where the value was stored.
type SetIfPresentStored struct{}

func (SetIfPresentStored) isSetIfPresentResponse() {}
