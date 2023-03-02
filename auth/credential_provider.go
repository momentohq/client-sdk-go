package auth

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/golang-jwt/jwt/v4"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

type Endpoints struct {
	ControlEndpoint string
	CacheEndpoint   string
}

type CredentialProvider interface {
	GetAuthToken() string
	GetControlEndpoint() string
	GetCacheEndpoint() string
	WithEndpoints(endpoints *Endpoints) (CredentialProvider, error)
}

type defaultCredentialProvider struct {
	authToken       string
	controlEndpoint string
	cacheEndpoint   string
}

func (credentialProvider defaultCredentialProvider) GetAuthToken() string {
	return credentialProvider.authToken
}

func (credentialProvider defaultCredentialProvider) GetControlEndpoint() string {
	return credentialProvider.controlEndpoint
}

func (credentialProvider defaultCredentialProvider) GetCacheEndpoint() string {
	return credentialProvider.cacheEndpoint
}

func FromEnvironmentVariable(envVar string) (CredentialProvider, error) {
	credentialProvider, err := NewEnvMomentoTokenProvider(envVar)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}

func FromString(authToken string) (CredentialProvider, error) {
	credentialProvider, err := NewStringMomentoTokenProvider(authToken)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}

// WithEndpoints overrides the cache and control endpoint URIs with those provided by the supplied Endpoints struct
// and returns a CredentialProvider with the new endpoint values
func (credentialProvider defaultCredentialProvider) WithEndpoints(endpoints *Endpoints) (CredentialProvider, error) {
	if endpoints == nil {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"endpoints cannot be nil",
			errors.New("invalid argument"),
		)
	}
	credentialProvider.cacheEndpoint = endpoints.CacheEndpoint
	credentialProvider.controlEndpoint = endpoints.ControlEndpoint
	return credentialProvider, nil
}

// NewEnvMomentoTokenProvider constructor for a CredentialProvider using an environment variable to store an
// authentication token
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

// NewStringMomentoTokenProvider constructor for a CredentialProvider using a string containing an
// authentication token
func NewStringMomentoTokenProvider(authToken string) (CredentialProvider, error) {
	endpoints, err := getEndpointsFromToken(authToken)
	if err != nil {
		return nil, err
	}
	provider := defaultCredentialProvider{
		authToken:       authToken,
		controlEndpoint: endpoints.ControlEndpoint,
		cacheEndpoint:   endpoints.CacheEndpoint,
	}
	return provider, nil
}

func getEndpointsFromToken(authToken string) (*Endpoints, momentoerrors.MomentoSvcErr) {
	token, _, err := new(jwt.Parser).ParseUnverified(authToken, jwt.MapClaims{})
	if err != nil {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Could not parse auth token.",
			err,
		)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &Endpoints{
			ControlEndpoint: reflect.ValueOf(claims["cp"]).String(),
			CacheEndpoint:   reflect.ValueOf(claims["c"]).String(),
		}, nil
	}
	return nil, momentoerrors.NewMomentoSvcErr(
		momentoerrors.InvalidArgumentError,
		"Invalid Auth token.",
		nil,
	)
}
