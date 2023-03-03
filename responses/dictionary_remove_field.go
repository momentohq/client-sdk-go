package responses

type DictionaryRemoveFieldResponse interface {
	isDictionaryRemoveFieldResponse()
}

type DictionaryRemoveFieldSuccess struct{}

func (DictionaryRemoveFieldSuccess) isDictionaryRemoveFieldResponse() {}
