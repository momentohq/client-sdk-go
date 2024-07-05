// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
)

type AuthClient interface {
	GenerateDisposableToken(ctx context.Context, request *GenerateDisposableTokenRequest) (responses.GenerateDisposableTokenResponse, error)

	Close()
}

// defaultAuthClient represents all information needed for momento client to enable api calls to our auth endpoints.
type defaultAuthClient struct {
	credentialProvider auth.CredentialProvider
	tokenClient        *tokenClient
	log                logger.MomentoLogger
}

// NewAuthClient returns a new AuthClient with provided configuration and credential provider arguments.
func NewAuthClient(authConfiguration config.AuthConfiguration, credentialProvider auth.CredentialProvider) (AuthClient, error) {
	client := &defaultAuthClient{
		credentialProvider: credentialProvider,
		log:                authConfiguration.GetLoggerFactory().GetLogger("auth-client"),
	}

	tokenClient, err := newTokenClient(&models.TokenClientRequest{
		CredentialProvider: credentialProvider,
		Log:                authConfiguration.GetLoggerFactory().GetLogger("token-client"),
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.tokenClient = tokenClient

	return client, nil
}

func (c defaultAuthClient) GenerateDisposableToken(ctx context.Context, request *GenerateDisposableTokenRequest) (responses.GenerateDisposableTokenResponse, error) {
	if err := utils.ValidateDisposableTokenExpiry(request.ExpiresIn); err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	requestMetadata := internal.CreateMetadata(ctx, internal.Auth)

	tokenResp, err := c.tokenClient.GenerateDisposableToken(requestMetadata, request)

	if err != nil {
		c.log.Debug("failed to generate disposable token...")
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	switch result := tokenResp.(type) {
	case *responses.GenerateDisposableTokenSuccess:
		return &responses.GenerateDisposableTokenSuccess{
			ApiKey:     result.ApiKey,
			Endpoint:   result.Endpoint,
			ValidUntil: result.ValidUntil,
		}, nil
	}
	return nil, convertMomentoSvcErrorToCustomerError(
		momentoerrors.NewMomentoSvcErr(
			momentoerrors.UnknownServiceError, "Unknown service error was returned when requesting a disposable token", nil,
		),
	)
}

func (c defaultAuthClient) Close() {
	defer c.tokenClient.close()
}
