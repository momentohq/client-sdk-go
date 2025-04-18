package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type SetAddElementRequest struct {
	CacheName string
	SetName   string
	Element   Value
	Ttl       *utils.CollectionTtl
}

func (r *SetAddElementRequest) cacheName() string { return r.CacheName }

func (c SetAddElementRequest) GetRequestName() string {
	return "SetAddElementRequest"
}
