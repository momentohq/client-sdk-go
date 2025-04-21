package responses

import (
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
)

// RefreshApiKeyResponse is the base response type for a refresh api key request.
type RefreshApiKeyResponse interface {
	responses.MomentoAuthResponse
	isRefreshApiKeyResponse()
}

// RefreshApiKeySuccess indicates a successful refresh api key request.
type RefreshApiKeySuccess struct {
	ApiKey       string
	RefreshToken string
	Endpoint     string
	ExpiresAt    utils.ExpiresAt
}

func (*RefreshApiKeySuccess) isRefreshApiKeyResponse() {}
