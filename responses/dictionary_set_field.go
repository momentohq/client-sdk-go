package responses

// DictionarySetFieldResponse is the base response type for a dictionary set field request.
type DictionarySetFieldResponse interface {
	MomentoCacheResponse
	isDictionarySetFieldResponse()
}

// DictionarySetFieldSuccess indicates a successful dictionary set field request.
type DictionarySetFieldSuccess struct{}

func (DictionarySetFieldSuccess) isDictionarySetFieldResponse() {}
