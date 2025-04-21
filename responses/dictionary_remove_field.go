package responses

// DictionaryRemoveFieldResponse is the base response type for a dictionary remove field request.
type DictionaryRemoveFieldResponse interface {
	MomentoCacheResponse
	isDictionaryRemoveFieldResponse()
}

// DictionaryRemoveFieldSuccess indicates a successful dictionary remove field request.
type DictionaryRemoveFieldSuccess struct{}

func (DictionaryRemoveFieldSuccess) isDictionaryRemoveFieldResponse() {}
