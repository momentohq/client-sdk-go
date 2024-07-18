package momento

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type GenerateApiKeyRequest struct {
	ExpiresIn utils.ExpiresIn
	Scope     PermissionScope
}
