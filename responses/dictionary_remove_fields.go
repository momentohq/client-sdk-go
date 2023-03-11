package responses

// DictionaryRemoveFieldsResponse is the base response type for a dictionary remove fields request.
type DictionaryRemoveFieldsResponse interface {
	isDictionaryRemoveFieldsResponse()
}

// DictionaryRemoveFieldsSuccess indicates a successful dictionary remove fields request.
type DictionaryRemoveFieldsSuccess struct{}

func (DictionaryRemoveFieldsSuccess) isDictionaryRemoveFieldsResponse() {}
