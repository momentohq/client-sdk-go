package momento

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

var errUnexpectedGrpcResponse = errors.New("unexpected gRPC response")

type requester interface {
	hasCacheName
	initGrpcRequest(client scsDataClient) error
	makeGrpcRequest(client scsDataClient, metadata context.Context) (grpcResponse, error)
	interpretGrpcResponse() error
	requestName() string
}

type grpcResponse interface {
	Reset()
	String() string
}

type hasCacheName interface {
	cacheName() string
}

type hasKey interface {
	key() Bytes
}

type hasValue interface {
	value() Bytes
}

type hasScalarTTL interface {
	ttl() time.Duration
}

func prepareName(name string, label string) (string, error) {
	if len(strings.TrimSpace(name)) < 1 {
		errStr := fmt.Sprintf("%v cannot be empty", label)
		return "", momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, errStr, nil)
	}
	return name, nil
}

func prepareCacheName(r hasCacheName) (string, error) {
	return prepareName(r.cacheName(), "Cache name")
}

func prepareKey(r hasKey) ([]byte, error) {
	key := r.key().AsBytes()

	if len(key) == 0 {
		err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "key cannot be empty", nil)
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return key, nil
}

func prepareValue(r hasValue) ([]byte, momentoerrors.MomentoSvcErr) {
	value := r.value().AsBytes()
	if len(value) == 0 {
		err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "value cannot be empty", nil)
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return value, nil
}

func prepareTTL(r hasScalarTTL, defaultTtl time.Duration) (uint64, error) {
	ttl := defaultTtl
	if r.ttl() != time.Duration(0) {
		ttl = r.ttl()
	}
	return uint64(ttl.Milliseconds()), nil
}
