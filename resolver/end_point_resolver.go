package resolver

import (
	"reflect"

	"github.com/golang-jwt/jwt/v4"
)

const (
	MOMENTO_CONTROL_ENDPOINT_PREFIX = "control."
	MOMENTO_CACHE_ENDPOINT_PREFIX   = "cache."
	CONTROL_ENDPOINT_CLAIM_ID       = "cp"
	CACHE_ENDPOINT_CLAIM_ID         = "c"
)

type EndPoints struct {
	ContorlEndPoint string
	CacheEndPoint   string
}

func Resolve(authToken string, endPointOverride ...string) (*EndPoints, error) {
	if len(endPointOverride) != 0 {
		return &EndPoints{ContorlEndPoint: MOMENTO_CONTROL_ENDPOINT_PREFIX + endPointOverride[0], CacheEndPoint: MOMENTO_CACHE_ENDPOINT_PREFIX + endPointOverride[0]}, nil
	}
	return getEndPointsFromToken(authToken)
}

func getEndPointsFromToken(authToken string) (*EndPoints, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(authToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		ctEndPoint := reflect.ValueOf(claims[CONTROL_ENDPOINT_CLAIM_ID]).String()
		cEndPoint := reflect.ValueOf(claims[CACHE_ENDPOINT_CLAIM_ID]).String()
		return &EndPoints{ContorlEndPoint: ctEndPoint, CacheEndPoint: cEndPoint}, nil
	} else {
		return nil, err
	}
}
