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

type hasKeys interface {
	keys() []Key
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

type hasDictionaryElements interface {
	dictionaryElements() []DictionaryElement
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
		errStr := fmt.Sprintf("%v cannot be empty or blank", label)
		return "", buildError(momentoerrors.InvalidArgumentError, errStr, nil)
	}
	return name, nil
}

func prepareCacheName(r hasCacheName) (string, error) {
	return prepareName(r.cacheName(), "Cache name")
}

func prepareKey(r hasKey) ([]byte, error) {
	if err := validateNotEmpty(r.key(), "key"); err != nil {
		return nil, err
	}

	return r.key().asBytes(), nil
}

func prepareKeys(r hasKeys) ([][]byte, error) {
	var keys [][]byte
	for _, key := range r.keys() {
		if err := validateNotEmpty(key, "key"); err != nil {
			return nil, err
		}
		keys = append(keys, key.asBytes())
	}
	return keys, nil
}

func prepareField(r hasField) ([]byte, error) {
	if err := validateNotEmpty(r.field(), "field"); err != nil {
		return nil, err
	}
	return r.field().asBytes(), nil
}

func prepareFields(r hasFields) ([][]byte, error) {
	if r.fields() == nil {
		return nil, buildError(InvalidArgumentError, "fields cannot be nil", nil)
	}

	var fields [][]byte
	for _, field := range r.fields() {
		if err := validateNotEmpty(field, "field"); err != nil {
			return nil, err
		}
		fields = append(fields, field.asBytes())
	}
	return fields, nil
}

func prepareValue(r hasValue) ([]byte, error) {
	if err := validateNotNil(r.value(), "value"); err != nil {
		return []byte{}, err
	}
	return r.value().asBytes(), nil
}

func prepareValues(r hasValues) ([][]byte, error) {
	values, err := momentoValuesToPrimitiveByteList(r.values())
	if err != nil {
		return [][]byte{}, err
	}
	return values, nil
}

func prepareDictionaryElements(r hasDictionaryElements) ([]DictionaryElement, error) {
	for _, v := range r.dictionaryElements() {
		if err := validateNotNil(v.Value, "value"); err != nil {
			return nil, err
		}
		if err := validateNotEmpty(v.Field, "element field"); err != nil {
			return nil, err
		}
	}
	return r.dictionaryElements(), nil
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

func momentoValuesToPrimitiveByteList(values []Value) ([][]byte, error) {
	if values == nil {
		return nil, buildError(momentoerrors.InvalidArgumentError, "values cannot be nil", nil)
	}

	var rList [][]byte
	for _, mb := range values {
		if err := validateNotNil(mb, "value"); err != nil {
			return [][]byte{}, err
		}
		rList = append(rList, mb.asBytes())
	}
	return rList, nil
}

func validateNotEmpty(thing Value, label string) error {
	if err := validateNotNil(thing, label); err != nil {
		return err
	}

	if len(thing.asBytes()) == 0 {
		return buildError(
			momentoerrors.InvalidArgumentError, fmt.Sprintf("%v cannot be empty", label), nil,
		)
	}
	return nil
}

func validateNotNil(value Value, label string) error {
	if value == nil {
		return buildError(
			momentoerrors.InvalidArgumentError, fmt.Sprintf("%v cannot be nil", label), nil,
		)
	}

	return nil
}

func ElementsFromMapStringString(theMap map[string]string) []DictionaryElement {
	var elements []DictionaryElement
	for k, v := range theMap {
		elements = append(elements, DictionaryElement{
			Field: String(k),
			Value: String(v),
		})
	}
	return elements
}

func ElementsFromMapStringBytes(theMap map[string][]byte) []DictionaryElement {
	var elements []DictionaryElement
	for k, v := range theMap {
		elements = append(elements, DictionaryElement{
			Field: String(k),
			Value: Bytes(v),
		})
	}
	return elements
}

func ElementsFromMapStringValue(theMap map[string]Value) []DictionaryElement {
	var elements []DictionaryElement
	for k, v := range theMap {
		elements = append(elements, DictionaryElement{
			Field: String(k),
			Value: v,
		})
	}
	return elements
}
