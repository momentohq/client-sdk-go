package responses

// SetRemoveElementResponse is the base response type for a set remobe element request.
type SetRemoveElementResponse interface {
	MomentoCacheResponse
	isSetRemoveElementResponse()
}

// SetRemoveElementSuccess indicates a successful set remove element request.
type SetRemoveElementSuccess struct{}

func (SetRemoveElementSuccess) isSetRemoveElementResponse() {}
