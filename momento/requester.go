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

type hasTtl interface {
	ttl() *time.Duration
}

type hasRefreshTtl interface {
	refreshTtl() *bool
}

func buildError(errorCode string, errorMessage string, originalError error) MomentoError {
	return convertMomentoSvcErrorToCustomerError(
		momentoerrors.NewMomentoSvcErr(errorCode, errorMessage, originalError),
	)
}

func prepareName(name string, label string) (string, error) {
	if len(strings.TrimSpace(name)) < 1 {
		errStr := fmt.Sprintf("%v cannot be empty", label)
		return "", buildError(momentoerrors.InvalidArgumentError, errStr, nil)
	}
	return name, nil
}

func prepareElementName(name Value) ([]byte, error) {
	if name == nil {
		return nil, buildError(
			momentoerrors.InvalidArgumentError, "element name cannot be nil or empty", nil,
		)
	}

	// just validate not empty using prepareName
	nameBytes := name.asBytes()
	_, err := prepareName(string(nameBytes), "element name")
	if err != nil {
		return nil, err
	}

	return nameBytes, nil
}

func prepareCacheName(r hasCacheName) (string, error) {
	return prepareName(r.cacheName(), "Cache name")
}

func prepareKey(r hasKey) ([]byte, error) {
	if r.key() == nil {
		return nil, buildError(momentoerrors.InvalidArgumentError, "key cannot be nil or empty", nil)
	}

	key := r.key().asBytes()
	if len(key) == 0 {
		return nil, buildError(momentoerrors.InvalidArgumentError, "key cannot be nil or empty", nil)
	}
	return key, nil
}

func prepareField(r hasField) ([]byte, error) {
	if r.field() == nil {
		return nil, buildError(
			momentoerrors.InvalidArgumentError, "field cannot be nil or empty", nil,
		)
	}
	field := r.field().asBytes()
	if err := validateNotEmpty(field, "field"); err != nil {
		return nil, err
	}
	return field, nil
}

func prepareFields(r hasFields) ([][]byte, error) {
	if r.fields() == nil {
		return nil, buildError(momentoerrors.InvalidArgumentError, "fields cannot be nil or empty", nil)
	}
	var fields [][]byte
	for _, valueField := range r.fields() {
		if valueField == nil {
			return nil, buildError(momentoerrors.InvalidArgumentError, "fields cannot be nil or empty", nil)
		}
		field := valueField.asBytes()
		if err := validateNotEmpty(field, "field"); err != nil {
			return nil, buildError(momentoerrors.InvalidArgumentError, "fields cannot be nil or empty", nil)
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func prepareValue(r hasValue) ([]byte, momentoerrors.MomentoSvcErr) {
	if r.value() == nil {
		return []byte{}, buildError(
			momentoerrors.InvalidArgumentError, "value may not be nil", nil,
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

func prepareElements(r hasElements) (map[string][]byte, error) {
	retMap := make(map[string][]byte)
	for k, v := range r.elements() {
		if v == nil {
			return map[string][]byte{}, buildError(
				momentoerrors.InvalidArgumentError, "item values may not be nil", nil,
			)
		}
		if err := validateNotEmpty([]byte(k), "item keys"); err != nil {
			return nil, err
		}
		retMap[k] = v.asBytes()
	}
	return retMap, nil
}

func prepareTtl(r hasTtl, defaultTtl time.Duration) (uint64, error) {
	ttl := r.ttl()
	if *r.ttl() == time.Duration(0) {
		ttl = &defaultTtl
	}
	if *ttl <= time.Duration(0) {
		return 0, buildError(
			momentoerrors.InvalidArgumentError, "ttl must be a non-zero positive value", nil,
		)
	}
	return uint64(ttl.Milliseconds()), nil
}

func prepareRefreshTtl(r hasRefreshTtl) *bool {
	if r.refreshTtl() == nil {
		t := true
		return &t
	}
	return r.refreshTtl()
}

func momentoValuesToPrimitiveByteList(i []Value) ([][]byte, momentoerrors.MomentoSvcErr) {
	if i == nil {
		return [][]byte{}, buildError(momentoerrors.InvalidArgumentError, "values may not be nil", nil)
	}
	var rList [][]byte
	for _, mb := range i {
		if mb == nil {
			return [][]byte{}, buildError(momentoerrors.InvalidArgumentError, "values may not be nil", nil)
		}
		rList = append(rList, mb.asBytes())
	}
	return rList, nil
}

func validateNotEmpty(field []byte, label string) error {
	if len(field) == 0 {
		return buildError(
			momentoerrors.InvalidArgumentError, fmt.Sprintf("%s cannot be empty", label), nil,
		)
	}
	return nil
}
