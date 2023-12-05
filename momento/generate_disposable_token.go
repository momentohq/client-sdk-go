package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type GenerateDisposableTokenRequest struct {
	ExpiresIn utils.ExpiresIn
	Scope     DisposableTokenScope
	Props     DisposableTokenProps
}
