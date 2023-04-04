package momento

type DictionaryGetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}

func (r *DictionaryGetFieldRequest) cacheName() string { return r.CacheName }
