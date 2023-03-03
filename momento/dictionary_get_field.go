package momento

type DictionaryGetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}
