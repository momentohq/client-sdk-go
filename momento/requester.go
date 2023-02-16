package momento

// The requester interface is implemented by individual
// method request objects, for example SetRequest.
// requester.template is a template file to help implement
// a requester.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/utils"
)

func errUnexpectedGrpcResponse(r requester, grpcResp grpcResponse) momentoerrors.MomentoSvcErr {
	return momentoerrors.NewMomentoSvcErr(
		momentoerrors.InternalServerError,
		fmt.Sprintf(
			"%s request got an unexpected response %T '%s'",
			r.requestName(), grpcResp, grpcResp,
		),
		nil,
	)
}

type requester interface {
	hasCacheName
	initGrpcRequest(client scsDataClient) error
	makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error)
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
	key() Key
}

type hasValue interface {
	value() Value
}

type hasValues interface {
	values() []Value
}

type hasField interface {
	field() Value
}

type hasFields interface {
	fields() []Value
}

type hasItems interface {
	items() map[string]Value
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
	key := r.key().asBytes()

	if len(key) == 0 {
		err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "key cannot be empty", nil)
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return key, nil
}

func prepareField(r hasField) ([]byte, error) {
	field := r.field().asBytes()

	if len(field) == 0 {
		err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "field cannot be empty", nil)
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return field, nil
}

func prepareFields(r hasFields) ([][]byte, error) {
	var fields [][]byte
	for _, field := range r.fields() {
		if len(field.asBytes()) == 0 {
			err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "field cannot be empty", nil)
			return nil, convertMomentoSvcErrorToCustomerError(err)
		}
		fields = append(fields, field.asBytes())
	}
	return fields, nil
}

func prepareValue(r hasValue) ([]byte, momentoerrors.MomentoSvcErr) {
	value := r.value().asBytes()
	if len(value) == 0 {
		err := momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "value cannot be empty", nil)
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return value, nil
}

func prepareValues(r hasValues) ([][]byte, momentoerrors.MomentoSvcErr) {
	values := momentoBytesListToPrimitiveByteList(r.values())
	for i := range values {
		if len(values[i]) == 0 {
			return nil, momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError,
				"value in list cannot be empty",
				nil,
			)
		}
	}
	return values, nil
}

func prepareItems(r hasItems) (map[string][]byte, momentoerrors.MomentoSvcErr) {
	retMap := make(map[string][]byte)
	for k, v := range r.items() {
		retMap[k] = v.asBytes()
	}
	return retMap, nil
}

func prepareTTL(r hasScalarTTL, defaultTtl time.Duration) (uint64, error) {
	ttl := defaultTtl
	if r.ttl() != time.Duration(0) {
		ttl = r.ttl()
	}
	return uint64(ttl.Milliseconds()), nil
}

func prepareCollectionTtl(ttl utils.CollectionTTL, defaultTtl time.Duration) (uint64, bool) {
	ttlDuration := defaultTtl
	if ttl.Ttl != time.Duration(0) {
		ttlDuration = ttl.Ttl
	}

	return uint64(ttlDuration.Milliseconds()), ttl.RefreshTtl
}

func momentoBytesListToPrimitiveByteList(i []Value) [][]byte {
	var rList [][]byte
	for _, mb := range i {
		rList = append(rList, mb.asBytes())
	}
	return rList
}
