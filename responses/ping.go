package responses

// PingResponse is the base response type for a ping request.
type PingResponse interface {
	MomentoCacheResponse
	isPingResponse()
}

// PingSuccess indicates a successful ping request.
type PingSuccess struct{}

func (PingSuccess) isPingResponse() {}
