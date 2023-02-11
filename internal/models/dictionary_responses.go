package models

type DictionaryFetchResponse interface {
	isDictionaryFetchResponse()
}

type DictionaryFetchMiss struct{}

func (DictionaryFetchMiss) isDictionaryFetchResponse() {}

type DictionaryFetchHit struct {
	Items map[string][]byte
}

func (DictionaryFetchHit) isDictionaryFetchResponse() {}

type DictionaryGetFieldResponse interface {
	isDictionaryGetFieldResponse()
}

type DictionaryGetFieldMiss struct{}

func (DictionaryGetFieldMiss) isDictionaryGetFieldResponse() {}

type DictionaryGetFieldHit struct {
	Value []byte
}

func (DictionaryGetFieldHit) isDictionaryGetFieldResponse() {}
