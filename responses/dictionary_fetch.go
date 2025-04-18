package responses

// DictionaryFetchResponse is the base response type for a dictionary fetch request.
type DictionaryFetchResponse interface {
	MomentoCacheResponse
	isDictionaryFetchResponse()
}

// DictionaryFetchHit indicates a hit dictionary fetch request.
type DictionaryFetchHit struct {
	elementsStringByte   map[string][]byte
	elementsStringString map[string]string
}

func (DictionaryFetchHit) isDictionaryFetchResponse() {}

// ValueMap returns the data as a Map whose keys and values are utf-8 strings, decoded from the underlying byte arrays.
// This is a convenience alias for ValueMapStringString.
func (resp DictionaryFetchHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

// ValueMapStringString returns the data as a Map whose keys and values are utf-8 strings, decoded from the underlying byte arrays.
func (resp DictionaryFetchHit) ValueMapStringString() map[string]string {
	if resp.elementsStringString == nil {
		resp.elementsStringString = make(map[string]string)
		for k, v := range resp.elementsStringByte {
			resp.elementsStringString[k] = string(v)
		}
	}
	return resp.elementsStringString
}

// ValueMapStringByte returns the data as a Map whose keys are utf-8 strings, decoded from the underlying byte array, and whose values are byte arrays.
func (resp DictionaryFetchHit) ValueMapStringByte() map[string][]byte {
	return resp.elementsStringByte
}

// DictionaryFetchMiss indicates a dictionary fetch request was a miss.
type DictionaryFetchMiss struct{}

func (DictionaryFetchMiss) isDictionaryFetchResponse() {}

// NewDictionaryFetchHit returns a new DictionaryFetchHit with the data as a Map whose keys are utf-8 strings.
func NewDictionaryFetchHit(elements map[string][]byte) *DictionaryFetchHit {
	return &DictionaryFetchHit{elementsStringByte: elements}
}
