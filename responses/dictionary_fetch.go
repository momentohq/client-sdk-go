package responses

type DictionaryFetchResponse interface {
	isDictionaryFetchResponse()
}

type DictionaryFetchHit struct {
	elementsStringByte   map[string][]byte
	elementsStringString map[string]string
}

func (DictionaryFetchHit) isDictionaryFetchResponse() {}

func (resp DictionaryFetchHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

func (resp DictionaryFetchHit) ValueMapStringString() map[string]string {
	if resp.elementsStringString == nil {
		resp.elementsStringString = make(map[string]string)
		for k, v := range resp.elementsStringByte {
			resp.elementsStringString[k] = string(v)
		}
	}
	return resp.elementsStringString
}

func (resp DictionaryFetchHit) ValueMapStringByte() map[string][]byte {
	return resp.elementsStringByte
}

type DictionaryFetchMiss struct{}

func (DictionaryFetchMiss) isDictionaryFetchResponse() {}

func NewDictionaryFetchHit(elements map[string][]byte) *DictionaryFetchHit {
	return &DictionaryFetchHit{elementsStringByte: elements}
}
