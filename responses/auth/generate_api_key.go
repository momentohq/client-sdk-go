package responses

import "github.com/momentohq/client-sdk-go/utils"

// GenerateApiKeyResponse is the base response type for a generate api key request.
type GenerateApiKeyResponse interface {
	isGenerateApiKeyResponse()
}

// GenerateApiKeySuccess indicates a successful generate api key request.
type GenerateApiKeySuccess struct {
	ApiKey       string
	RefreshToken string
	Endpoint     string
	ExpiresAt    utils.ExpiresAt
}

func (*GenerateApiKeySuccess) isGenerateApiKeyResponse() {}
