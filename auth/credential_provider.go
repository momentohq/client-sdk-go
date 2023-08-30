package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/golang-jwt/jwt/v4"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

type Endpoints struct {
	// ControlEndpoint is the host which the Momento client will connect to the Momento control plane
	ControlEndpoint string
	// CacheEndpoint is the host which the Momento client will connect to the Momento data plane
	CacheEndpoint string
}

type tokenAndEndpoints struct {
	Endpoints
	AuthToken string
}

type CredentialProvider interface {
	GetAuthToken() string
	GetControlEndpoint() string
	GetCacheEndpoint() string
	WithEndpoints(endpoints Endpoints) (CredentialProvider, error)
}

type defaultCredentialProvider struct {
	authToken       string
	controlEndpoint string
	cacheEndpoint   string
}

// GetAuthToken returns user's auth token.
func (credentialProvider defaultCredentialProvider) GetAuthToken() string {
	return credentialProvider.authToken
}

// GetControlEndpoint returns Endpoints.ControlEndpoint.
func (credentialProvider defaultCredentialProvider) GetControlEndpoint() string {
	return credentialProvider.controlEndpoint
}

// GetCacheEndpoint returns Endpoints.CacheEndpoint.
func (credentialProvider defaultCredentialProvider) GetCacheEndpoint() string {
	return credentialProvider.cacheEndpoint
}

// FromEnvironmentVariable returns a new CredentialProvider using an auth token stored in the provided environment variable.
func FromEnvironmentVariable(envVar string) (CredentialProvider, error) {
	credentialProvider, err := NewEnvMomentoTokenProvider(envVar)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}

// FromString returns a new CredentialProvider with the provided user auth token.
func FromString(authToken string) (CredentialProvider, error) {
	credentialProvider, err := NewStringMomentoTokenProvider(authToken)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}

// WithEndpoints overrides the cache and control endpoint URIs with those provided by the supplied Endpoints struct
// and returns a CredentialProvider with the new endpoint values. An endpoint supplied as an empty string is ignored
// and the existing value for that endpoint is retained.
func (credentialProvider defaultCredentialProvider) WithEndpoints(endpoints Endpoints) (CredentialProvider, error) {
	if endpoints.CacheEndpoint != "" {
		credentialProvider.cacheEndpoint = endpoints.CacheEndpoint
	}
	if endpoints.ControlEndpoint != "" {
		credentialProvider.controlEndpoint = endpoints.ControlEndpoint
	}
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
	print(authToken)
	return NewStringMomentoTokenProvider(authToken)
}

// NewStringMomentoTokenProvider constructor for a CredentialProvider using a string containing an
// authentication token
func NewStringMomentoTokenProvider(authToken string) (CredentialProvider, error) {
	tokenAndEndpoints, err := decodeAuthToken(authToken)
	if err != nil {
		return nil, err
	}
	provider := defaultCredentialProvider{
		authToken:       tokenAndEndpoints.AuthToken,
		controlEndpoint: tokenAndEndpoints.ControlEndpoint,
		cacheEndpoint:   tokenAndEndpoints.CacheEndpoint,
	}
	return provider, nil
}

func decodeAuthToken(authToken string) (*tokenAndEndpoints, momentoerrors.MomentoSvcErr) {
	decodedBase64Token, err := b64.StdEncoding.DecodeString(authToken)
	if err != nil {
		return processJwtToken(authToken)
	}
	return processV1Token(decodedBase64Token)
}

func processV1Token(decodedBase64Token []byte) (*tokenAndEndpoints, momentoerrors.MomentoSvcErr) {
	var tokenData map[string]string
	if err := json.Unmarshal(decodedBase64Token, &tokenData); err != nil {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"malformed auth token",
			nil,
		)
	}

	if tokenData["endpoint"] == "" || tokenData["api_key"] == "" {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"failed to parse token",
			nil,
		)
	}

	return &tokenAndEndpoints{
		Endpoints: Endpoints{
			ControlEndpoint: fmt.Sprintf("control.%s", tokenData["endpoint"]),
			CacheEndpoint:   fmt.Sprintf("cache.%s", tokenData["endpoint"]),
		},
		AuthToken: tokenData["api_key"],
	}, nil
}

func processJwtToken(authToken string) (*tokenAndEndpoints, momentoerrors.MomentoSvcErr) {
	token, _, err := new(jwt.Parser).ParseUnverified(authToken, jwt.MapClaims{})
	if err != nil {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Could not parse auth token.",
			err,
		)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		controlEndpoint := reflect.ValueOf(claims["cp"]).String()
		cacheEndpoint := reflect.ValueOf(claims["c"]).String()
		return &tokenAndEndpoints{
			Endpoints: Endpoints{
				ControlEndpoint: controlEndpoint,
				CacheEndpoint:   cacheEndpoint,
			},
			AuthToken: authToken,
		}, nil
	}
	return nil, momentoerrors.NewMomentoSvcErr(
		momentoerrors.InvalidArgumentError,
		"Invalid Auth token.",
		nil,
	)
}
