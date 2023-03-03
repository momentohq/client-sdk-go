package responses

type SetResponse interface {
	isSetResponse()
}

type SetSuccess struct{}

func (SetSuccess) isSetResponse() {}
