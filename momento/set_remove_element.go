package momento

// SetRemoveElementResponse

type SetRemoveElementResponse interface {
	isSetRemoveElementResponse()
}

type SetRemoveElementSuccess struct{}

func (SetRemoveElementSuccess) isSetRemoveElementResponse() {}

// SetRemoveElementRequest

type SetRemoveElementRequest struct {
	CacheName string
	SetName   string
	Element   Value
}
