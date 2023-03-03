package responses

type DictionarySetFieldResponse interface {
	isDictionarySetFieldResponse()
}

type DictionarySetFieldSuccess struct{}

func (DictionarySetFieldSuccess) isDictionarySetFieldResponse() {}
