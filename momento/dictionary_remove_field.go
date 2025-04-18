package momento

type DictionaryRemoveFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}

func (r *DictionaryRemoveFieldRequest) cacheName() string { return r.CacheName }

func (c DictionaryRemoveFieldRequest) GetRequestName() string {
	return "DictionaryRemoveFieldRequest"
}
