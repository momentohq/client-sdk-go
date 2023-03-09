package responses

// DictionarySetFieldResponse is a base response type for a dictionary set field request.
type DictionarySetFieldResponse interface {
	isDictionarySetFieldResponse()
}

// DictionarySetFieldSuccess indicates a successful dictionary set field request.
type DictionarySetFieldSuccess struct{}

func (DictionarySetFieldSuccess) isDictionarySetFieldResponse() {}
