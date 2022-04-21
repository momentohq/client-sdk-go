package utility

import (
	"strings"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

func IsCacheNameValid(cacheName string) bool {
	return len(strings.TrimSpace(cacheName)) != 0
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

func IsTtlValid(ttl uint64) momentoerrors.MomentoSvcErr {
	if ttl < 0 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "ttl cannot be negative", nil)
	}
	return nil
}
