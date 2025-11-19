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

type Endpoint struct {
	// Endpoint is the host which the Momento client will connect to
	Endpoint string
	// InsecureConnection is a flag to indicate whether the connection to the endpoint should be insecure. The zero value for a bool in Go is false, so we default to a secure connection (InsecureConnection==false) if a value is not provided.
	InsecureConnection bool
}

func (endpoint Endpoint) String() string {
	return fmt.Sprintf(
		"Endpoint{Endpoint=%s, InsecureConnection=%t}",
		endpoint.Endpoint,
		endpoint.InsecureConnection)
}

type AllEndpoints struct {
	// ControlEndpoint is the host which the Momento client will connect to the Momento control plane
	ControlEndpoint Endpoint
	// CacheEndpoint is the host which the Momento client will connect to the Momento data plane
	CacheEndpoint Endpoint
	// TokenEndpoint is the host which the Momento client will connect to for generating disposable auth tokens
	TokenEndpoint Endpoint
	// StorageEndpoint is the host which the Momento client will connect to the Momento storage data plane
	StorageEndpoint Endpoint
}

type tokenAndEndpoints struct {
	Endpoints AllEndpoints
	AuthToken string
}

type CredentialProvider interface {
	GetAuthToken() string
	GetCacheTlsHostname() string
	GetControlEndpoint() string
	IsControlEndpointSecure() bool
	GetCacheEndpoint() string
	IsCacheEndpointSecure() bool
	GetTokenEndpoint() string
	IsTokenEndpointSecure() bool
	GetStorageEndpoint() string
	IsStorageEndpointSecure() bool
	WithEndpoints(endpoints AllEndpoints) (CredentialProvider, error)
}

type defaultCredentialProvider struct {
	authToken       string
	controlEndpoint Endpoint
	cacheEndpoint   Endpoint
	tokenEndpoint   Endpoint
	storageEndpoint Endpoint
	// TLS hostname we want to preserve even if we override the cache endpoint
	cacheTlsHostname string
}

// GetAuthToken returns user's auth token.
func (credentialProvider defaultCredentialProvider) GetAuthToken() string {
	return credentialProvider.authToken
}

// GetCacheTlsHostname returns the TLS hostname for the cache endpoint.
func (credentialProvider defaultCredentialProvider) GetCacheTlsHostname() string {
	return credentialProvider.cacheTlsHostname
}

// GetControlEndpoint returns AllEndpoints.ControlEndpoint.Endpoint.
func (credentialProvider defaultCredentialProvider) GetControlEndpoint() string {
	return credentialProvider.controlEndpoint.Endpoint
}

// IsControlEndpointSecure returns true if the control endpoint is secure.
func (credentialProvider defaultCredentialProvider) IsControlEndpointSecure() bool {
	return !credentialProvider.controlEndpoint.InsecureConnection
}

// GetCacheEndpoint returns AllEndpoints.CacheEndpoint.Endpoint.
func (credentialProvider defaultCredentialProvider) GetCacheEndpoint() string {
	return credentialProvider.cacheEndpoint.Endpoint
}

// IsCacheEndpointSecure returns true if the cace endpoint is secure.
func (credentialProvider defaultCredentialProvider) IsCacheEndpointSecure() bool {
	return !credentialProvider.cacheEndpoint.InsecureConnection
}

// GetTokenEndpoint returns AllEndpoints.TokenEndpoint.Endpoint.
func (credentialProvider defaultCredentialProvider) GetTokenEndpoint() string {
	return credentialProvider.tokenEndpoint.Endpoint
}

// IsTokenEndpointSecure returns true if the token endpoint is secure.
func (credentialProvider defaultCredentialProvider) IsTokenEndpointSecure() bool {
	return !credentialProvider.tokenEndpoint.InsecureConnection
}

// GetStorageEndpoint returns AllEndpoints.StorageEndpoint.Endpoint.
func (credentialProvider defaultCredentialProvider) GetStorageEndpoint() string {
	return credentialProvider.storageEndpoint.Endpoint
}

// IsStorageEndpointSecure returns true if the storage endpoint is secure.
func (credentialProvider defaultCredentialProvider) IsStorageEndpointSecure() bool {
	return !credentialProvider.storageEndpoint.InsecureConnection
}

