package models

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type DictionarySetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          []byte
	Value          []byte
	CollectionTtl  utils.CollectionTTL
}

type DictionaryGetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          []byte
}

type DictionarySetFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Items          map[string][]byte
	CollectionTTL  utils.CollectionTTL
}

type DictionaryIncrement struct {
	CacheName      string
	DictionaryName string
	Field          []byte
	CollectionTTL  utils.CollectionTTL
}

type DictionaryGetFields struct {
	CacheName      string
	DictionaryName string
	Fields         [][]byte
	CollectionTTL  utils.CollectionTTL
}

type DictionaryFetchRequest struct {
	CacheName      string
	DictionaryName string
}

type DictionaryDeleteRequest struct {
	CacheName      string
	DictionaryName string
}

type DictionaryRemoveField struct {
	CacheName      string
	DictionaryName string
	Field          []byte
}

type DictionaryRemoveFields struct {
	CacheName      string
	DictionaryName string
	Fields         [][]byte
}
