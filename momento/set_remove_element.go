package momento

type SetRemoveElementRequest struct {
	CacheName string
	SetName   string
	Element   Value
}
