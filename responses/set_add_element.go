package responses

// SetAddElementResponse is a base response type for a set add element request.
type SetAddElementResponse interface {
	isSetAddElementResponse()
}

// SetAddElementSuccess indicates a successful set add element request.
type SetAddElementSuccess struct{}

func (SetAddElementSuccess) isSetAddElementResponse() {}
