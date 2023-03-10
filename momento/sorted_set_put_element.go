package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type SortedSetPutElementRequest struct {
	CacheName string
	SetName   string
	Value     Value
	Score     float64
	Ttl       *utils.CollectionTtl
}
