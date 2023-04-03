package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type DictionarySetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
	Value          Value
	Ttl            *utils.CollectionTtl
}

func (r *DictionarySetFieldRequest) cacheName() string { return r.CacheName }
