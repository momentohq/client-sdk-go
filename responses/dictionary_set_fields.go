package responses

// DictionarySetFieldsResponse is a base response type for a dictionary set fields request.
type DictionarySetFieldsResponse interface {
	isDictionarySetFieldsResponse()
}

// DictionarySetFieldsSuccess indicates a successful dictionary set fields request.
type DictionarySetFieldsSuccess struct{}

func (DictionarySetFieldsSuccess) isDictionarySetFieldsResponse() {}
