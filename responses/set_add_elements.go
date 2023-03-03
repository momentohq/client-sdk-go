package responses

type SetAddElementsResponse interface {
	isSetAddElementResponse()
}

type SetAddElementsSuccess struct{}

func (SetAddElementsSuccess) isSetAddElementResponse() {}
