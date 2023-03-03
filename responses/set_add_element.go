package responses

type SetAddElementResponse interface {
	isSetAddElementResponse()
}

type SetAddElementSuccess struct{}

func (SetAddElementSuccess) isSetAddElementResponse() {}
