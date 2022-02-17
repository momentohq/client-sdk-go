package resolver

import (
	"reflect"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"github.com/golang-jwt/jwt/v4"
)

const (
	momentoControlEndpointPrefix = "control."
	momentoCacheEndpointPrefix   = "cache."
	controlEndpointClaimId       = "cp"
	cacheEndpointClaimId         = "c"
)

type Endpoints struct {
	ControlEndpoint string
	CacheEndpoint   string
}

func Resolve(request *models.ResolveRequest) (*Endpoints, momentoerrors.MomentoSvcErr) {
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
