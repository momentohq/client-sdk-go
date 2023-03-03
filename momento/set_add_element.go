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
