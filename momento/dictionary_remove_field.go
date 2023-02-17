package momento

// DictionaryRemoveFieldResponse

type DictionaryRemoveFieldResponse interface {
	isDictionaryRemoveFieldResponse()
}

type DictionaryRemoveFieldSuccess struct{}

func (DictionaryRemoveFieldSuccess) isDictionaryRemoveFieldResponse() {}

// DictionaryRemoveFieldRequest

type DictionaryRemoveFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}
