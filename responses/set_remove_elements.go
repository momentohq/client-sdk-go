package responses

// SetRemoveElementsResponse is the base response type for a set remove elements request.
type SetRemoveElementsResponse interface {
	isSetRemoveElementsResponse()
}

// SetRemoveElementsSuccess indicates a successful set remove elements request.
type SetRemoveElementsSuccess struct{}

func (SetRemoveElementsSuccess) isSetRemoveElementsResponse() {}
