package responses

type DeleteResponse interface {
	isDeleteResponse()
}

type DeleteSuccess struct{}

func (DeleteSuccess) isDeleteResponse() {}
