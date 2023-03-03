package momento

type DictionaryRemoveFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}
