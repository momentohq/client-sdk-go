package resolver

import (
	"reflect"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"github.com/golang-jwt/jwt/v4"
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

func Resolve(request *models.ResolveRequest) (*Endpoints, momentoerrors.MomentoBaseError) {
	if request.EndpointOverride != "" {
		return &Endpoints{ControlEndpoint: MOMENTO_CONTROL_ENDPOINT_PREFIX + request.EndpointOverride, CacheEndpoint: MOMENTO_CACHE_ENDPOINT_PREFIX + request.EndpointOverride}, nil
	}
	return getEndpointsFromToken(request.AuthToken)
}

func getEndpointsFromToken(authToken string) (*Endpoints, momentoerrors.MomentoBaseError) {
	token, _, err := new(jwt.Parser).ParseUnverified(authToken, jwt.MapClaims{})
	if err != nil {
		return nil, momentoerrors.ConvertError(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &Endpoints{ControlEndpoint: reflect.ValueOf(claims[CONTROL_ENDPOINT_CLAIM_ID]).String(), CacheEndpoint: reflect.ValueOf(claims[CACHE_ENDPOINT_CLAIM_ID]).String()}, nil
	} else {
		return nil, momentoerrors.NewMomentoBaseError("InvalidArgumentError", "Invalid Auth token.")
	}
}
