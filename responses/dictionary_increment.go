package responses

type DictionaryIncrementResponse interface {
	isDictionaryIncrementResponse()
}

type DictionaryIncrementSuccess struct {
	value int64
}

func (DictionaryIncrementSuccess) isDictionaryIncrementResponse() {}

func (resp DictionaryIncrementSuccess) Value() int64 {
	return resp.value
}

func NewDictionaryIncrementSuccess(value int64) *DictionaryIncrementSuccess {
	return &DictionaryIncrementSuccess{value: value}
}
