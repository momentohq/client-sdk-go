package momento

type SortedSetRemoveElementRequest struct {
	CacheName string
	SetName   string
	Value     Value
}

func (r *SortedSetRemoveElementRequest) cacheName() string { return r.CacheName }

func (c SortedSetRemoveElementRequest) GetRequestName() string {
	return "SortedSetRemoveElementRequest"
}
