package responses

type SetRemoveElementResponse interface {
	isSetRemoveElementResponse()
}

type SetRemoveElementSuccess struct{}

func (SetRemoveElementSuccess) isSetRemoveElementResponse() {}
