package responses

type DictionarySetFieldsResponse interface {
	isDictionarySetFieldsResponse()
}

type DictionarySetFieldsSuccess struct{}

func (DictionarySetFieldsSuccess) isDictionarySetFieldsResponse() {}
