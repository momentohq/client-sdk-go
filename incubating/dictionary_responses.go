package incubating

type DictionaryFetchResponse interface {
	isDictionaryFetchResponse()
}

type DictionaryFetchMiss struct{}

func (DictionaryFetchMiss) isDictionaryFetchResponse() {}

type DictionaryFetchHit struct {
	items             map[string][]byte
	itemsStringString map[string]string
}

func (DictionaryFetchHit) isDictionaryFetchResponse() {}

func (resp DictionaryFetchHit) ValueDictionaryStringByteArray() map[string][]byte {
	return resp.items
}

func (resp DictionaryFetchHit) ValueDictionaryStringString() map[string]string {
	if resp.itemsStringString == nil {
		for k, v := range resp.items {
			resp.itemsStringString[k] = string(v)
		}
	}
	return resp.itemsStringString
}

type DictionaryGetFieldResponse interface {
	isDictionaryGetFieldResponse()
}

type DictionaryGetFieldMiss struct{}

func (DictionaryGetFieldMiss) isDictionaryGetFieldResponse() {}

type DictionaryGetFieldHit struct {
	value []byte
}

func (DictionaryGetFieldHit) isDictionaryGetFieldResponse() {}

func (resp DictionaryGetFieldHit) ValueByteArray() []byte {
	return resp.value
}

func (resp DictionaryGetFieldHit) ValueString() string {
	return string(resp.value)
}
