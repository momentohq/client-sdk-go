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

type hasElements interface {
	elements() map[string]Value
}

type hasTTL interface {
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
	if err := validateNotEmpty(field, "field"); err != nil {
		return nil, err
	}
	return field, nil
}

func prepareFields(r hasFields) ([][]byte, error) {
	var fields [][]byte
	for _, valueField := range r.fields() {
		field := valueField.asBytes()
		if err := validateNotEmpty(field, "field"); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func prepareValue(r hasValue) ([]byte, momentoerrors.MomentoSvcErr) {
	if r.value() == nil {
		return []byte{}, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"value may not be nil",
			nil,
		)
	}
	return r.value().asBytes(), nil
}

func prepareValues(r hasValues) ([][]byte, momentoerrors.MomentoSvcErr) {
	values, err := momentoValuesToPrimitiveByteList(r.values())
	if err != nil {
		return [][]byte{}, err
	}
	return values, nil
}

func prepareItems(r hasElements) (map[string][]byte, error) {
	retMap := make(map[string][]byte)
	for k, v := range r.elements() {
		if v == nil {
			return map[string][]byte{}, momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError,
				"item values may not be nil",
				nil,
			)
		}
		if err := validateNotEmpty([]byte(k), "item keys"); err != nil {
			return nil, err
		}
		retMap[k] = v.asBytes()
	}
	return retMap, nil
}

func prepareTTL(r hasTTL, defaultTtl time.Duration) (uint64, error) {
	ttl := r.ttl()
	if r.ttl() == time.Duration(0) {
		ttl = defaultTtl
	}
	if ttl <= time.Duration(0) {
		return 0, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"ttl must be a non-zero positive value",
			nil,
		)
	}
	return uint64(ttl.Milliseconds()), nil
}

func momentoValuesToPrimitiveByteList(i []Value) ([][]byte, momentoerrors.MomentoSvcErr) {
	var rList [][]byte
	for _, mb := range i {
		if mb == nil {
			return [][]byte{}, momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError,
				"values may not be nil",
				nil,
			)
		}
		rList = append(rList, mb.asBytes())
	}
	return rList, nil
}

func validateNotEmpty(field []byte, label string) error {
	if len(field) == 0 {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError, fmt.Sprintf("%s cannot be empty", label), nil,
		)
	}
	return nil
}
