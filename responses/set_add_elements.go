package responses

// SetAddElementsResponse is a base response type for a set add elements request.
type SetAddElementsResponse interface {
	isSetAddElementResponse()
}

// SetAddElementsSuccess indicates a successful set add elements request.
type SetAddElementsSuccess struct{}

func (SetAddElementsSuccess) isSetAddElementResponse() {}
