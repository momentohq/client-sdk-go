package responses

// SetRemoveElementResponse is a base response type for a set remobe element request.
type SetRemoveElementResponse interface {
	isSetRemoveElementResponse()
}

// SetRemoveElementSuccess indicates a successful set remove element request.
type SetRemoveElementSuccess struct{}

func (SetRemoveElementSuccess) isSetRemoveElementResponse() {}
