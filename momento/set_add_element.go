package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

// SetAddElementResponse

type SetAddElementResponse interface {
	isSetAddElementResponse()
}

type SetAddElementSuccess struct{}

func (SetAddElementSuccess) isSetAddElementResponse() {}

// SetAddElementRequest

type SetAddElementRequest struct {
	CacheName     string
	SetName       string
	Element       Value
	CollectionTtl utils.CollectionTtl
}
