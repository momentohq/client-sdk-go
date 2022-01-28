package resolver

import (
	"reflect"

	"github.com/golang-jwt/jwt/v4"
	"github.com/momentohq/client-sdk-go/internal/requests"
)

const (
	MOMENTO_CONTROL_ENDPOINT_PREFIX = "control."
	MOMENTO_CACHE_ENDPOINT_PREFIX   = "cache."
	CONTROL_ENDPOINT_CLAIM_ID       = "cp"
	CACHE_ENDPOINT_CLAIM_ID         = "c"
)

type Endpoints struct {
	ControlEndpoint string
	CacheEndpoint   string
}

func Resolve(rr requests.ResolveRequest) (*Endpoints, error) {
	if rr.EndpointOverride != "" {
		return &Endpoints{ControlEndpoint: MOMENTO_CONTROL_ENDPOINT_PREFIX + rr.EndpointOverride, CacheEndpoint: MOMENTO_CACHE_ENDPOINT_PREFIX + rr.EndpointOverride}, nil
	}
	return getEndpointsFromToken(rr.AuthToken)
}

func getEndpointsFromToken(authToken string) (*Endpoints, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(authToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		ctEndpoint := reflect.ValueOf(claims[CONTROL_ENDPOINT_CLAIM_ID]).String()
		cEndpoint := reflect.ValueOf(claims[CACHE_ENDPOINT_CLAIM_ID]).String()
		return &Endpoints{ControlEndpoint: ctEndpoint, CacheEndpoint: cEndpoint}, nil
	} else {
		return nil, err
	}
}
