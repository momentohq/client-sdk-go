package responses

import "github.com/momentohq/client-sdk-go/responses"

// GenerateDisposableTokenResponse is the base response type for a generate disposable token request.
type GenerateDisposableTokenResponse interface {
	responses.MomentoAuthResponse
	isGenerateDisposableTokenResponse()
}

// GenerateDisposableTokenSuccess indicates a successful generate disposable token request.
type GenerateDisposableTokenSuccess struct {
	ApiKey     string
	Endpoint   string
	ValidUntil uint64
}

func (*GenerateDisposableTokenSuccess) isGenerateDisposableTokenResponse() {}
