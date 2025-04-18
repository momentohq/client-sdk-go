package momento

type DictionaryGetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}

func (r *DictionaryGetFieldRequest) cacheName() string { return r.CacheName }

func (c DictionaryGetFieldRequest) GetRequestName() string {
	return "DictionaryGetFieldRequest"
}
