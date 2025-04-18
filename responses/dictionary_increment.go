package responses

// DictionaryIncrementResponse is the base response type for a dictionary increment request.
type DictionaryIncrementResponse interface {
	MomentoCacheResponse
	isDictionaryIncrementResponse()
}

// DictionaryIncrementSuccess indicates a successful dictionary increment success.
type DictionaryIncrementSuccess struct {
	value int64
}

func (DictionaryIncrementSuccess) isDictionaryIncrementResponse() {}

// Value returns the new value of the element after incrementing.
func (resp DictionaryIncrementSuccess) Value() int64 {
	return resp.value
}

// NewDictionaryIncrementSuccess returns a new DictionaryIncrementSuccess containing the supplied value.
func NewDictionaryIncrementSuccess(value int64) *DictionaryIncrementSuccess {
	return &DictionaryIncrementSuccess{value: value}
}
