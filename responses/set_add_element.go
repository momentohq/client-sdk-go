package responses

// SetAddElementResponse is the base response type for a set add element request.
type SetAddElementResponse interface {
	MomentoCacheResponse
	isSetAddElementResponse()
}

// SetAddElementSuccess indicates a successful set add element request.
type SetAddElementSuccess struct{}

func (SetAddElementSuccess) isSetAddElementResponse() {}