func (credentialProvider defaultCredentialProvider) String() string {
	authToken := credentialProvider.authToken
	if len(authToken) > 4 {
		authToken = fmt.Sprintf("%s***%s", authToken[:2], authToken[len(authToken)-2:])
	} else {
		authToken = "***"
	}

	return fmt.Sprintf(
		"CredentialProvider{authToken=%s, controlEndpoint=%s, cacheEndpoint=%s, tokenEndpoint=%s, storageEndpoint=%s}",
		authToken,
		credentialProvider.controlEndpoint,
		credentialProvider.cacheEndpoint,
		credentialProvider.tokenEndpoint,
		credentialProvider.storageEndpoint,
	)
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
func (credentialProvider defaultCredentialProvider) WithEndpoints(endpoints AllEndpoints) (CredentialProvider, error) {
	if endpoints.CacheEndpoint.Endpoint != "" {
		credentialProvider.cacheEndpoint = endpoints.CacheEndpoint
	}
	if endpoints.ControlEndpoint.Endpoint != "" {
		credentialProvider.controlEndpoint = endpoints.ControlEndpoint
	}
	if endpoints.TokenEndpoint.Endpoint != "" {
		credentialProvider.tokenEndpoint = endpoints.TokenEndpoint
	}
	if endpoints.StorageEndpoint.Endpoint != "" {
		credentialProvider.storageEndpoint = endpoints.StorageEndpoint
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
	return NewStringMomentoTokenProvider(authToken)
}

// NewStringMomentoTokenProvider constructor for a CredentialProvider using a string containing an
// authentication token
func NewStringMomentoTokenProvider(authToken string) (CredentialProvider, error) {
	tokenAndEndpoints, err := decodeAuthToken(authToken)
	if err != nil {
		return nil, err
	}
	port := 443
	provider := defaultCredentialProvider{
		authToken:        tokenAndEndpoints.AuthToken,
		cacheTlsHostname: tokenAndEndpoints.Endpoints.CacheEndpoint.Endpoint,
		controlEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("%s:%d", tokenAndEndpoints.Endpoints.ControlEndpoint.Endpoint, port),
		},
		cacheEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("%s:%d", tokenAndEndpoints.Endpoints.CacheEndpoint.Endpoint, port),
		},
		tokenEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("%s:%d", tokenAndEndpoints.Endpoints.TokenEndpoint.Endpoint, port),
		},
		storageEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("%s:%d", tokenAndEndpoints.Endpoints.StorageEndpoint.Endpoint, port),
		},
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
		Endpoints: AllEndpoints{
			ControlEndpoint: Endpoint{Endpoint: fmt.Sprintf("control.%s", tokenData["endpoint"])},
			CacheEndpoint:   Endpoint{Endpoint: fmt.Sprintf("cache.%s", tokenData["endpoint"])},
			TokenEndpoint:   Endpoint{Endpoint: fmt.Sprintf("token.%s", tokenData["endpoint"])},
			StorageEndpoint: Endpoint{Endpoint: fmt.Sprintf("storage.%s", tokenData["endpoint"])},
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
			Endpoints: AllEndpoints{
				ControlEndpoint: Endpoint{Endpoint: controlEndpoint},
				CacheEndpoint:   Endpoint{Endpoint: cacheEndpoint},
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

type MomentoLocalConfig struct {
	Hostname string
	Port     uint
}

func NewMomentoLocalProvider(config *MomentoLocalConfig) (CredentialProvider, error) {
	hostname := "127.0.0.1"
	port := uint(8080)
	if config != nil {
		if config.Hostname != "" {
			hostname = config.Hostname
		}
		if config.Port != 0 {
			port = config.Port
		}
	}

	momentoLocalEndpoint := Endpoint{
		Endpoint:           fmt.Sprintf("%s:%d", hostname, port),
		InsecureConnection: true,
	}
	tokenAndEndpoints := &tokenAndEndpoints{
		Endpoints: AllEndpoints{
			ControlEndpoint: momentoLocalEndpoint,
			CacheEndpoint:   momentoLocalEndpoint,
			TokenEndpoint:   momentoLocalEndpoint,
			StorageEndpoint: momentoLocalEndpoint,
		},
	}

	provider := defaultCredentialProvider{
		authToken:       tokenAndEndpoints.AuthToken,
		controlEndpoint: tokenAndEndpoints.Endpoints.ControlEndpoint,
		cacheEndpoint:   tokenAndEndpoints.Endpoints.CacheEndpoint,
		tokenEndpoint:   tokenAndEndpoints.Endpoints.TokenEndpoint,
		storageEndpoint: tokenAndEndpoints.Endpoints.StorageEndpoint,
	}
	return provider, nil
}

type GlobalKeyFromStringProps struct {
	ApiKey string
	Endpoint  string
}

type GlobalKeyFromEnvVarProps struct {
	EnvVarName string
	Endpoint   string
}

// NewGlobalEnvMomentoTokenProvider constructor for a CredentialProvider using an endpoint and an environment
// variable to store a global api key.
func NewGlobalEnvMomentoTokenProvider(props GlobalKeyFromEnvVarProps) (CredentialProvider, error) {
	if props.EnvVarName == "" {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Environment variable name is empty",
			errors.New("invalid argument"),
		)
	}
	var authToken = os.Getenv(props.EnvVarName)
	if authToken == "" {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			fmt.Sprintf("Missing required environment variable %s", props.EnvVarName),
			errors.New("invalid argument"),
		)
	}
	return NewGlobalStringMomentoTokenProvider(GlobalKeyFromStringProps{
		ApiKey:  authToken,
		Endpoint: props.Endpoint,
	})
}

// NewGlobalStringMomentoTokenProvider constructor for a CredentialProvider using an endpoint and a string
// containing a global api key.
func NewGlobalStringMomentoTokenProvider(props GlobalKeyFromStringProps) (CredentialProvider, error) {
	if props.ApiKey == "" {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Auth token is an empty string",
			errors.New("invalid argument"),
		)
	}
	if props.Endpoint == "" {
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"Endpoint is an empty string",
			errors.New("invalid argument"),
		)
	}
	port := 443
	provider := defaultCredentialProvider{
		authToken:        props.ApiKey,
		cacheTlsHostname: props.Endpoint,
		controlEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("control.%s:%d", props.Endpoint, port),
		},
		cacheEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("cache.%s:%d", props.Endpoint, port),
		},
		tokenEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("token.%s:%d", props.Endpoint, port),
		},
		storageEndpoint: Endpoint{
			Endpoint: fmt.Sprintf("storage.%s:%d", props.Endpoint, port),
		},
	}
	return provider, nil
}

// GlobalKeyFromEnvironmentVariable returns a new CredentialProvider using a global api key stored
// in the provided environment variable, as well as an endpoint.
func GlobalKeyFromEnvironmentVariable(props GlobalKeyFromEnvVarProps) (CredentialProvider, error) {
	credentialProvider, err := NewGlobalEnvMomentoTokenProvider(props)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}

// GlobalKeyFromString returns a new CredentialProvider with the provided global api key and endpoint.
func GlobalKeyFromString(props GlobalKeyFromStringProps) (CredentialProvider, error) {
	credentialProvider, err := NewGlobalStringMomentoTokenProvider(props)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}
