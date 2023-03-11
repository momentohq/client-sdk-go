package responses

// SetResponse is the base response type for a set request.
type SetResponse interface {
	isSetResponse()
}

// SetSuccess indicates a successful set request.
type SetSuccess struct{}

func (SetSuccess) isSetResponse() {}
