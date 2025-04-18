package momento

type SetRemoveElementRequest struct {
	CacheName string
	SetName   string
	Element   Value
}

func (r *SetRemoveElementRequest) cacheName() string { return r.CacheName }

func (c SetRemoveElementRequest) GetRequestName() string {
	return "SetRemoveElementRequest"
}
