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

	"github.com/momentohq/client-sdk-go/utils"

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
	elements() []Element
}

type hasTtl interface {
	ttl() time.Duration
}

type hasCollectionTtl interface {
	collectionTtl() *utils.CollectionTtl
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

func prepareElementValue(value Value) ([]byte, error) {
	if value == nil {
		return nil, buildError(
			momentoerrors.InvalidArgumentError, "element value cannot be nil", nil,
		)
	}

	// just validate not empty using prepareName
	_, err := prepareName(value.asString(), "element value")
	if err != nil {
		return nil, err
	}

	return value.asBytes(), nil
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

func prepareElements(r hasElements) ([]Element, error) {
	for _, v := range r.elements() {
		if v.ElemValue == nil || v.ElemField == nil {
			return nil, buildError(
				momentoerrors.InvalidArgumentError, "element fields and values may not be nil", nil,
			)
		}
		if err := validateNotEmpty(v.ElemField.asBytes(), "element field"); err != nil {
			return nil, err
		}
	}
	return r.elements(), nil
}

func prepareCollectionTtl(r hasCollectionTtl, defaultTtl time.Duration) (uint64, bool, error) {
	if r.collectionTtl() == nil {
		return uint64(defaultTtl.Milliseconds()), true, nil
	} else if r.collectionTtl().Ttl == time.Duration(0) {
		return uint64(defaultTtl.Milliseconds()), r.collectionTtl().RefreshTtl, nil
	} else if r.collectionTtl().Ttl <= time.Duration(0) {
		return 0, false, buildError(
			momentoerrors.InvalidArgumentError, "ttl must be a non-zero positive value", nil,
		)
	}
	return uint64(r.collectionTtl().Ttl.Milliseconds()), r.collectionTtl().RefreshTtl, nil
}

func prepareTtl(r hasTtl, defaultTtl time.Duration) (uint64, error) {
	ttl := r.ttl()
	if r.ttl() == time.Duration(0) {
		ttl = defaultTtl
	}
	if ttl <= time.Duration(0) {
		return 0, buildError(
			momentoerrors.InvalidArgumentError, "ttl must be a non-zero positive value", nil,
		)
	}
	return uint64(ttl.Milliseconds()), nil
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

func ElementsFromMapStringString(theMap map[string]string) []Element {
	var elements []Element
	for k, v := range theMap {
		elements = append(elements, Element{
			ElemField: String(k),
			ElemValue: String(v),
		})
	}
	return elements
}

func ElementsFromMapStringBytes(theMap map[string][]byte) []Element {
	var elements []Element
	for k, v := range theMap {
		elements = append(elements, Element{
			ElemField: String(k),
			ElemValue: Bytes(v),
		})
	}
	return elements
}

func ElementsFromMapStringValue(theMap map[string]Value) []Element {
	var elements []Element
	for k, v := range theMap {
		elements = append(elements, Element{
			ElemField: String(k),
			ElemValue: v,
		})
	}
	return elements
}
