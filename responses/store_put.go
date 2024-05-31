package responses

type StorePutResponse interface {
	isStorePutResponse()
}

type StorePutSuccess struct{}

func (StorePutSuccess) isStorePutResponse() {}
