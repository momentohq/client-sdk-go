package responses

// DictionaryRemoveFieldResponse is a base response type for a dictionary remove field request.
type DictionaryRemoveFieldResponse interface {
	isDictionaryRemoveFieldResponse()
}

// DictionaryRemoveFieldSuccess indicates a successful dictionary remove field request.
type DictionaryRemoveFieldSuccess struct{}

func (DictionaryRemoveFieldSuccess) isDictionaryRemoveFieldResponse() {}
