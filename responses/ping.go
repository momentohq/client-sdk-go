package responses

type PingResponse interface {
	isPingResponse()
}

type PingSuccess struct{}

func (PingSuccess) isPingResponse() {}
