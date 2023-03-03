package responses

type SetRemoveElementsResponse interface {
	isSetRemoveElementsResponse()
}

type SetRemoveElementsSuccess struct{}

func (SetRemoveElementsSuccess) isSetRemoveElementsResponse() {}
