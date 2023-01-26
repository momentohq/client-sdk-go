package auth

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/golang-jwt/jwt/v4"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

type ResolveRequest struct {
	AuthToken        string
	EndpointOverride string
}

type Endpoints struct {
	ControlEndpoint string
	CacheEndpoint   string
}

type CredentialProvider interface {
	GetAuthToken() string
	GetControlEndpoint() string
	GetCacheEndpoint() string
}

type DefaultCredentialProvider struct {
	authToken       string
	controlEndpoint string
	cacheEndpoint   string
}

func (credentialProvider DefaultCredentialProvider) GetAuthToken() string {
	return credentialProvider.authToken
}

func (credentialProvider DefaultCredentialProvider) GetControlEndpoint() string {
	return credentialProvider.controlEndpoint
}

func (credentialProvider DefaultCredentialProvider) GetCacheEndpoint() string {
	return credentialProvider.cacheEndpoint
}

// NewEnvMomentoTokenProvider
// TODO: add overrides for endpoints
func NewEnvMomentoTokenProvider(envVariableName string) (CredentialProvider, error) {
	var authToken = os.Getenv(envVariableName)
	if authToken == "" {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			fmt.Sprintf("Missing required environment variable %s", envVariableName),
			errors.New("invalid argument"),
		)
	}
	return NewStringMomentoTokenProvider(authToken)
}

// NewStringMomentoTokenProvider
// TODO: add overrides for endpoints
func NewStringMomentoTokenProvider(authToken string) (CredentialProvider, error) {
	endpoints, err := resolve(&ResolveRequest{
		AuthToken: authToken,
	})
	if err != nil {
		return nil, err
	}
	provider := DefaultCredentialProvider{
		authToken:       authToken,
		controlEndpoint: endpoints.ControlEndpoint,
		cacheEndpoint:   endpoints.CacheEndpoint,
	}
	return provider, nil
}

const (
	momentoControlEndpointPrefix = "control."
	momentoCacheEndpointPrefix   = "cache."
	controlEndpointClaimId       = "cp"
	cacheEndpointClaimId         = "c"
)

func resolve(request *ResolveRequest) (*Endpoints, momentoerrors.MomentoSvcErr) {
	if request.EndpointOverride != "" {
		return &Endpoints{
			ControlEndpoint: momentoControlEndpointPrefix + request.EndpointOverride,
			CacheEndpoint:   momentoCacheEndpointPrefix + request.EndpointOverride,
		}, nil
	}
	return getEndpointsFromToken(request.AuthToken)
}

func getEndpointsFromToken(authToken string) (*Endpoints, momentoerrors.MomentoSvcErr) {
	token, _, err := new(jwt.Parser).ParseUnverified(authToken, jwt.MapClaims{})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &Endpoints{
			ControlEndpoint: reflect.ValueOf(claims[controlEndpointClaimId]).String(),
			CacheEndpoint:   reflect.ValueOf(claims[cacheEndpointClaimId]).String(),
		}, nil
	}
	return nil, momentoerrors.NewMomentoSvcErr(
		momentoerrors.InvalidArgumentError,
		"Invalid Auth token.",
		nil,
	)
}
