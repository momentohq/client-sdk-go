package momento

import (
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

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

func prepareCacheName(r hasCacheName) (string, error) {
	name := r.cacheName()

	if len(strings.TrimSpace(name)) < 1 {
		err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
		return "", err
	}
	return name, nil
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
