package utility

import (
	"reflect"
	"strings"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

func IsCacheNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	return nil
}

func IsKeyValid(key interface{}) momentoerrors.MomentoSvcErr {
	if key == nil {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "key cannot be nil", nil)
	}
	return nil
}

func IsValueValid(value interface{}) momentoerrors.MomentoSvcErr {
	if value == nil {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "value cannot be nil", nil)
	}
	return nil
}

func EncodeKey(key interface{}) ([]byte, momentoerrors.MomentoSvcErr) {
	switch key.(type) {
	case string:
		return []byte(reflect.ValueOf(key).String()), nil
	case []byte:
		return reflect.ValueOf(key).Bytes(), nil
	default:
		// If target is not string or byte[] then use default gob encoder.
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"error encoding cache key only support []byte or string currently",
			nil,
		)
	}
}

func EncodeValue(value interface{}) ([]byte, momentoerrors.MomentoSvcErr) {
	switch value.(type) {
	case string:
		return []byte(reflect.ValueOf(value).String()), nil
	case []byte:
		return reflect.ValueOf(value).Bytes(), nil
	default:
		// If target is not string or byte[] then use default gob encoder.
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"error encoding cache value  only support []byte or string currently", nil,
		)
	}
}

func DecodeValue(rawBytes []byte, target interface{}) momentoerrors.MomentoSvcErr {
	switch target.(type) {
	case []byte:
		target = rawBytes
	case string:
		target = string(rawBytes)
	default:
		// If target is not string or byte[]
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"error decoding cache value only support []byte or string currently",
			nil,
		)
	}
	return nil
}
