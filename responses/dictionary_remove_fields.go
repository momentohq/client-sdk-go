package responses

type DictionaryRemoveFieldsResponse interface {
	isDictionaryRemoveFieldsResponse()
}

type DictionaryRemoveFieldsSuccess struct{}

func (DictionaryRemoveFieldsSuccess) isDictionaryRemoveFieldsResponse() {}
