package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

// DictionarySetFieldResponse

type DictionarySetFieldResponse interface {
	isDictionarySetFieldResponse()
}

type DictionarySetFieldSuccess struct{}

func (DictionarySetFieldSuccess) isDictionarySetFieldResponse() {}

// DictionarySetFieldRequest

type DictionarySetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
	Value          Value
	CollectionTTL  utils.CollectionTTL
}
